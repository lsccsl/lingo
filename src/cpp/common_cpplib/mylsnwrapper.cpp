/**
* @file mylsnwrapper.cpp
* @brief multi thread select listerner wrapper,
*        wrapper tcp detail, provide network data to up layer
* @author linsc
*/
#include "mylsnwrapper.h"
#include "mylogex.h"
#include "channel.h"
#include "AutoLock.h"
extern "C"
{
    #include "mydefmempool.h"
}

/**
* @brief constructor
*/
mylsnwrapper::mylsnwrapper(uint32 thread_count, uint32 max_fd_count,
    uint32 max_msg_count, uint32 bufsz_reserve):hm_(NULL),
	    bufsz_reserve_(bufsz_reserve),
		real_thrd_count_(thread_count),
		max_msg_count_(max_msg_count)
{
	MYLOG_INFO(("mylsnwrapper::mylsnwrapper real_thrd_count_:%d max_msg_count_:%d bufsz_reserve_:%d",
		this->real_thrd_count_, this->max_msg_count_, this->bufsz_reserve_));

	this->hm_ = RhapsodyMemPoolConstruct();

	pthread_mutex_init(&this->map_thrd_protector_, NULL);
}

/**
* @brief destructor
*/
mylsnwrapper::~mylsnwrapper()
{
	MYLOG_TRACE(("mylsnwrapper::~mylsnwrapper"));
	
	uint32 i = 0;
	for(i = 0; i < this->vthrds_.size(); i ++)
	{
		MYLOG_INFO(("thrd %d stop", i));
		MyListernerDestruct(this->vthrds_[i].thrd_);
		this->vthrds_[i].thrd_ = NULL;
	}

	for(i = 0; i < this->vthrds_.size(); i ++)
	{
		MYLOG_INFO(("free thrd %d context", i));
		delete this->vthrds_[i].thrd_cxt_;
	}

	this->vthrds_.clear();

	{
		CAutoLock alock(&this->map_thrd_protector_);
		for(std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.begin();
			it != this->map_thrd_.end();
			it ++)
		{
			if(it->second)
				delete it->second;

			it->second = NULL;
		}

		this->map_thrd_.clear();
	}

	pthread_mutex_destroy(&this->map_thrd_protector_);

	if(this->hm_)
		MyMemePoolDestruct(this->hm_);
}

/**
* @brief 初始化
*/
int32 mylsnwrapper::init()
{
	MYLOG_DEBUG(("mylsnwrapper::init thrd count:%d", this->real_thrd_count_));
	/* 创建线程 */
	uint32 real_count = 0;
	uint32 i = 0;
	for(i = 0; i < this->real_thrd_count_; i ++)
	{
		HMYLISTERNER hlsn = MyListernerConstruct(this->hm_, this->max_msg_count_);
		if(NULL == hlsn)
			continue;

		thrd_cxt_t ti;
		ti.thrd_ = hlsn;
		ti.thrd_cxt_ = NULL;
		this->InitThrdHandle(ti.thrd_cxt_, i);

		this->vthrds_.push_back(ti);
		MYLOG_DEBUG(("add thread obj:%x", hlsn));
	}

	this->real_thrd_count_ = this->vthrds_.size();
	MYLOG_DEBUG(("real thrd count:%d", this->real_thrd_count_));

	/* 运行线程 */
	for(i = 0; i < this->real_thrd_count_; i ++)
	{
		MyListernerRun(this->vthrds_[i].thrd_);
		MYLOG_DEBUG(("run thread obj:%x", this->vthrds_[i]));
	}

	return 0;
}

/**
* @brief add tcp srv to listern
*/
int32 mylsnwrapper::add_tcp_srv_fd(int32 tcp_fd, uint32& thrd_index, int32 bauto_add)
{
	MYLOG_DEBUG(("mylsnwrapper::add_tcp_srv_fd tcp_fd:%d bauto_add:%d", tcp_fd, bauto_add));

	if(0 != this->_addToMapThrd(tcp_fd, TCP_ACCEPTOR_FD, thrd_index, E_FD_READ | E_FD_EXCEPTION))
		return -1;

	this->_set_accept_auto_add_to_listern_or_not(tcp_fd, bauto_add);

	return 0;
}

/**
* @brief add tcp cli/conn
*/
int32 mylsnwrapper::add_tcp_cli_fd(int32 tcp_fd, uint32& thrd_index)
{
	MYLOG_DEBUG(("mylsnwrapper::add_tcp_cli_fd tcp_fd:%d", tcp_fd));

	if(0 != this->_addToMapThrd(tcp_fd, TCP_CONNECTOR_CLI_FD, thrd_index, E_FD_READ | E_FD_WRITE | E_FD_EXCEPTION))
		return -1;

	return 0;
}

/**
* @brief add udp fd to listener
*/
int32 mylsnwrapper::add_udp_fd(int32 udp_fd, uint32& thrd_index)
{
	MYLOG_DEBUG(("mylsnwrapper::add_udp_fd udp_fd:%d", udp_fd));
	assert(0);
	return 0;
}

/**
* @brief close and remove from listener
*/
int32 mylsnwrapper::del_and_close_fd(int32 fd)
{
	MYLOG_ERR(("mylsnwrapper::del_and_close_fd fd:%d", fd));

	HMYLISTERNER hlsn = NULL;
	if(0 != this->_getFdThrd(fd, hlsn))
	{
		MYLOG_WARN(("can't find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	this->_push_inter_msg(hlsn, fd, E_DEL_AND_CLOSE_FD, NULL);

	return 0;
}

/**
* @brief 将要发送的数据压进缓冲区,并更改epoll的wait状态,等于合适的时候再发送数据
* @param fd:句柄
* @param data:数据缓冲区
* @param data_sz:data的大小
*/
int32 mylsnwrapper::data_tcp_out(int32 fd, std::vector<uint8>& data)
{
	MYLOG_DEBUG(("CMyThrdPoll::data_tcp_out fd:%d sz:%d", fd, data.size()));

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != this->_getFdThrd(fd, hlsn))
	{
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(mylsnwrapper::thrd_msg_user_output));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("alloc msg body err"));
		return -1;
	}
	mylsnwrapper::thrd_msg_user_output * mb = new (place_ment)mylsnwrapper::thrd_msg_user_output;

	mb->msg.swap(data);

	if(0 != this->_push_inter_msg(hlsn, fd, E_USER_OUTPUT_TCP, mb))
	{
		if(mb)
		{
			mb->~thrd_msg_user_output();
			MyMemPoolFree(this->hm_, mb);
		}
	}

	return 0;
}
/**
* @brief 直接将数据发送出去,不缓存,用户保证对同一个fd的读写是串行的(或者加锁,或者该fd永远只被一个线程写)
*/
int32 mylsnwrapper::data_tcp_out_sync(int32 fd, std::vector<uint8>& data)
{
	MYLOG_DEBUG(("mylsnwrapper::data_tcp_out_sync fd:%d sz:%d", fd, data.size()));

	if(0 == data.size())
	{
		MYLOG_INFO(("data buf is 0, don't send"));
		return 0;
	}

	int32 data_len = data.size();
	uint32 data_pos = 0;

	while(data_len > 0)
	{
		MYLOG_DEBUG(("data_len:%d data_pos:%d", data_len, data_pos));
		int32 ret = CChannel::TcpWrite(fd, &data[data_pos], data_len);
		if(ret < 0)
		{
			MYLOG_INFO(("write to peer err data_len:%d ret:%d", data_len, ret));
			return -1;
		}

		data_len -= ret;
		data_pos += ret;
	}

	return 0;
}

/**
* @brief 给指定的线程发消息
*/
int32 mylsnwrapper::push_msg(void * msg, uint32 thrd_to, int32 thrd_from)
{
	MYLOG_DEBUG(("mylsnwrapper::push_msg from:%d to:%d %x", thrd_from, thrd_to, msg));

	if(thrd_to >= this->vthrds_.size())
	{
		MYLOG_WARN(("out of range %d - %d", thrd_to, this->vthrds_.size()));
		return -1;
	}

	MYLOG_DEBUG(("malloc msg"));

	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(mylsnwrapper::thrd_user_msg_t_));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("no memory"));
		return -1;
	}
	mylsnwrapper::thrd_user_msg_t_ * m = new (place_ment)mylsnwrapper::thrd_user_msg_t_;
	m->thrd_from_ = thrd_from;
	m->thrd_to_ = thrd_to;
	m->msg_body_ = msg;

	MYLOG_DEBUG(("real push msg"));

	MyListernerAddMsg(this->vthrds_[thrd_to].thrd_, m, (unsigned long)this, mylsnwrapper::_thrdsUserMsgCb);

	return 0;
}

/**
 * @brief 添加定时器
 */
int32 mylsnwrapper::add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period)
{
	MYLOG_DEBUG(("mylsnwrapper::add_time time_second:%d thrd_index:%d timer_data:%x", time_second, thrd_index, timer_data));
	if(thrd_index >= this->vthrds_.size())
	{
		MYLOG_WARN(("out of range %d - %d", thrd_index, this->vthrds_.size()));
		return -1;
	}

	mytimer_node_t node = {0};

	node.context_data = (unsigned long)this->vthrds_[thrd_index].thrd_cxt_;
	node.timer_user_data = timer_data;
	node.first_expire.tv_sec = time_second;
	/* 是否周期触发 */
	if(period)
		node.period.tv_sec = time_second;
	node.timeout_cb = mylsnwrapper::__handle_timeout;

	timer_id = MyListernerAddTimer(this->vthrds_[thrd_index].thrd_, &node);

	return 0;
}

/**
 * @brief 删除定时器
 */
int32 mylsnwrapper::del_time(uint32 thrd_index, HTIMERID timer_id)
{
	MYLOG_DEBUG(("mylsnwrapper::del_time"));

	MyListernerDelTimer(this->vthrds_[thrd_index].thrd_, timer_id);

	return 0;
}

/**
* @brief handle network input event callback
*/
int mylsnwrapper::__handle_tcp_input(unsigned long context_data, int fd)
{
	MYLOG_DEBUG(("mylsnwrapper::__handle_input context_data:%x fd:%d ==================", context_data, fd));

	mylsnwrapper * pthis = (mylsnwrapper *)context_data;

	/* 由于fd对应的fi是会在同一个线程的上下文中被删除,被访问,所以可以不加锁
	* read and,when egain,stop
	*/
	mylsnwrapper::fd_info_t * pfi = NULL;
	if(0 != pthis->_getFdInfoFromMapThrd(fd, pfi))
	{
		//this->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info =================="));
		return -1;
	}
	if(NULL == pfi)
	{
		MYLOG_WARN(("fd info is null =================="));
		return -1;
	}

	MYLOG_DEBUG(("fi buf size:%d pos:%d", pfi->rbuf_.size(), pfi->rcvpos_));

	int32 need_close = 0;
	do
	{
		if(pfi->rcvpos_ >= pfi->rbuf_.size())
		{
			pfi->rbuf_.resize(pfi->rbuf_.size() + 1024);
		}

		int32 ret = CChannel::TcpRead(pfi->fd_, &pfi->rbuf_[pfi->rcvpos_], pfi->rbuf_.size() - pfi->rcvpos_);

		MYLOG_DEBUG(("read ret:%d", ret));

		if(ret < 0)
		{
			MYLOG_ERR(("err need close"));
			need_close = 1;
			break;
		}
		else if(0 == ret)
		{
			MYLOG_DEBUG(("no data"));
			break;
		}
		else
		{
			pfi->rcvpos_ += ret;

			MYLOG_DUMP_BIN(&pfi->rbuf_[0], pfi->rcvpos_);

		}
	}while(1);

	uint32 recved = 0;
	MYLOG_DEBUG(("fd:%d rcvpos:%d readpos:%d %x size:%d", pfi->fd_, pfi->rcvpos_, pfi->readpos_, &pfi->rbuf_[0], pfi->rbuf_.size()));

	/* we call back data in event here, 
	* the same fd's event will be call back in the same thread, 
	* so app all always think it was in the same thread context 
	* when call back, we don't take any lock
	*/
	MYLOG_DEBUG(("thrd_context_data_:%x", pfi->thrd_context_data_));

	MYLOG_DEBUG(("fd:%d local[%s:%d] remote[%s:%d]",
		fd,
		pfi->local_ip_.c_str(), pfi->local_port_,
		pfi->remote_ip_.c_str(), pfi->remote_port_));

	if(pfi->thrd_context_data_ && pfi->rcvpos_ > pfi->readpos_)
		pfi->thrd_context_data_->data_tcp_in(pfi->fd_, &pfi->rbuf_[pfi->readpos_], pfi->rcvpos_ - pfi->readpos_, recved);

	MYLOG_DEBUG(("_data_in end recved:%d readpos_:%d rcvpos_:%d bufsz_reserve_:%d", recved, pfi->readpos_, pfi->rcvpos_, pthis->bufsz_reserve_));
	pfi->readpos_ += recved;
	if(pfi->readpos_ >= pfi->rcvpos_)
	{
		MYLOG_DEBUG(("all data in buf has been recved by app, now pos go back"));
		pfi->readpos_ = 0;
		pfi->rcvpos_ = 0;
	}
	else if(((pfi->rcvpos_ - pfi->readpos_) < pthis->bufsz_reserve_) && (pfi->readpos_ > pthis->bufsz_reserve_))
	{
		/* 未被上层接收的缓冲小于指定大小bufsz_reserve_,并且pfi->rrpos_位置大于bufsz_reserve_,将缓冲区整到开头去 */
		MYLOG_DEBUG(("all reserve data copy to the front of the buf, and pos go back"));

		/* move buf_sz back */
		memmove(&pfi->rbuf_[0], &pfi->rbuf_[pfi->readpos_], pfi->rcvpos_ - pfi->readpos_);
		pfi->rcvpos_ -= pfi->readpos_;
		pfi->readpos_ = 0;
	}
	else
	{
		MYLOG_DEBUG(("not all been recv,but don't move buf"));
	}

	if(need_close)
	{
		MYLOG_ERR(("fd need to close"));

		pthis->_delFromMapThrd(fd);
	}

	MYLOG_DEBUG(("mylsnwrapper::_thrdsDoInput %d end =================", fd));

	return 0;
}
int mylsnwrapper::__handle_tcpsrv_input(unsigned long context_data, int fd)
{
	mylsnwrapper * pthis = (mylsnwrapper *)context_data;

	do
	{
		MYLOG_DEBUG(("accept loop"));

		int8 ip[32] = {0};
		unsigned short port = 0;
		int32 accept_fd = -1;
		accept_fd = CChannel::TcpAccept(fd);
		if(accept_fd < 0)
		{
			MYLOG_DEBUG(("accept fail, loop break"));
			break;
		}

		CChannel::set_no_block(accept_fd);

		MYLOG_DEBUG(("accept new connection %d", accept_fd));

		fd_info_t * pfi = NULL;
		pthis->_getFdInfoFromMapThrd(fd, pfi);
		if(NULL == pfi)
		{
			MYLOG_WARN(("accept end fd info is null"));
			continue;
		}

		if(pfi->bauto_add_to_listern_)
		{
			MYLOG_DEBUG(("need auto add to listern"));

			uint32 thrd_index = 0;
			pthis->_addToMapThrd(accept_fd, TCP_CONNECTOR_FD, thrd_index, E_FD_READ | E_FD_WRITE | E_FD_EXCEPTION, fd);
			MYLOG_DEBUG(("_addToMapThrd end"));

			/* 推入消息
			*/
			pthis->_push_inter_msg(pthis->vthrds_[thrd_index].thrd_, accept_fd, E_NEW_CONN, NULL);
		}
		else
		{
			MYLOG_DEBUG(("not need auto add to listern"));
			pfi->thrd_context_data_->accept_have_conn(fd, accept_fd);
		}

		MYLOG_DEBUG(("accept end"));
	}while(1);

	return 0;
}

/** 
* @brief handle out event callback
*/
int mylsnwrapper::__handle_tcp_output(unsigned long context_data, int fd)
{
	MYLOG_DEBUG(("mylsnwrapper::__handle_output context_data:%x fd:%d", context_data, fd));
	mylsnwrapper * pthis = (mylsnwrapper *)context_data;
	return pthis->_inter_tcp_output(fd);
}

/**
* @brief handle exception callback
*/
int mylsnwrapper::__handle_exception(unsigned long context_data, int fd)
{
	MYLOG_ERR(("mylsnwrapper::__handle_exception context_data:%x fd:%d", context_data, fd));
	mylsnwrapper * pthis = (mylsnwrapper *)context_data;
	pthis->_delFromMapThrd(fd);
	return 0;
}

/**
* @brief handle timer callback
*/
int mylsnwrapper::__handle_timeout(unsigned long context_data,  unsigned long timer_user_data, HTIMERID timerid)
{
	//MYLOG_DEBUG(("mylsnwrapper::__handle_timeout context_data:%x timer_user_data:%x timerid:%x",
	//	context_data, timer_user_data, timerid));

	LsnThrdHandle * thrd_context_data = (LsnThrdHandle *)context_data;
	if(NULL == thrd_context_data)
	{
		MYLOG_WARN(("err context data"));
		return -1;
	}

	thrd_context_data->time_out(timer_user_data, timerid);

	return 0;
}

/**
* @brief 获取fd的上下文数据,由于上下文数据只会被固定的一个线程所操作,所以在操作上下文数据时不加锁
*/
int32 mylsnwrapper::_getFdInfoFromMapThrd(int32 fd, mylsnwrapper::fd_info_t*& i)
{
	MYLOG_DEBUG(("mylsnwrapper::_getFdInfoFromMapThrd %d", fd));

	CAutoLock alock(&this->map_thrd_protector_);
	std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
	if(this->map_thrd_.end() == it)
	{
		MYLOG_WARN(("fd:%d is not in map", fd));
		return -1;
	}

	i = it->second;

	MYLOG_DEBUG(("thrd_context_data_:%x", i->thrd_context_data_));

	return 0;
}
int32 mylsnwrapper::_getFdThrd(int32 fd, HMYLISTERNER& hlsn)
{
	MYLOG_DEBUG(("mylsnwrapper::_getFdThrd %d", fd));

	CAutoLock alock(&this->map_thrd_protector_);
	std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
	if(this->map_thrd_.end() == it)
	{
		MYLOG_WARN(("fd:%d is not in map", fd));
		return -1;
	}

	MYLOG_DEBUG(("thrd_context_data_:%x", it->second->thrd_context_data_));

	hlsn = it->second->hlsn_;

	return 0;
}
/**
* @brief 添加至 map_thrd_ need lock
*/
int32 mylsnwrapper::_addToMapThrd(int32 fd, uint32 fd_type, uint32& thrd_index, uint32 mask, int32 fd_master)
{
	MYLOG_DEBUG(("mylsnwrapper::_addToMapThrd fd:%d fd_type:%d fd_master:%d", fd, fd_type, fd_master));

	CAutoLock alock(&this->map_thrd_protector_);

	std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
	if(this->map_thrd_.end() != it)
	{
		MYLOG_WARN(("fd:%d is in map", fd));
		return -1;
	}

	thrd_index = fd % this->vthrds_.size();

	mylsnwrapper::fd_info_t * i = new mylsnwrapper::fd_info_t;
	i->fd_ = fd;
	i->type_ = fd_type;

	/* fd was process by one thread only, 
	*/
	i->hlsn_ = this->vthrds_[thrd_index].thrd_;
	i->thrd_context_data_ = this->vthrds_[thrd_index].thrd_cxt_;
	i->rcvpos_ = 0;
	i->readpos_ = 0;
	i->rbuf_.resize(4096);

	/* fill another sock info
	*/
	CChannel::getSocketName(fd, i->local_ip_, i->local_port_);
	CChannel::getPeerName(fd, i->remote_ip_, i->remote_port_);
	i->fd_master_ = fd_master;
	i->bauto_add_to_listern_ = 1;

	MYLOG_DEBUG(("fd:%d local[%s:%d] remote[%s:%d]",
		fd,
		i->local_ip_.c_str(), i->local_port_,
		i->remote_ip_.c_str(), i->remote_port_));

	this->map_thrd_[fd] = i;

	MYLOG_DEBUG(("i->thrd_context_data_ %x", i->thrd_context_data_));

	switch(fd_type)
	{
	case TCP_ACCEPTOR_FD:
		{
			MYLOG_DEBUG(("fd is tcp type"));
			event_handle_t evt_handle = {
				mylsnwrapper::__handle_tcpsrv_input,
				NULL,
				mylsnwrapper::__handle_exception,
			};

			evt_handle.context_data = (unsigned long)this;
			MyListernerAddFD(i->hlsn_, fd, (E_HANDLE_SET_MASK)mask, &evt_handle);
		}
		break;

	case TCP_CONNECTOR_FD:
	case TCP_CONNECTOR_CLI_FD:
		{
			MYLOG_DEBUG(("fd is tcp type"));
			event_handle_t evt_handle = {
				mylsnwrapper::__handle_tcp_input,
				mylsnwrapper::__handle_tcp_output,
				mylsnwrapper::__handle_exception,
			};

			evt_handle.context_data = (unsigned long)this;
			MyListernerAddFD(i->hlsn_, fd, (E_HANDLE_SET_MASK)mask, &evt_handle);
		}
		break;

	case UDP_FD:
	default:
		{
			MYLOG_DEBUG(("fd is not tcp type"));
		}
		break;
	}

	return 0;
}

/**
* @brief 删除
*/
int32 mylsnwrapper::_delFromMapThrd(int32 fd)
{
	MYLOG_DEBUG(("mylsnwrapper::_delFromMapThrd %d", fd));

	LsnThrdHandle * h = NULL;

	{
		CAutoLock alock(&this->map_thrd_protector_);
		std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
		if(this->map_thrd_.end() == it)
		{
			MYLOG_WARN(("fd:%d is not in map", fd));
			return -1;
		}

		if(NULL == it->second)
		{
			MYLOG_WARN(("fd info is null,bug here..."));
			this->map_thrd_.erase(it);
			return -1;
		}

		MyListernerDelFD(it->second->hlsn_, fd);

		MYLOG_INFO(("fd:%d local[%s:%d] remote[%s:%d]",
			fd,
			it->second->local_ip_.c_str(), it->second->local_port_,
			it->second->remote_ip_.c_str(), it->second->remote_port_));

		h = it->second->thrd_context_data_;
		MYLOG_DEBUG(("thrd_context_data_:%x", it->second->thrd_context_data_));
		delete it->second;

		this->map_thrd_.erase(it);
	}

	MYLOG_DEBUG(("%x", h));

	if(h)
	{
		MYLOG_DEBUG(("close fd call back"));
		h->close_fd(fd);
		MYLOG_DEBUG(("close fd call back end"));
	}
	else
	{
		MYLOG_INFO(("no user context_data"));
	}

	CChannel::CloseFd(fd);

	return 0;
}

/**
* @brief set auto listen or not
*/
int32 mylsnwrapper::_set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add)
{
	MYLOG_DEBUG(("mylsnwrapper::_set_accept_auto_add_to_listern_or_not fd:%d bauto_add:%d", fd, bauto_add));

	CAutoLock alock(&this->map_thrd_protector_);
	std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
	if(this->map_thrd_.end() == it)
	{
		MYLOG_WARN(("fd:%d is not in map", fd));
		return -1;
	}

	if(NULL == it->second)
	{
		MYLOG_WARN(("fd info is null"));
		return -1;
	}

	it->second->bauto_add_to_listern_ = bauto_add;
	return 0;
}

/**
* @brief 推入线程池内部线程之间的消息
*/
int32 mylsnwrapper::_push_inter_msg(HMYLISTERNER hlsn, int32 fd, uint32 ev, void * pb)
{
	MYLOG_DEBUG(("mylsnwrapper::_push_inter_msg %x fd:%d event:%d pb:%x", hlsn, fd, ev, pb));

	/* 推入消息
	*/
	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(mylsnwrapper::thrd_msg_t));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("no memory"));
		return -1;
	}
	mylsnwrapper::thrd_msg_t * msg = new(place_ment)mylsnwrapper::thrd_msg_t;

	if(NULL == msg)
	{
		MYLOG_ERR(("impossible err"));
		return -1;
	}

	msg->fd_ = fd;
	msg->event_ = ev;
	msg->msg_body_ = pb;
	MyListernerAddMsg(hlsn, msg, (unsigned long)this, mylsnwrapper::_thrdsMsgCb);

	return 0;
}

/**
* @brief 处理消息队列里的消息,this function was invoke in the multi thrds context
*/
int32 mylsnwrapper::_thrdsMsgCb(unsigned long context_data, void * msg)
{
	MYLOG_DEBUG(("mylsnwrapper::_thrdsMsgCb msg:%x", msg));
	mylsnwrapper * pthis = (mylsnwrapper *)context_data;
	if(NULL == pthis || NULL == msg)
	{
		MYLOG_WARN(("pthis is null or msg is null"));
		return -1;
	}

	mylsnwrapper::thrd_msg_t * e = (mylsnwrapper::thrd_msg_t *)msg;
	if(NULL == e)
	{
		MYLOG_WARN(("msg is null"));
		return -1;
	}

	MYLOG_DEBUG(("event:%d fd:%d", e->event_, e->fd_));

	switch(e->event_)
	{
	case E_USER_OUTPUT_TCP:
		pthis->_thrdsDoUserOutputTcp(e->fd_, e);
		break;

	case E_NEW_CONN:
		pthis->_thrdsDoNewConn(e->fd_);
		break;

	case E_DEL_AND_CLOSE_FD:
		pthis->_delFromMapThrd(e->fd_);
		break;

	default:
		MYLOG_WARN(("we have unkown event..."));
		break;
	}

	e->free_body(pthis->hm_);
	MyMemPoolFree(pthis->hm_, e);

	MYLOG_DEBUG(("_thrdsMsgCb end\r\n\r\n"));

	return 0;
}

/**
* @brief 用户消息回调函数
*/
int32 mylsnwrapper::_thrdsUserMsgCb(unsigned long context_data, void * msg)
{
	MYLOG_DEBUG(("mylsnwrapper::_thrdsUserMsgCb msg:%x", msg));
	mylsnwrapper * tp = (mylsnwrapper *)context_data;
	if(NULL == tp || NULL == msg)
	{
		MYLOG_WARN(("tp is null or msg is null"));
		return -1;
	}

	mylsnwrapper::thrd_user_msg_t_ * m = (mylsnwrapper::thrd_user_msg_t_ *)msg;
	if(NULL == m)
	{
		MYLOG_DEBUG(("msg is null"));
		return -1;
	}

	/* 消息要入口时被限制了,一定会有线程处理 */
	assert(m->thrd_to_ < tp->vthrds_.size());

	MYLOG_DEBUG(("m->msg_body_:%x, m->thrd_from_:%d, m->thrd_to_:%d", m->msg_body_, m->thrd_from_, m->thrd_to_));
	tp->vthrds_[m->thrd_to_].thrd_cxt_->msg_callback(m->thrd_from_, m->msg_body_);

	/* 交由产生m->msg_body_的代码来释放m->msg_body_,此处不释放m->msg_body_ */
	m->msg_body_ = NULL;

	/* 没有析构,无需要呼叫 */
	//m->~thrd_user_msg_t_();
	MyMemPoolFree(tp->hm_, m);
	return 0;
}

/**
* @brief do new connection in event, 在线程池中某个线程的上下文中被呼叫
*/
int32 mylsnwrapper::_thrdsDoNewConn(int32 fd)
{
	MYLOG_DEBUG(("mylsnwrapper::_thrdsDoNewConn %d", fd));

	mylsnwrapper::fd_info_t * pfi = NULL;
	this->_getFdInfoFromMapThrd(fd, pfi);
	if(NULL == pfi)
		return -1;

	MYLOG_INFO(("fd:%d local[%s:%d] remote[%s:%d]",
		fd,
		pfi->local_ip_.c_str(), pfi->local_port_,
		pfi->remote_ip_.c_str(), pfi->remote_port_));

	if(pfi->thrd_context_data_)
	{
		MYLOG_DEBUG(("new connect cb"));
		pfi->thrd_context_data_->tcp_conn(fd, pfi->fd_master_);

		/* 缓冲区里是否已经有数据了,如果有数据,要通知 */
		if(pfi->rcvpos_ > pfi->readpos_)
		{
			MYLOG_DEBUG(("before connect cb,there has been data in this channel already,process it"));

			uint32 recved = 0;
			pfi->thrd_context_data_->data_tcp_in(fd, &pfi->rbuf_[pfi->readpos_], pfi->rcvpos_ - pfi->readpos_, recved);
			pfi->readpos_ += recved;
		}
	}
	else
		MYLOG_INFO(("new connection in,but context data is null, so no call back"));
	return 0;
}

/**
* @brief do output, 在线程池中某个线程的上下文中被呼叫
*/
int32 mylsnwrapper::_thrdsDoUserOutputTcp(int32 fd, mylsnwrapper::thrd_msg_t * msg)
{
	MYLOG_DEBUG(("mylsnwrapper::_thrdsDoUserOutputTcp fd:%d msg:%x", fd, msg));

	if(NULL == msg)
	{
		MYLOG_WARN(("msg is null"));
		return -1;
	}
	if(NULL == msg->msg_body_)
	{
		MYLOG_WARN(("msg body is null"));
		return -1;
	}

	thrd_msg_user_output * mb = (thrd_msg_user_output *)msg->msg_body_;

	mylsnwrapper::fd_info_t * pfi = NULL;
	this->_getFdInfoFromMapThrd(fd, pfi);
	if(NULL == pfi)
		return -1;

	MYLOG_DEBUG(("size:%d", pfi->sbuf_.size()));

	if(0 == pfi->sbuf_.size())
	{
		MYLOG_DEBUG(("send buf has nothing,swap it"));
		pfi->sbuf_.swap(mb->msg);
	}
	else
	{
		MYLOG_DEBUG(("send buf has something,append it"));
		pfi->sbuf_.insert(pfi->sbuf_.end(), mb->msg.begin(), mb->msg.end());
	}

	this->_inter_tcp_output(fd);

	return 0;
}

/**
* @brief 将缓存的tcp内容发出去
*/
int32 mylsnwrapper::_inter_tcp_output(int32 fd)
{
	MYLOG_DEBUG(("mylsnwrapper::_inter_tcp_output fd:%d", fd));

	/* 循环发送
	* 送完为止
	*/
	mylsnwrapper::fd_info_t * pfi = NULL;
	this->_getFdInfoFromMapThrd(fd, pfi);
	if(NULL == pfi)
		return -1;

	int32 data_len = pfi->sbuf_.size();
	if(data_len <= 0)
	{
		MYLOG_DEBUG(("have not date to write"));
		return 0;
	}

	uint32 data_pos = 0;

	while(data_len > 0)
	{
		int32 ret = CChannel::TcpWrite(pfi->fd_, &pfi->sbuf_[data_pos], data_len);
		if(ret < 0)
		{
			break;
		}

		data_len -= ret;
		data_pos += ret;
	}

	//int32 need_close = 0;
	/* move unsend data to buf head */
	if(data_len <= 0)
	{
		MYLOG_DEBUG(("all data has been send, mod epoll to read"));
		pfi->sbuf_.clear();
	}
	else
	{
		MYLOG_INFO(("not all data has been send, not mod to read"));
		memmove(&pfi->sbuf_[0], &pfi->sbuf_[data_pos], data_len);
		pfi->sbuf_.resize(data_len);
	}

	//if(need_close)
	//{
	//	this->del_and_close_fd(fd);
	//	return 0;
	//}

	return 0;
}






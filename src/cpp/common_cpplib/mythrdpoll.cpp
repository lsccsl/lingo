/**
 * @file mythrdpoll.cpp 
 * @brief 处理sock消息线程池
 *
 * @author linshaochuan
 */
#include "mythrdpoll.h"
#include <assert.h>
#include <errno.h>
#include "myepoll.h"
#include "mylogex.h"
#include "AutoLock.h"
#include "channel.h"
extern "C"
{
	#include "mydefmempool.h"
}
#define MY_MTU (1024)

/**
 * @brief constructor
 */
CMyThrdPoll::CMyThrdPoll(const uint32 thread_count,
	const uint32 max_fd_count,
	const uint32 max_msg_count,
	const uint32 bufsz_reserve,
	const uint32 epoll_thrd_count):
		real_thrd_count_(thread_count),
		epoll_(new CMyEPoll(max_fd_count, epoll_thrd_count)),
		max_msg_count_(max_msg_count),
		bufsz_reserve_(bufsz_reserve)
{
	MYLOG_INFO(("real_thrd_count_:%d max_msg_count_:%d bufsz_reserve_:%d",
		this->real_thrd_count_, this->max_msg_count_, this->bufsz_reserve_));

	this->hm_ = RhapsodyMemPoolConstruct();

	pthread_mutex_init(&this->map_thrd_protector_, NULL);

	if(NULL == epoll_)
		MYLOG_WARN(("create poll obj err"));
}

/**
 * @brief destructor
 */
CMyThrdPoll::~CMyThrdPoll()
{
	MYLOG_DEBUG(("CMyThrdPoll::~CMyThrdPoll"));

	if(this->epoll_)
		delete this->epoll_;

	uint32 i = 0;
	for(i = 0; i < this->thrds_.size(); i ++)
	{
		MYLOG_INFO(("thrd %d stop", i));
		MyListernerDestruct(this->thrds_[i].thrd_);
	}

	for(i = 0; i < this->thrds_.size(); i ++)
	{
		MYLOG_INFO(("free thrd %d context", i));
		delete this->thrds_[i].thrd_context_data_;
	}

	this->thrds_.clear();

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
 * @brief 让用户反初始化线程数据
 */
//int32 CMyThrdPoll::UninitThrdHandle(ThrdHandle * thrd_handle, uint32 thrd_index)
//{
//	MYLOG_DEBUG(("CMyThrdPoll::UninitThrdHandle"));
//	return 0;
//}

/**
 * @brief 初始化
 */
int32 CMyThrdPoll::init()
{
	MYLOG_DEBUG(("CMyThrdPoll::init"));

	if(NULL == epoll_)
	{
		MYLOG_WARN(("create poll obj err"));
		return -1;
	}

	MYLOG_DEBUG(("thrd count:%d", this->real_thrd_count_));
	/* 创建线程 */
	uint32 real_count = 0;
	uint32 i = 0;
	for(i = 0; i < this->real_thrd_count_; i ++)
	{
		HMYLISTERNER hlsn = MyListernerConstruct(this->hm_, this->max_msg_count_);
		if(NULL == hlsn)
			continue;

		thrd_info_t ti;
		ti.thrd_ = hlsn;
		ti.thrd_context_data_ = NULL;
		this->InitThrdHandle(ti.thrd_context_data_, i);

		this->thrds_.push_back(ti);
		MYLOG_DEBUG(("add thread obj:%x", hlsn));
	}

	this->real_thrd_count_ = this->thrds_.size();
	MYLOG_DEBUG(("real thrd count:%d", this->real_thrd_count_));

	/* 运行线程 */
	for(i = 0; i < this->real_thrd_count_; i ++)
	{
		MyListernerRun(this->thrds_[i].thrd_);
		MYLOG_DEBUG(("run thread obj:%x", this->thrds_[i]));
	}

	return 0;
}

/**
 * @brief run all thrd
 */
int32 CMyThrdPoll::work_loop(int32 timeout)
{
	/* when accept msg coming, it is processed in this context */
	return this->epoll_->work_loop(timeout);
}

/**
 * @brief 将要发送的数据压进缓冲区,此函数供外部调用,为了防止调用者处于不同的线程上下文件,采用压消息的方式...
 * @param fd:句柄
 * @param data:数据缓冲区
 * @param data_sz:data的大小
 */
int32 CMyThrdPoll::data_tcp_out(int32 fd, std::vector<uint8>& data)
{
	MYLOG_DEBUG(("CMyThrdPoll::data_tcp_out fd:%d sz:%d", fd, data.size()));

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != this->_getFdThrd(fd, hlsn))
	{
		//this->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(thrd_msg_user_output));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("alloc msg body err"));
		return -1;
	}
	thrd_msg_user_output * mb = new (place_ment)thrd_msg_user_output;

	mb->msg.swap(data);

	if(0 != this->_push_inter_msg(hlsn, fd, E_USER_OUTPUT_TCP, mb))
	{
		if(mb)
		{
			mb->~thrd_msg_user_output();
			MyMemPoolFree(this->hm_, mb);
			//delete mb;
		}
	}

	return 0;
}
/**
 * @brief 直接将数据发送出去,不缓存,用户保证对同一个fd的读写是串行的(或者加锁,或者该fd永远只被一个线程写)
 */
int32 CMyThrdPoll::data_tcp_out_sync(int32 fd, std::vector<uint8>& data)
{
	MYLOG_DEBUG(("CMyThrdPoll::data_tcp_out_sync fd:%d sz:%d", fd, data.size()));

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

		int32 ret = 0;
		if(data_len < MY_MTU)
			ret = CChannel::TcpWrite(fd, &data[data_pos], data_len);
		else
			ret = CChannel::TcpWrite(fd, &data[data_pos], MY_MTU);
		MYLOG_DEBUG(("write ret:%d", ret));
		if(ret < 0)
		{
			MYLOG_INFO(("write to peer err"));
			//this->del_and_close_fd(fd);
			return -1;
		}

		if(0 == ret)
		{
			MYLOG_INFO(("can't write now, try write ayn data_pos:%d data_len:%d", data_pos, data_len));

			memmove(&data[0], &data[data_pos], data_len);
			data.resize(data_len);

			this->data_tcp_out(fd, data);
			return 0;
		}

		data_len -= ret;
		data_pos += ret;
	}

	return 0;
}

/**
 * @brief 给指定的线程发消息
 */
int32 CMyThrdPoll::push_msg(void * msg, uint32 thrd_to, int32 thrd_from)
{
	MYLOG_DEBUG(("CMyThrdPoll::push_msg from:%d to:%d %x", thrd_from, thrd_to, msg));

	if(thrd_to >= this->thrds_.size())
	{
		MYLOG_WARN(("out of range %d - %d", thrd_to, this->thrds_.size()));
		return -1;
	}

	MYLOG_DEBUG(("malloc msg"));

	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(CMyThrdPoll::thrd_user_msg_t_));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("no memory"));
		return -1;
	}
	CMyThrdPoll::thrd_user_msg_t_ * m = new (place_ment)CMyThrdPoll::thrd_user_msg_t_;
	m->thrd_from_ = thrd_from;
	m->thrd_to_ = thrd_to;
	m->msg_body_ = msg;

	MYLOG_DEBUG(("real push msg"));

	MyListernerAddMsg(this->thrds_[thrd_to].thrd_, m, (unsigned long)this, CMyThrdPoll::_thrdsUserMsgCb);

	return 0;
}
/**
* @brief 根据fd来推消息
*/
int32 CMyThrdPoll::push_msg_by_fd(void * msg, int32 fd_to, int32 thrd_from)
{
	MYLOG_DEBUG(("CMyThrdPoll::push_msg_by fd_to:%d from:%d %x", fd_to, thrd_from, msg));

	assert(this->thrds_.size());
	if(0 == this->thrds_.size())
	{
		MYLOG_WARN(("not init yet ..."));
		return -1;
	}
	uint32 thrd_to = fd_to % this->thrds_.size();

	return this->push_msg(msg, thrd_to, thrd_from);
}

/**
 * @brief 添加定时器
 */
int32 CMyThrdPoll::add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period)
{
	MYLOG_DEBUG(("CMyThrdPoll::add_time time_second:%d thrd_index:%d timer_data:%x", time_second, thrd_index, timer_data));
	if(thrd_index >= this->thrds_.size())
	{
		MYLOG_WARN(("out of range %d - %d", thrd_index, this->thrds_.size()));
		return -1;
	}

	mytimer_node_t node = {0};

	node.context_data = (unsigned long)this->thrds_[thrd_index].thrd_context_data_;
	node.timer_user_data = timer_data;
	node.first_expire.tv_sec = time_second;
	/* 是否周期触发 */
	if(period)
		node.period.tv_sec = time_second;
	node.timeout_cb = CMyThrdPoll::_thrdsTimeOut;

	timer_id = MyListernerAddTimer(this->thrds_[thrd_index].thrd_, &node);

	return 0;
}

/**
 * @brief 删除定时器
 */
int32 CMyThrdPoll::del_time(uint32 thrd_index, HTIMERID timer_id)
{
	MYLOG_DEBUG(("CMyThrdPoll::del_time"));

	MyListernerDelTimer(this->thrds_[thrd_index].thrd_, timer_id);

	return 0;
}

/**
 * @brief 将一个udp fd加入监听
 */
int32 CMyThrdPoll::add_udp_fd(int32 udp_fd, uint32& thrd_index)
{
	MYLOG_DEBUG(("CMyThrdPoll::add_udp_fd %d", udp_fd));

	CChannel::set_no_block(udp_fd);

	this->_addToMapThrd(udp_fd, UDP_FD, thrd_index);

	CMyEPoll::event_handle e = {0};
	/* udp的特点就是直接扔,不管其它是不是处于可写状态,不需要做等到其可写时再发送 */
	e.input = CMyThrdPoll::_EpollInputUdp;
	e.context_data = this;
	this->epoll_->addfd(udp_fd, CMyEPoll::EVENT_INPUT, &e);

	return 0;
}

/**
 * @brief 将一个tcp fd加入监听
 */
int32 CMyThrdPoll::add_tcp_srv_fd(int32 tcp_fd, uint32& thrd_index, int32 bauto_add)
{
	MYLOG_DEBUG(("CMyThrdPoll::add_tcp_srv_fd %d", tcp_fd));

	CChannel::set_no_block(tcp_fd);

	if(0 != this->_addToMapThrd(tcp_fd, TCP_ACCEPTOR_FD, thrd_index))
	{
		MYLOG_WARN(("_addToMapThrd err"));
		return -1;
	}
	MYLOG_DEBUG(("_addToMapThrd end"));

	this->_set_accept_auto_add_to_listern_or_not(tcp_fd, bauto_add);

	CMyEPoll::event_handle e = {0};
	e.input = CMyThrdPoll::_EpollInputAccept;
	e.context_data = this;
	if(0 != this->epoll_->addfd(tcp_fd, CMyEPoll::EVENT_INPUT, &e))
	{
		MYLOG_WARN(("epoll_ addfd err"));
		return -1;
	}

	return 0;
}

/**
 * @brief 将一个tcp client fd 加入监听
 */
int32 CMyThrdPoll::add_tcp_cli_fd(int32 tcp_fd, uint32& thrd_index)
{
	MYLOG_DEBUG(("CMyThrdPoll::add_tcp_cli_fd %d", tcp_fd));

	CChannel::set_no_block(tcp_fd);

	if(0 != this->_addToMapThrd(tcp_fd, TCP_CONNECTOR_CLI_FD, thrd_index))
	{
		MYLOG_WARN(("_addToMapThrd err"));
		return -1;
	}
	MYLOG_DEBUG(("_addToMapThrd end"));

	CMyEPoll::event_handle e = {0};
	e.input = CMyThrdPoll::_EpollInputTcp;
	e.output = CMyThrdPoll::_EpollOutputTcp;
	e.exception = CMyThrdPoll::_EpollErr;
	e.context_data = this;
	if(0 != this->epoll_->addfd(tcp_fd, CMyEPoll::EVENT_INPUT | CMyEPoll::EVENT_OUTPUT, &e))
	{
		MYLOG_WARN(("epoll_ addfd err"));
		return -1;
	}

	return 0;
}

/**
 * @brief 取消对一个句柄的监听
 */
int32 CMyThrdPoll::del_and_close_fd(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::del_and_close_fd:%d", fd));

	HMYLISTERNER hlsn = NULL;
	if(0 != this->_getFdThrd(fd, hlsn))
	{
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	this->_push_inter_msg(hlsn, fd, E_DEL_AND_CLOSE_FD, NULL);

	return 0;
}

/**
 * @brief process tcp input,this function was invoke in the CMyThrdPoll::work_loop context
 */
int32 CMyThrdPoll::_EpollInputTcp(void * context_data, int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_InputTcp fd:%d", fd));

	/* find fd info,and send msg */
	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp)
	{
		MYLOG_WARN(("context data is null"));
		return -1;
	}

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != tp->_getFdThrd(fd, hlsn))
	{
		//tp->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	tp->_push_inter_msg(hlsn, fd, E_INPUT_TCP, NULL);

	MYLOG_DEBUG(("push msg end:%d\r\n", fd));

	return 0;
}

/**
 * @brief process accept event,this function was invoke in the CMyThrdPoll::work_loop context
 */
int32 CMyThrdPoll::_EpollInputAccept(void * context_data, int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_InputAccept fd:%d \r\n", fd));

	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp)
	{
		MYLOG_WARN(("context data is null"));
		return -1;
	}

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != tp->_getFdThrd(fd, hlsn))
	{
		//tp->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	tp->_push_inter_msg(hlsn, fd, E_ACCEPT_TCP, NULL);

	return 0;
}

/**
 * @brief process output,this function was invoke in the CMyThrdPoll::work_loop context
 */
int32 CMyThrdPoll::_EpollOutputTcp(void * context_data, int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_Output fd:%d \r\n", fd));

	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp)
	{
		MYLOG_WARN(("context data is null"));
		return -1;
	}

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != tp->_getFdThrd(fd, hlsn))
	{
		//tp->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	tp->_push_inter_msg(hlsn, fd, E_OUTPUT_TCP, NULL);

	return 0;
}

/**
 * @brief process udp input,在work_loop所在的线程中被呼叫
 */
int32 CMyThrdPoll::_EpollInputUdp(void * context_data, int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_EpollInputUdp %d", fd));

	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp)
	{
		MYLOG_WARN(("context data is null"));
		return -1;
	}

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != tp->_getFdThrd(fd, hlsn))
	{
		//tp->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	tp->_push_inter_msg(hlsn, fd, E_INPUT_UDP, NULL);

	return 0;
}

/**
 * @brief process err 在work_loop所在的线程中被呼叫
 */
int32 CMyThrdPoll::_EpollErr(void * context_data, int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_EpollErr fd:%d \r\n", fd));

	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp)
	{
		MYLOG_WARN(("context data is null"));
		return -1;
	}

	/* hlsn是不会被销毁的,可以在不加锁的情况下使用 */
	HMYLISTERNER hlsn = NULL;
	if(0 != tp->_getFdThrd(fd, hlsn))
	{
		//tp->epoll_->delfd(fd);
		MYLOG_WARN(("can find fd info\r\n"));
		return -1;
	}
	if(NULL == hlsn)
		return -1;

	tp->_push_inter_msg(hlsn, fd, E_ERR, NULL);

	return 0;
}


/**
 * @brief 处理消息队列里的消息,this function was invoke in the multi thrds context
 */
int32 CMyThrdPoll::_thrdsMsgCb(unsigned long context_data, void * msg)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsMsgCb msg:%x", msg));
	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp || NULL == msg)
	{
		MYLOG_WARN(("tp is null or msg is null"));
		return -1;
	}

	CMyThrdPoll::thrd_msg_t * e = (CMyThrdPoll::thrd_msg_t *)msg;
	if(NULL == e)
	{
		MYLOG_WARN(("msg is null"));
		return -1;
	}

	MYLOG_DEBUG(("event:%d fd:%d", e->event_, e->fd_));

	switch(e->event_)
	{
	case E_ERR:
		tp->_thrdsDoErr(e->fd_);
		break;

	case E_INPUT_TCP:
		tp->_thrdsDoInputTcp(e->fd_);
		break;

	case E_OUTPUT_TCP:
		tp->_thrdsDoOutputTcp(e->fd_);
		break;

	case E_ACCEPT_TCP:
		tp->_thrdsDoAcceptTcp(e->fd_);
		break;

	case E_USER_OUTPUT_TCP:
		tp->_thrdsDoUserOutputTcp(e->fd_, e);
		break;

	case E_INPUT_UDP:
		tp->_thrdsDoInputUdp(e->fd_);
		break;

	case E_NEW_CONN:
		tp->_thrdsDoNewConn(e->fd_);
		break;

	case E_DEL_AND_CLOSE_FD:
		tp->_thrdsDoDelAndClose(e->fd_);
		break;

	default:
		MYLOG_WARN(("we have unkown event..."));
		break;
	}

	e->free_body(tp->hm_);
	MyMemPoolFree(tp->hm_, e);

	MYLOG_DEBUG(("_thrdsMsgCb end\r\n\r\n"));

	return 0;
}

/**
 * @brief 用户消息回调函数
 */
int32 CMyThrdPoll::_thrdsUserMsgCb(unsigned long context_data, void * msg)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsUserMsgCb msg:%x", msg));
	CMyThrdPoll * tp = (CMyThrdPoll *)context_data;
	if(NULL == tp || NULL == msg)
	{
		MYLOG_WARN(("tp is null or msg is null"));
		return -1;
	}

	CMyThrdPoll::thrd_user_msg_t_ * m = (CMyThrdPoll::thrd_user_msg_t_ *)msg;
	if(NULL == m)
	{
		MYLOG_DEBUG(("msg is null"));
		return -1;
	}

	/* 消息要入口时被限制了,一定会有线程处理 */
	assert(m->thrd_to_ < tp->thrds_.size());

	MYLOG_DEBUG(("m->msg_body_:%x, m->thrd_from_:%d, m->thrd_to_:%d", m->msg_body_, m->thrd_from_, m->thrd_to_));
	tp->thrds_[m->thrd_to_].thrd_context_data_->msg_callback(m->thrd_from_, m->msg_body_);

	/* 交由产生m->msg_body_的代码来释放m->msg_body_,此处不释放m->msg_body_ */
	m->msg_body_ = NULL;

	/* 没有析构,无需要呼叫 */
	//m->~thrd_user_msg_t_();
	MyMemPoolFree(tp->hm_, m);
	return 0;
}

/**
 * @brief 超时回调
 */
int32 CMyThrdPoll::_thrdsTimeOut(unsigned long context_data, unsigned long timer_user_data, HTIMERID timerid)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsTimeOut"));

	ThrdHandle * thrd_context_data = (ThrdHandle *)context_data;
	if(NULL == thrd_context_data)
	{
		MYLOG_WARN(("err context data"));
		return -1;
	}

	thrd_context_data->time_out(timer_user_data, timerid);

	return -1;
}

/**
 * @brief do err, 在线程池中某个线程的上下文中被呼叫
 */
int32 CMyThrdPoll::_thrdsDoErr(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoErr %d", fd));

	this->_delFromMapThrd(fd);

	return 0;
}

/**
 * @brief 删除链接
 */
int32 CMyThrdPoll::_thrdsDoDelAndClose(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoDelAndClose %d", fd));

	this->_delFromMapThrd(fd);

	return 0;
}

/**
 * @brief accept new connection,this function was invoke in the thrds context
 */
int32 CMyThrdPoll::_thrdsDoAcceptTcp(int32 fd)
{
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
		CChannel::set_no_linger(accept_fd);

		MYLOG_DEBUG(("accept new connection %d", accept_fd));

		fd_info_t * pfi = NULL;
		/* pfi的释放保证了是在同一个线程上下文当中的,所以可以在不加锁的情况下使用 */
		this->_getFdInfoFromMapThrd(fd, pfi);
		if(NULL == pfi)
		{
			MYLOG_WARN(("accept end fd info is null"));
			continue;
		}

		if(pfi->bauto_add_to_listern_)
		{
			MYLOG_DEBUG(("need auto add to listern"));

			uint32 thrd_index = 0;
			this->_addToMapThrd(accept_fd, TCP_CONNECTOR_FD, thrd_index, fd);
			MYLOG_DEBUG(("_addToMapThrd end"));

			CMyEPoll::event_handle e = {0};
			e.input = CMyThrdPoll::_EpollInputTcp;
			e.output = CMyThrdPoll::_EpollOutputTcp;
			e.exception = CMyThrdPoll::_EpollErr;
			e.context_data = this;
			this->epoll_->addfd(accept_fd, CMyEPoll::EVENT_INPUT | CMyEPoll::EVENT_OUTPUT, &e);

			/* 推入消息
			*/
			this->_push_inter_msg(this->thrds_[thrd_index].thrd_, accept_fd, E_NEW_CONN, NULL);
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
 * @brief do input,this function was invoke in the thrds context
 */
int32 CMyThrdPoll::_thrdsDoInputTcp(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoInput %d ==================", fd));
	
	/* 由于fd对应的fi是会在同一个线程的上下文中被删除,被访问,所以可以不加锁
	 * read and,when egain,stop
	 */
	CMyThrdPoll::fd_info_t * pfi = NULL;
	if(0 != this->_getFdInfoFromMapThrd(fd, pfi))
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

	MYLOG_DEBUG(("fi buf size:%d pos:%d", pfi->rbuf_.size(), pfi->rpos_));

	int32 need_close = 0;
	do
	{
		if(pfi->rpos_ >= pfi->rbuf_.size())
		{
			pfi->rbuf_.resize(pfi->rbuf_.size() + 1024);
		}
		
		int32 ret = CChannel::TcpRead(pfi->fd_, &pfi->rbuf_[pfi->rpos_], pfi->rbuf_.size() - pfi->rpos_);

		MYLOG_DEBUG(("read ret:%d", ret));

		if(ret < 0)
		{
			MYLOG_DEBUG(("err need close"));
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
			pfi->rpos_ += ret;

			pfi->byte_read_ += ret;

			MYLOG_DUMP_BIN(&pfi->rbuf_[0], pfi->rpos_);

		}
	}while(1);

	uint32 recved = 0;
	MYLOG_DEBUG(("fd:%d rpos:%d rrpos:%d %x size:%d", pfi->fd_, pfi->rpos_, pfi->rrpos_, &pfi->rbuf_[0], pfi->rbuf_.size()));

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

	if(pfi->thrd_context_data_ && pfi->rpos_ > pfi->rrpos_)
		pfi->thrd_context_data_->data_tcp_in(pfi->fd_, &pfi->rbuf_[pfi->rrpos_], pfi->rpos_ - pfi->rrpos_, recved);

	MYLOG_DEBUG(("_data_in end recved:%d rrpos_:%d rpos_:%d bufsz_reserve_:%d", recved, pfi->rrpos_, pfi->rpos_, this->bufsz_reserve_));
	pfi->rrpos_ += recved;
	if(pfi->rrpos_ >= pfi->rpos_)
	{
		MYLOG_DEBUG(("all data in buf has been recved by app, now pos go back"));
		pfi->rrpos_ = 0;
		pfi->rpos_ = 0;

		if(pfi->rbuf_.capacity() > MY_MTU)
		{
			pfi->rbuf_.clear();
			pfi->rbuf_.resize(MY_MTU);
		}
	}
	else if(((pfi->rpos_ - pfi->rrpos_) < this->bufsz_reserve_) && (pfi->rrpos_ > this->bufsz_reserve_))
	{
		/* 未被上层接收的缓冲小于指定大小bufsz_reserve_,并且pfi->rrpos_位置大于bufsz_reserve_,将缓冲区整到开头去 */
		MYLOG_DEBUG(("all reserve data copy to the front of the buf, and pos go back"));

		/* move buf_sz back */
		memmove(&pfi->rbuf_[0], &pfi->rbuf_[pfi->rrpos_], pfi->rpos_ - pfi->rrpos_);
		pfi->rpos_ -= pfi->rrpos_;
		pfi->rrpos_ = 0;
	}
	else
	{
		MYLOG_DEBUG(("not all been recv,but don't move buf"));
	}

	if(need_close)
	{
		MYLOG_DEBUG(("fd need to close"));
		
		this->_delFromMapThrd(fd);
	}

	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoInput %d end =================", fd));

	return 0;
}

/**
 * @brief do output
 */
int32 CMyThrdPoll::_thrdsDoOutputTcp(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoOutput:%d", fd));

	/* 先发sbuf_,然后再发lst_sbuf_
	* 如果送不动,则返回,等下一次可写事件发生时,再发送
	*/

	CMyThrdPoll::fd_info_t * pfi = NULL;
	this->_getFdInfoFromMapThrd(fd, pfi);
	if(NULL == pfi)
		return -1;

	MYLOG_DEBUG(("lst_sbuf_:%d", pfi->lst_sbuf_.size()));

	while(pfi->sbuf_.size())
	{
		MYLOG_DEBUG(("lst_sbuf_:%d", pfi->lst_sbuf_.size()));

		int32 data_len = pfi->sbuf_.size();
		uint32 data_pos = 0;

		/* 循环发送
		* 送完为止
		*/
		while(data_len > 0)
		{
			int32 ret = 0;
			if(data_len < MY_MTU)
			{
				MYLOG_DEBUG(("data_pos:%d write data:%d", data_pos, data_len));
				ret = CChannel::TcpWrite(pfi->fd_, &pfi->sbuf_[data_pos], data_len);
			}
			else
			{
				MYLOG_DEBUG(("data_pos:%d write data:%d", data_pos, MY_MTU));
				ret = CChannel::TcpWrite(pfi->fd_, &pfi->sbuf_[data_pos], MY_MTU);
			}
			MYLOG_DEBUG(("write data ret:%d", ret));

			if(ret < 0)
			{
				MYLOG_INFO(("write err, fd:%d need close", fd));
				//this->del_and_close_fd(fd);
				return 0;
			}
			if(0 == ret)
			{
				MYLOG_INFO(("can't write any more need wait data_len:%d lst_sbuf_:%d fd:%d", data_len, pfi->lst_sbuf_.size(), pfi->fd_));
				memmove(&pfi->sbuf_[0], &pfi->sbuf_[data_pos], data_len);
				pfi->sbuf_.resize(data_len);

				return 0;
			}

			pfi->byte_write_ += ret;

			data_len -= ret;
			data_pos += ret;
		}

		/* move unsend data to buf head */
		assert(data_len <= 0);

		MYLOG_DEBUG(("current buf have been done"));
		pfi->sbuf_.clear();

		if(pfi->lst_sbuf_.empty())
		{
			MYLOG_DEBUG(("all data has been send"));
			//if(0 != this->epoll_->modfd(fd, CMyEPoll::EVENT_INPUT /*| CMyEPoll::EVENT_OUTPUT*/))
			//	return -1;
			break;
		}

		MYLOG_DEBUG(("still have buf:%d", pfi->lst_sbuf_.size()));

		assert(pfi->lst_sbuf_.size());
		pfi->sbuf_.swap(*pfi->lst_sbuf_.begin());
		assert(pfi->sbuf_.size());
		//todo 这里可以再次做累积,直到pfi->sbuf_的大小为MY_MTU

		pfi->lst_sbuf_.pop_front();
	}

	MYLOG_DEBUG(("write data done:%d", pfi->fd_));
	return 0;
}

/**
 * @brief do user output
 */
int32 CMyThrdPoll::_thrdsDoUserOutputTcp(int32 fd, CMyThrdPoll::thrd_msg_t * msg)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoUserOutput fd:%d msg:%x", fd, msg));

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

	CMyThrdPoll::fd_info_t * pfi = NULL;
	this->_getFdInfoFromMapThrd(fd, pfi);
	if(NULL == pfi)
		return -1;

#ifdef WIN32
	/* why vc can't receive writable event... */
	CChannel::TcpSelectWrite(fd, &mb->msg[0], mb->msg.size());
#else

	MYLOG_DEBUG(("current buf:size:%d send:%d lst_sbuf_:%d",
		pfi->sbuf_.size(), mb->msg.size(), pfi->lst_sbuf_.size()));

	int32 byte_left = mb->msg.size();
	uint32 pos = 0;

	if(pfi->lst_sbuf_.empty())
	{
		MYLOG_DEBUG(("lst_sbuf_ is empty"));

		if(0 == pfi->sbuf_.size())
		{
			MYLOG_DEBUG(("send buf has nothing,swap it"));

			if(mb->msg.size() <= MY_MTU)
			{
				MYLOG_DEBUG(("send buf less the mtu, swap it"));
				pfi->sbuf_.swap(mb->msg);/* 大部分情况,代码会走这个路径 */
				byte_left = 0;
			}
			else
			{
				MYLOG_DEBUG(("send buf more the mtu, append it"));
				pfi->sbuf_.insert(pfi->sbuf_.end(), mb->msg.begin(), mb->msg.begin() + MY_MTU);
				byte_left -= MY_MTU;
				pos += MY_MTU;
			}
		}
		else
		{
			MYLOG_DEBUG(("send buf has something,append it"));
			if(mb->msg.size() <= MY_MTU)
			{
				MYLOG_DEBUG(("send buf less the mtu, append it"));
				pfi->sbuf_.insert(pfi->sbuf_.end(), mb->msg.begin(), mb->msg.end());
				byte_left = 0;
			}
			else
			{
				MYLOG_DEBUG(("send buf more the mtu, append it"));
				pfi->sbuf_.insert(pfi->sbuf_.end(), mb->msg.begin(), mb->msg.begin() + MY_MTU);
				byte_left -= MY_MTU;
				pos += MY_MTU;
			}
		}
	}

	while(byte_left > 0)
	{
		MYLOG_DEBUG(("byte_left:%d", byte_left));

		/* todo:we need option here, */
		std::vector<uint8> vbuf(2 * MY_MTU);
		vbuf.resize(0);

		if((mb->msg.size() - pos) <= MY_MTU)
		{
			MYLOG_DEBUG(("left send buf less than mtu"));
			vbuf.insert(vbuf.end(), mb->msg.begin() + pos, mb->msg.end());

			pfi->lst_sbuf_.push_back(vbuf);
			break;
		}
		else
		{
			MYLOG_DEBUG(("left send buf more than mtu"));
			vbuf.insert(vbuf.end(), mb->msg.begin() + pos, mb->msg.begin() + pos + MY_MTU);
			byte_left -= MY_MTU;
			pos += MY_MTU;

			pfi->lst_sbuf_.push_back(vbuf);
			continue;
		}
	}

	this->_thrdsDoOutputTcp(fd);

	//if(0 != this->epoll_->modfd(fd, /*CMyEPoll::EVENT_INPUT |*/ CMyEPoll::EVENT_OUTPUT))
	//	return -1;
#endif

	return 0;
}

/**
 * @brief do new connection in event, 在线程池中某个线程的上下文中被呼叫
 */
int32 CMyThrdPoll::_thrdsDoNewConn(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoNewConn %d", fd));

	CMyThrdPoll::fd_info_t * pfi = NULL;
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
		if(pfi->rpos_ > pfi->rrpos_)
		{
			MYLOG_DEBUG(("before connect cb,there has been data in this channel already,process it"));

			uint32 recved = 0;
			pfi->thrd_context_data_->data_tcp_in(fd, &pfi->rbuf_[pfi->rrpos_], pfi->rpos_ - pfi->rrpos_, recved);
			pfi->rrpos_ += recved;
		}
	}
	else
		MYLOG_INFO(("new connection in,but context data is null, so no call back"));
	return 0;
}

/**
 * @brief do udp input
 */
int32 CMyThrdPoll::_thrdsDoInputUdp(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_thrdsDoInputUdp %d", fd));

	/* 由于fd对应的fi是会在同一个线程的上下文中被删除,被访问,所以可以不加锁
	 * read and,when egain,stop
	 */
	CMyThrdPoll::fd_info_t * pfi = NULL;
	if(0 != this->_getFdInfoFromMapThrd(fd, pfi))
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

	MYLOG_DEBUG(("fi buf size:%d pos:%d", pfi->rbuf_.size(), pfi->rpos_));

	/* 读取报文 */
	/* 回调 */
	do
	{
		int8 acip[32] = {0};
		uint16 port = 0;
		int32 ret = CChannel::UdpRead(fd, &pfi->rbuf_[0], pfi->rbuf_.size(), acip, sizeof(acip), &port);

		MYLOG_DEBUG(("%s:%d ret:%d", acip, port, ret));
		MYLOG_DUMP_BIN(&pfi->rbuf_[0], ret);

		if(ret <= 0)
			break;

		pfi->thrd_context_data_->data_udp_in(fd, &pfi->rbuf_[0], ret, acip, port);

	}while(1);

	return 0;
}


/**
* @brief 添加至 map_thrd_
*/
int32 CMyThrdPoll::_addToMapThrd(int32 fd, SOCKET_FD_TYPE_T fd_type, uint32& thrd_index, int32 fd_master)
{
	MYLOG_DEBUG(("CMyThrdPoll::_addToMapThrd %d fd_type:%d", fd, fd_type));

	{
		CAutoLock alock(&this->map_thrd_protector_);

		std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.find(fd);
		if(this->map_thrd_.end() != it)
		{
			MYLOG_WARN(("fd:%d is in map", fd));
			return -1;
		}

		thrd_index = fd % this->thrds_.size();

		fd_info_t * i = new fd_info_t;
		i->fd_ = fd;
		i->type_ = fd_type;
		/* fd固定由一个线程处理,
		* 对fd对应的事件时,都添加进此消息队列 
		*/
		i->hlsn_ = this->thrds_[thrd_index].thrd_;
		i->thrd_context_data_ = this->thrds_[thrd_index].thrd_context_data_;
		i->rpos_ = 0;
		i->rrpos_ = 0;
		i->rbuf_.resize(MY_MTU);

		/* 填充供分析日志时使用的杂项信息
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
	}

	return 0;
}

/**
 * @brief 删除
 */
int32 CMyThrdPoll::_delFromMapThrd(int32 fd)
{
	MYLOG_DEBUG(("CMyThrdPoll::_delFromMapThrd %d", fd));

	ThrdHandle * h = NULL;

	this->epoll_->delfd(fd);

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

		MYLOG_INFO(("fd:%d local[%s:%d] remote[%s:%d] byte_write:%u",
			fd,
			it->second->local_ip_.c_str(), it->second->local_port_,
			it->second->remote_ip_.c_str(), it->second->remote_port_,
			it->second->byte_write_));

		//CChannel::CloseFd(it->second->fd_);

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
 * @brief 获取fd的上下文数据,由于上下文数据只会被固定的一个线程所操作,所以在操作上下文数据时不加锁
 */
int32 CMyThrdPoll::_getFdInfoFromMapThrd(int32 fd, fd_info_t*& i)
{
	MYLOG_DEBUG(("CMyThrdPoll::_getFdInfoFromMapThrd %d", fd));

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
int32 CMyThrdPoll::_getFdThrd(int32 fd, HMYLISTERNER& hlsn)
{
	MYLOG_DEBUG(("CMyThrdPoll::_getFdThrd %d", fd));

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
* @brief 设置成自动加入或非自动加入监听
*/
int32 CMyThrdPoll::_set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add)
{
	MYLOG_DEBUG(("CMyThrdPoll::_set_accept_auto_add_to_listern_or_not"));

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
int32 CMyThrdPoll::_push_inter_msg(HMYLISTERNER hlsn, int32 fd, uint32 ev, void * pb)
{
	MYLOG_DEBUG(("CMyThrdPoll::_push_inter_msg %x fd:%d event:%d pb:%x", hlsn, fd, ev, pb));

	/* 推入消息
	*/
	void * place_ment = MyMemPoolMalloc(this->hm_, sizeof(CMyThrdPoll::thrd_msg_t));
	if(NULL == place_ment)
	{
		MYLOG_ERR(("no memory"));
		return -1;
	}
	CMyThrdPoll::thrd_msg_t * msg = new(place_ment)CMyThrdPoll::thrd_msg_t;

	if(NULL == msg)
	{
		MYLOG_ERR(("impossible err"));
		return -1;
	}

	msg->fd_ = fd;
	msg->event_ = ev;
	msg->msg_body_ = pb;
	MyListernerAddMsg(hlsn, msg, (unsigned long)this, CMyThrdPoll::_thrdsMsgCb);

	return 0;
}

/**
* @brief close all open tcp connect(listerner tcp sock and udp don't close)
*/
void CMyThrdPoll::close_tcp_connect()
{
	MYLOG_ERR(("CMyThrdPoll::close_tcp_connect"));

	CAutoLock alock(&this->map_thrd_protector_);

	for(std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.begin(); it != this->map_thrd_.end(); it ++)
	{
		if(NULL == it->second)
			continue;

		switch(it->second->type_)
		{
		case TCP_CONNECTOR_FD:
		case TCP_CONNECTOR_CLI_FD:
			{
				MYLOG_ERR(("close fd:%d", it->second->fd_));

				this->_push_inter_msg(it->second->hlsn_, it->second->fd_, E_DEL_AND_CLOSE_FD, NULL);
			}
			break;
		}
	}
}

/**
 * @brief runtime_view 多线程不安全函数,仅供调试时使用
 */
void CMyThrdPoll::runtime_view()
{
	this->epoll_->runtime_view();

	{
		CAutoLock alock(&this->map_thrd_protector_);
		for(std::map<int32, fd_info_t *>::iterator it = this->map_thrd_.begin();
			it != this->map_thrd_.end();
			it ++)
		{
			std::string peer_ip;
			uint32 peer_port;
			CChannel::getPeerName(it->second->fd_, peer_ip, peer_port);
			MYLOG_ERREX(("view", "sock:%d-%d type_:%d remote:[%s:%d]-[%s:%d] local[%s:%d] hlsn_:%x thrd_context_data_:%x rpos_:%d rrpos_:%d rbuf_:%d sbuf_:%d lst_sbuf_:%d byte_write:%u byte_read:%u",
				it->first,
				it->second->fd_,
				it->second->type_,
				it->second->remote_ip_.c_str(), it->second->remote_port_,
				peer_ip.c_str(), peer_port,
				it->second->local_ip_.c_str(), it->second->local_port_,
				it->second->hlsn_,
				it->second->thrd_context_data_,
				it->second->rpos_,
				it->second->rrpos_,
				it->second->rbuf_.size(),
				it->second->sbuf_.size(),
				it->second->lst_sbuf_.size(),
				it->second->byte_write_,
				it->second->byte_read_));

			/* 打出接收缓冲区里的内容 */
			MYLOG_ERREX(("view", "recv buf"));

			MYLOG_DUMP_BIN(&it->second->rbuf_[it->second->rrpos_], it->second->rpos_ - it->second->rrpos_);
			MYLOG_DUMP_BIN(&it->second->sbuf_[0], it->second->sbuf_.size());
		}
		MYLOG_ERREX(("view", "map_thrd_.size:%d", this->map_thrd_.size()));
	}

	for(uint32 i = 0; i < this->thrds_.size(); i ++)
	{
		char actemp[256] = {0};
		MyListernerPrint(this->thrds_[i].thrd_, actemp, sizeof(actemp) - 1);
		MYLOG_ERREX(("view", actemp));
	}
	MYLOG_ERREX(("view", "END\r\n\r\n"));
}


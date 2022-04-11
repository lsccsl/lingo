/**
* @file tcpclient_pool.cpp
* @brief tcp client connect pool mgr module
* @author linsc
*/
#include "tcpclient_pool.h"
#include <assert.h>
#include <vector>
#include "mylogex.h"
#include "channel.h"

/**
* @brief constructor
* @param srv_ip:tcp服务端的ip
* @param srv_port:服务端port
*/
tcpclient_pool::tcpclient_pool(const int8 * srv_ip, const uint32 srv_port, const uint32 max_conn_count, const uint32 min_conn_count):
	srv_ip_(srv_ip ? srv_ip : ""),srv_port_(srv_port),max_conn_count_(max_conn_count),min_conn_count_(min_conn_count)
{
	MYLOG_DEBUG(("tcpclient_pool::tcpclient_pool max_conn_count_:%u min_conn_count_:%u [%s:%u]",
		this->max_conn_count_, this->min_conn_count_, this->srv_ip_.c_str(), this->srv_port_));

	if(this->min_conn_count_ > this->max_conn_count_)
		this->min_conn_count_ = this->max_conn_count_;

	assert(srv_ip_.size());
	assert(srv_port <= 65535);

	pthread_mutex_init(&this->fd_protector_, NULL);
	sem_init(&this->sem_fd_, 0, 0);
}

/**
* @brief destructor
*/
tcpclient_pool::~tcpclient_pool()
{
	MYLOG_DEBUG(("tcpclient_pool::~tcpclient_pool"));

	this->max_conn_count_ = 0;
	this->min_conn_count_ = 0;

	sem_post(&this->sem_fd_);
	pthread_mutex_lock(&this->fd_protector_);
	this->lst_idle_fd_.clear();

	for(std::set<int32>::iterator it = this->s_fd_.begin(); it != this->s_fd_.end(); it ++)
	{
		MYLOG_DEBUG(("close fd:%d", *it));
		CChannel::CloseFd(*it);
	}
	this->s_fd_.clear();

	pthread_mutex_unlock(&this->fd_protector_);

	pthread_mutex_destroy(&this->fd_protector_);
	sem_destroy(&this->sem_fd_);
}

/**
* @brief send data
*/
int32 tcpclient_pool::send_data(const void * buf, uint32 buf_sz, const uint32 count, const uint32 timeout_second,
	const tcpclient_send_end_cb * cb, const void * context_data, const void * user_data)
{
	MYLOG_DEBUG(("tcpclient_pool::send_data buf:%x buf_sz:%u count:%u timeout_second:%u cb:%x context_data:%x user_data:%x", buf, buf_sz, count, timeout_second,
		cb, context_data, user_data));

	int32 fd = -1;
	for(uint32 i = 0; i < 2; i ++)
	{
		/* 获取连接 */
		if(0 != this->_get_conn(fd))
		{
			MYLOG_INFO(("get connect fail ..."));
			return -1;
		}

		/* 发送数据 */
		int32 ret = CChannel::TcpSelectWrite(fd, buf, buf_sz, timeout_second, count);

		MYLOG_DEBUG(("ret:%d", ret));

		if(ret != buf_sz)
		{
			MYLOG_INFO(("fail"));

			pthread_mutex_lock(&this->fd_protector_);
			this->s_fd_.erase(fd);
			CChannel::CloseFd(fd);
			pthread_mutex_unlock(&this->fd_protector_);

			fd = -1;
			continue;
		}
		else
		{
			MYLOG_DEBUG(("write:%d suc", fd));
			break;
		}
	}

	MYLOG_DEBUG(("fd:%d", fd));

	if(fd < 0)
	{
		MYLOG_INFO(("fd:%d", fd));
		return -1;
	}

	/* 呼叫回调 */
	if(cb)
	{
		MYLOG_DEBUG(("need call back"));
		if(0 != cb->send_end_call_back(fd, context_data, user_data))
		{
			MYLOG_DEBUG(("call back end, need close fd, don't recyc"));

			pthread_mutex_lock(&this->fd_protector_);
			this->s_fd_.erase(fd);
			CChannel::CloseFd(fd);
			pthread_mutex_unlock(&this->fd_protector_);

			return -1;
		}
	}

	MYLOG_DEBUG(("call back end, recyc fd"));
	this->_recyc_fd(fd);

	return 0;
}

/**
* @brief read data
*/
int32 tcpclient_pool::read_data(int32 fd, const void * buf, uint32 buf_sz, uint32 count, uint32 timeout_second)
{
	MYLOG_DEBUG(("tcpclient_pool::read_data fd:%d buf:%x buf_sz:%u count:%u timeout_second:%u", fd, buf, buf_sz, count, timeout_second));

	int32 ret = CChannel::TcpSelectRead(fd, (void *)buf, buf_sz, timeout_second, count);
	if(ret < 0)
		return -1;
	if(ret != buf_sz)
		return -1;

	return 0;
}

/**
* @brief 取出或者生成一个tcp链接
*/
int32 tcpclient_pool::_get_conn(int32& fd)
{
	MYLOG_DEBUG(("tcpclient_pool::_get_conn"));

	fd = -1;
	int32 ret = 0;

	while(1)
	{
		MYLOG_DEBUG(("loop"));

		/* 有无空闲连接 */
		pthread_mutex_lock(&this->fd_protector_);
		if(this->lst_idle_fd_.empty())
		{
			MYLOG_DEBUG(("no idle connect %d", this->s_fd_.size()));

			/* 无,则看是否达到上限 */
			if(this->s_fd_.size() >= this->max_conn_count_)
			{
				MYLOG_DEBUG(("reach max %d", this->max_conn_count_));

				/* 达到上限则等待返回空闲连接 */
				pthread_mutex_unlock(&this->fd_protector_);

				MYLOG_DEBUG(("wait for idle"));
				sem_wait(&this->sem_fd_);
				continue;
			}
			else
			{
				MYLOG_DEBUG(("not reach max %d", this->max_conn_count_));

				pthread_mutex_unlock(&this->fd_protector_);

				MYLOG_DEBUG(("connect to [%s:%u]", this->srv_ip_.c_str(), this->srv_port_));
				/* 生成新链接 */
				fd = CChannel::TcpConnect(this->srv_ip_.c_str(), this->srv_port_, 0, 0);
				if(fd < 0)
				{
					MYLOG_DEBUG(("connect to [%s:%u] fail ... ", this->srv_ip_.c_str(), this->srv_port_));
					ret = -1;
					break;
				}

				/* 不置成阻塞,断开时,发送仍成功 */
				/*CChannel::set_no_block(fd);*/

				MYLOG_DEBUG(("connect to [%s:%u] ok", this->srv_ip_.c_str(), this->srv_port_));

				pthread_mutex_lock(&this->fd_protector_);
				this->s_fd_.insert(fd);
				pthread_mutex_unlock(&this->fd_protector_);

				break;
			}
		}
		else
		{
			MYLOG_DEBUG(("have idle connect"));

			/* 有则取之 */
			fd = *(this->lst_idle_fd_.begin());
			this->lst_idle_fd_.pop_front();
			pthread_mutex_unlock(&this->fd_protector_);
			break;
		}
	}

	MYLOG_DEBUG(("fd:%d ret:%d", fd, ret));

	return ret;
}

/**
* @brief 回取tcp链接
*/
int32 tcpclient_pool::_recyc_fd(int32 fd)
{
	MYLOG_DEBUG(("tcpclient_pool::_recyc_fd fd:%d", fd));

	pthread_mutex_lock(&this->fd_protector_);
	this->lst_idle_fd_.push_front(fd);

	if(1 == this->lst_idle_fd_.size())
	{
		MYLOG_DEBUG(("need wake up"));
		sem_post(&this->sem_fd_);
	}

	pthread_mutex_unlock(&this->fd_protector_);

	return 0;
}

/**
* @brief for debug
*/
void tcpclient_pool::view()
{
	MYLOG_INFOEX(("view", "tcp client pool view begin ================================"));

	MYLOG_INFOEX(("view", "[%s:%u] max_conn_count_:%u min_conn_count_:%u fd count:%d idle fd:%d", this->srv_ip_.c_str(), this->srv_port_,
		this->max_conn_count_, this->min_conn_count_, this->s_fd_.size(), this->lst_idle_fd_.size()));

	pthread_mutex_lock(&this->fd_protector_);
	
	{
		for(std::set<int32>::iterator it = this->s_fd_.begin(); it != this->s_fd_.end(); it ++)
		{
			std::string peer_ip;
			uint32 peer_port;
			CChannel::getSocketName((*it), peer_ip, peer_port);
			MYLOG_DEBUGEX(("view", "fd:%d [%s:%u]", (*it), peer_ip.c_str(), peer_port));
		}
	}

	{
		for(std::list<int32>::iterator it = this->lst_idle_fd_.begin(); it != this->lst_idle_fd_.end(); it ++)
		{
			MYLOG_DEBUGEX(("view", "fd:%d", (*it)));
		}
	}

	pthread_mutex_unlock(&this->fd_protector_);

	MYLOG_INFOEX(("view", "tcp client pool view begin ================================\r\n"));
}











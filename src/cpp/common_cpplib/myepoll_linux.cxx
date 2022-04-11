/**
 * @file myepoll.cpp
 * @brief wrapper event poll
 *
 * @author linshaochuan
 */
#include "myepoll.h"
#ifndef WIN32
#include <sys/epoll.h>
#include <unistd.h>
#endif
#include <errno.h>
#include "mylogex.h"
#include "AutoLock.h"

/**
 * @brief 构造
 */
CMyEPoll::CMyEPoll(uint32 max_fd_count, int32 wait_thrd_count)
{
	MYLOG_INFO(("wait_thrd_count:%d",
		wait_thrd_count));

	for(int32 i = 0; i < wait_thrd_count - 1; i ++)
	{
		epoll_wait_thrd * ewt = new epoll_wait_thrd(max_fd_count);
		if(0 == ewt->work())
		{
			MYLOG_INFO(("epoll wait thrd:%d", i));
			this->v_thrds_.push_back(ewt);
		}
		else
		{
			MYLOG_WARN(("epoll wait thrd fail:%d", i));
			delete ewt;
		}
	}

	this->ewt_ = new epoll_wait_thrd(max_fd_count);
	this->v_thrds_.push_back(this->ewt_);

	MYLOG_INFO(("v_thrds_:%d", this->v_thrds_.size()));
}

/**
 * @brief 析构
 */
CMyEPoll::~CMyEPoll()
{
	MYLOG_DEBUG(("CMyEPoll::~CMyEPoll"));
	for(uint32 i = 0; i < this->v_thrds_.size(); i ++)
	{
		MYLOG_DEBUG(("stop:%d", i));

		if(i < (this->v_thrds_.size() - 1))
			this->v_thrds_[i]->stop();

		delete this->v_thrds_[i];
	}
}

/**
 * @brief add fd
 */
int32 CMyEPoll::addfd(int32 fd, uint64 event_mask, event_handle * eh)
{
	MYLOG_DEBUG(("CMyEPoll::addfd"));

	uint32 idx = fd % this->v_thrds_.size();

	return this->v_thrds_[idx]->addfd(fd, event_mask, eh);
}

/**
 * @brief del fd
 */
int32 CMyEPoll::delfd(int32 fd)
{
	MYLOG_DEBUG(("CMyEPoll::delfd"));

	uint32 idx = fd % this->v_thrds_.size();

	return this->v_thrds_[idx]->delfd(fd);
}

/**
 * @brief modify fd
 */
int32 CMyEPoll::modfd(int32 fd, uint64 event_mask, event_handle * eh)
{
	MYLOG_DEBUG(("CMyEPoll::modfd"));

	uint32 idx = fd % this->v_thrds_.size();

	return this->v_thrds_[idx]->modfd(fd, event_mask, eh);
}

/**
 * @brief work 循环
 */
int32 CMyEPoll::work_loop(int32 timeout)
{
	return this->ewt_->work_loop(timeout);
}

/**
 * @brief view
 */
int32 CMyEPoll::runtime_view()
{
	for(uint32 i = 0; i < this->v_thrds_.size(); i ++)
	{
		this->v_thrds_[i]->runtime_view();
	}

	return 0;
}


/**
 * @brief 构造
 */
epoll_wait_thrd::epoll_wait_thrd(uint32 max_fd_count):
	efd_(epoll_create(max_fd_count)),
	vevent_(max_fd_count),
	max_fd_count_(max_fd_count),
	brun_(1)
{
	MYLOG_INFO(("max_fd_count:%d",
		this->max_fd_count_));

	if(-1 == this->efd_)
		MYLOG_WARN(("create epoll fd err, errno:%d", errno));

	pthread_mutex_init(&this->fd_map_protect_,NULL);
}

/**
 * @brief 析构
 */
epoll_wait_thrd::~epoll_wait_thrd()
{
	MYLOG_DEBUG(("epoll_wait_thrd::~epoll_wait_thrd"));

	if(-1 != this->efd_)
		close(this->efd_);

	pthread_mutex_destroy(&this->fd_map_protect_);
}

/**
 * @brief add fd
 */
int32 epoll_wait_thrd::addfd(int32 fd, uint64 event_mask, CMyEPoll::event_handle * eh)
{
	MYLOG_DEBUG(("CMyEPoll::addfd %d evenet_mask:%lld", fd, event_mask));

	if(-1 == this->efd_ || NULL == eh)
	{
		MYLOG_WARN(("init epoll fd err or bad param"));
		return -1;
	}

	CAutoLock alock(&this->fd_map_protect_);

	std::map<int32, CMyEPoll::event_handle>::iterator it = this->fd_map_.find(fd);
	if(it != this->fd_map_.end())
	{
		MYLOG_WARN(("fd has been add fd:%d", fd));
		return -1;
	}

	struct epoll_event e = {0};
	e.data.fd = fd;

	e.events = EPOLLHUP | EPOLLERR | EPOLLET | EPOLLPRI | EPOLLRDHUP;
	if(event_mask & CMyEPoll::EVENT_INPUT)
		e.events |= EPOLLIN;
	if(event_mask & CMyEPoll::EVENT_OUTPUT)
		e.events |= EPOLLOUT;

	if(0 != epoll_ctl(this->efd_, EPOLL_CTL_ADD, fd, &e))
	{
		MYLOG_WARN(("add to epoll set err :%d errno:%", fd, errno));
		return -1;
	}

	this->fd_map_[fd] = *eh;

	return 0;
}

/**
 * @brief del fd
 */
int32 epoll_wait_thrd::delfd(int32 fd)
{
	MYLOG_DEBUG(("CMyEPoll::delfd %d", fd));

	if(-1 == this->efd_)
	{
		MYLOG_WARN(("init epoll fd err"));
		return -1;
	}

	CAutoLock alock(&this->fd_map_protect_);

	struct epoll_event e = {0};
	if(0 != epoll_ctl(this->efd_, EPOLL_CTL_DEL, fd, NULL/*NULL this param meanless when del,but in low version linux2.6.X, if it's null, will fail...dam... */))
	{
		MYLOG_WARN(("del from epoll set fail :%d errno:%d", fd, errno));
		//return -1;
	}

	std::map<int32, CMyEPoll::event_handle>::iterator it = this->fd_map_.find(fd);
	if(it == this->fd_map_.end())
	{
		MYLOG_WARN(("fd not in poll set :%d", fd));
		return -1;
	}

	MYLOG_DEBUG(("call epoll_ctl %d", fd));

	this->fd_map_.erase(it);

	return 0;
}

/**
 * @brief modify fd
 */
int32 epoll_wait_thrd::modfd(int32 fd, uint64 event_mask, CMyEPoll::event_handle * eh)
{
	MYLOG_DEBUG(("CMyEPoll::modfd fd:%d event_mask:%lld", fd, event_mask));

	if(-1 == this->efd_)
	{
		MYLOG_WARN(("init epoll fd err"));
		return -1;
	}

	CAutoLock alock(&this->fd_map_protect_);

	std::map<int32, CMyEPoll::event_handle>::iterator it = this->fd_map_.find(fd);
	if(it == this->fd_map_.end())
	{
		MYLOG_WARN(("fd not in poll set :%d", fd));
		return -1;
	}

	if(eh)
	{
		MYLOG_DEBUG(("modify event call back"));

		it->second.context_data = eh->context_data;
		it->second.input = eh->input;
		it->second.output = eh->output;
		it->second.exception = eh->exception;
	}

	struct epoll_event e = {0};
	e.data.fd = fd;

	e.events = EPOLLHUP | EPOLLERR | EPOLLET | EPOLLPRI |EPOLLRDHUP;
	if(event_mask & CMyEPoll::EVENT_INPUT)
		e.events |= EPOLLIN;
	if(event_mask & CMyEPoll::EVENT_OUTPUT)
		e.events |= EPOLLOUT;

	if(0 != epoll_ctl(this->efd_, EPOLL_CTL_MOD, fd, &e))
	{
		MYLOG_WARN(("del from epoll set fail :%d errno:%d", fd, errno));
		return -1;
	}

	return 0;
}

/**
 * @brief 启动工作线程
 */
int32 epoll_wait_thrd::work()
{
	MYLOG_DEBUG(("epoll_wait_thrd::work"));

	return pthread_create(&this->thrd_, NULL, epoll_wait_thrd::thrd_fun, this);
}

/**
 * @brief 停止工作线程
 */
void epoll_wait_thrd::stop()
{
	MYLOG_DEBUG(("epoll_wait_thrd::stop"));

	this->brun_ = 0;

	if(-1 != this->efd_)
		close(this->efd_);

	pthread_cancel(this->thrd_);
	pthread_join(this->thrd_, NULL);
}

/**
 * @brief view
 */
int32 epoll_wait_thrd::runtime_view()
{
	CAutoLock alock(&this->fd_map_protect_);

	MYLOG_ERREX(("view", "fd map:%d efd:%d vsz:%d max_fd_count_:%d", 
		this->fd_map_.size(), this->efd_, this->vevent_.size(), this->max_fd_count_));

	for(std::map<int32, CMyEPoll::event_handle>::iterator it = this->fd_map_.begin();
		it != this->fd_map_.end();
		it ++)
	{
		MYLOG_ERREX(("view", "%d - %x", it->first, it->second.context_data));
	}

	MYLOG_ERREX(("view", "view end"));

}

/**
 * @brief 工作线程函数
 */
int32 epoll_wait_thrd::work_loop(int32 timeout)
{
	//MYLOG_DEBUG(("_work_loop begin"));
	if(-1 == this->efd_)
	{
		MYLOG_WARN(("init epoll fd err"));
		MYLOG_DEBUG(("loop end"));
		return -1;
	}

	this->vevent_.clear();
	this->vevent_.resize(this->max_fd_count_);

	MYLOG_DEBUG(("epoll_wait begin efd:%d sz:%d time:%d", this->efd_, this->vevent_.size(), timeout));
	int32 ret = epoll_wait(this->efd_, &this->vevent_[0], this->vevent_.size(), timeout);
	MYLOG_DEBUG(("epoll_wait end %d efd:%d", ret, this->efd_));

	if(ret < 0)
	{
		MYLOG_WARN(("epoll_wait err %d", errno));
		MYLOG_DEBUG(("loop end"));
		return -1;
	}

	if(0 == ret)
		return 0;

	//MYLOG_DEBUG(("poll event count:%d", ret));

	for(int i = 0; i < ret; i ++)
	{
		//MYLOG_DEBUG(("fd:%d events:%x", vevent_[i].data.fd, vevent_[i].events));

		CMyEPoll::event_handle e = {0};
		{
			CAutoLock alock(&this->fd_map_protect_);
			std::map<int32, CMyEPoll::event_handle>::iterator it =  this->fd_map_.find(vevent_[i].data.fd);
			if(this->fd_map_.end() == it)
				continue;

			e = it->second;
		}

		if(vevent_[i].events & EPOLLOUT)
		{
			//MYLOG_DEBUG(("handle output"));
			if(e.output)
				e.output(e.context_data, vevent_[i].data.fd);
			else
				MYLOG_WARN(("user not reg the input event call back"));
		}

		if(vevent_[i].events & EPOLLIN)
		{
			//MYLOG_DEBUG(("handle input"));
			if(e.input)
				e.input(e.context_data, vevent_[i].data.fd);
			else
				MYLOG_WARN(("user not reg the output event call back"));
		}

		if(vevent_[i].events & EPOLLERR)
		{
			//MYLOG_DEBUG(("handle exception"));
			if(e.exception)
				e.exception(e.context_data, vevent_[i].data.fd);
			else
				MYLOG_WARN(("user not reg the err event call back"));

			this->delfd(vevent_[i].data.fd);
		}
	}

	//MYLOG_DEBUG(("loop end"));

	return 0;
}

/**
 * @brief 
 */
void * epoll_wait_thrd::thrd_fun(void * param)
{
	MYLOG_DEBUG(("epoll_wait_thrd::thrd_fun"));

	epoll_wait_thrd * ewt = (epoll_wait_thrd *)param;
	while(ewt->brun_)
	{
		ewt->work_loop();
	}

	return NULL;
}






#include "myepoll.h"
#include "AutoLock.h"
#include "mylogex.h"
#include <windows.h>
#define sleep(x) Sleep(x * 1000)

/**
* @brief 构造
*/
CMyEPoll::CMyEPoll(uint32 max_fd_count, int32 wait_thrd_count):
	hlsn_(MyListernerConstruct(NULL, 1024))
{
	::pthread_mutex_init(&this->fd_map_protect_, NULL);
	MyListernerRun(this->hlsn_);
}

/**
* @brief 析构
*/
CMyEPoll::~CMyEPoll()
{
	if(hlsn_)
		MyListernerDestruct(this->hlsn_);
	::pthread_mutex_destroy(&this->fd_map_protect_);
}

/**
* @brief work 循环
*/
int32 CMyEPoll::work_loop(int32 timeout)
{
	sleep(timeout);
	//MyListernerWait(this->hlsn_);
	return 0;
}

/**
* @brief add fd
*/
int32 CMyEPoll::addfd(int32 fd, uint64 event_mask, event_handle * eh)
{
	CAutoLock alock(&this->fd_map_protect_);

	std::map<int32, event_handle>::iterator it = this->fd_map_.find(fd);
	if(it != this->fd_map_.end())
	{
		MYLOG_WARN(("fd has been add fd:%d", fd));
		return -1;
	}

	uint32 mask = E_FD_EXCEPTION;
	if(event_mask & CMyEPoll::EVENT_INPUT)
		mask |= E_FD_READ;
	else if(event_mask & CMyEPoll::EVENT_OUTPUT)
		mask |= E_FD_WRITE;

	event_handle_t e;
	e.context_data = (unsigned long)this;
	e.input = CMyEPoll::_lsn_handle_input;
	e.output = CMyEPoll::_lsn_handle_output;
	e.exception = CMyEPoll::_lsn_handle_err;
	MyListernerAddFD(this->hlsn_, fd, (E_HANDLE_SET_MASK)mask, &e);

	this->fd_map_[fd] = *eh;

	return 0;
}

/**
* @brief del fd
*/
int32 CMyEPoll::delfd(int32 fd)
{
	CAutoLock alock(&this->fd_map_protect_);

	MyListernerDelFD(this->hlsn_, fd);

	std::map<int32, event_handle>::iterator it = this->fd_map_.find(fd);
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
int32 CMyEPoll::modfd(int32 fd, uint64 event_mask, event_handle * eh)
{
	CAutoLock alock(&this->fd_map_protect_);

	std::map<int32, event_handle>::iterator it = this->fd_map_.find(fd);
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

	MyListernerDelFD(this->hlsn_, fd);

	uint32 mask = E_FD_EXCEPTION;
	if(event_mask & CMyEPoll::EVENT_INPUT)
		mask |= E_FD_READ;
	else if(event_mask & CMyEPoll::EVENT_OUTPUT)
		mask |= E_FD_WRITE;

	event_handle_t e;
	e.context_data = (unsigned long)this;
	e.input = CMyEPoll::_lsn_handle_input;
	e.output = CMyEPoll::_lsn_handle_output;
	e.exception = CMyEPoll::_lsn_handle_err;
	MyListernerAddFD(this->hlsn_, fd, (E_HANDLE_SET_MASK)mask, &e);

	return 0;
}

/**
 * @brief 处理有输入事件的回调函数
 */
int32 CMyEPoll::_lsn_handle_input(unsigned long context_data, int fd)
{
	CMyEPoll * ep = (CMyEPoll *)context_data;

	event_handle e = {0};
	{
		CAutoLock alock(&ep->fd_map_protect_);
		std::map<int32, event_handle>::iterator it =  ep->fd_map_.find(fd);
		if(ep->fd_map_.end() == it)
			return 0;

		e = it->second;
	}

	if(e.input)
		e.input(e.context_data, fd);
	else
		MYLOG_WARN(("user not reg the input event call back"));

	return 0;
}

/**
 * @brief 处理有输出事件的回调函数
 */
int32 CMyEPoll::_lsn_handle_output(unsigned long context_data, int fd)
{
	CMyEPoll * ep = (CMyEPoll *)context_data;

	event_handle e = {0};
	{
		CAutoLock alock(&ep->fd_map_protect_);
		std::map<int32, event_handle>::iterator it =  ep->fd_map_.find(fd);
		if(ep->fd_map_.end() == it)
			return 0;

		e = it->second;
	}

	if(e.output)
		e.output(e.context_data, fd);
	else
		MYLOG_WARN(("user not reg the input event call back"));

	return 0;
}

/**
 * @brief 处理有异常事件的回调函数
 */
int32 CMyEPoll::_lsn_handle_err(unsigned long context_data, int fd)
{
	CMyEPoll * ep = (CMyEPoll *)context_data;

	event_handle e = {0};
	{
		CAutoLock alock(&ep->fd_map_protect_);
		std::map<int32, event_handle>::iterator it =  ep->fd_map_.find(fd);
		if(ep->fd_map_.end() == it)
			return 0;

		e = it->second;
	}

	if(e.exception)
		e.exception(e.context_data, fd);
	else
		MYLOG_WARN(("user not reg the input event call back"));

	ep->delfd(fd);

	return 0;
}

/**
* @brief view
*/
int32 CMyEPoll::runtime_view()
{
	CAutoLock alock(&this->fd_map_protect_);

	for(std::map<int32, event_handle>::iterator it = this->fd_map_.begin();
		it != this->fd_map_.end();
		it ++)
	{
		MYLOG_INFO(("%d - %x", it->first, it->second.context_data));
	}

	MYLOG_INFO(("view end fd map:%d ", this->fd_map_.size()));

	return 0;
}

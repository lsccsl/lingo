/**
 * @file mypoll.h 
 * @brief wrapper event poll
 *
 * @author linshaochuan
 * @blog http://blog.csdn.net/lsccsl
 */
#ifndef __MYEPOLL_H__
#define __MYEPOLL_H__


#include <map>
#include <vector>
#include <pthread.h>
#include "type_def.h"


/**
 * @brief 事件处理回调函数
 */
typedef int32 (*POLL_EVENT_CB)(void * context_data, int32 fd);

class epoll_wait_thrd;

/**
 * @brief 
 */
class CMyEPoll
{
public:

	/**
	* @brief 事件回调函数集
	*/
	struct event_handle
	{
		/**
		* @brief 输入事件回调
		*/
		POLL_EVENT_CB input;

		/**
		* @brief 输出事件回调
		*/
		POLL_EVENT_CB output;

		/**
		* @brief 异常事件回调
		*/
		POLL_EVENT_CB exception;

		/**
		* @brief 用户处理事件时的上下文数据
		*/
		void * context_data;
	};

	enum{
		/**
		 * @brief 需要输入事件
		 */
		EVENT_INPUT = 0x01,

		/**
		 * @brief 需要输出事件
		 */
		EVENT_OUTPUT = 0x02,

		/**
		 * @brief 需要异常事件
		 */
		EVENT_ERR = 0x04,
	};

	/**
	 * @brief 构造
	 */
	CMyEPoll(uint32 max_fd_count = 1024, int32 wait_thrd_count = 10);

	/**
	 * @brief 析构
	 */
	~CMyEPoll();

	/**
	 * @brief work 循环
	 */
	int32 work_loop(int32 timeout = -1);

	/**
	 * @brief add fd
	 */
	int32 addfd(int32 fd, uint64 event_mask, event_handle * eh);

	/**
	 * @brief del fd
	 */
	int32 delfd(int32 fd);

	/**
	 * @brief modify fd
	 */
	int32 modfd(int32 fd, uint64 event_mask, event_handle * eh = NULL);

	/**
	 * @brief view
	 */
	int32 runtime_view();

protected:

	/**
	 * @brief epoll_wait线程组,启动后此数组不改变,数组里的指针也不改变,不加锁
	 */
	std::vector<epoll_wait_thrd *> v_thrds_;

	/**
	 * @brief 
	 */
	epoll_wait_thrd * ewt_;
};

/**
 * @brief epoll_wait线程
 */
class epoll_wait_thrd
{
public:

	/**
	 * @brief 构造
	 */
	epoll_wait_thrd(uint32 max_fd_count = 1024);

	/**
	 * @brief 析构
	 */
	~epoll_wait_thrd();

	/**
	 * @brief add fd
	 */
	int32 addfd(int32 fd, uint64 event_mask, CMyEPoll::event_handle * eh);

	/**
	 * @brief del fd
	 */
	int32 delfd(int32 fd);

	/**
	 * @brief modify fd
	 */
	int32 modfd(int32 fd, uint64 event_mask, CMyEPoll::event_handle * eh = NULL);

	/**
	 * @brief 启动工作线程
	 */
	int32 work();

	/**
	 * @brief 工作线程函数
	 */
	int32 work_loop(int32 timeout = -1);

	/**
	 * @brief 停止工作线程
	 */
	void stop();

	/**
	 * @brief view
	 */
	int32 runtime_view();

private:

	/**
	 * @brief 
	 */
	static void * thrd_fun(void * param);

private:

	/**
	* @brief epoll fd
	*/
	int32 efd_;

	/**
	* @brief fd map
	*/
	std::map<int32, CMyEPoll::event_handle> fd_map_;
	/* fd_map_的保护锁 */
	pthread_mutex_t fd_map_protect_;

	/**
	* @brief 容纳被触发的事件数据
	*/
	std::vector<struct epoll_event> vevent_;

	/**
	* @brief 最大监听数
	*/
	uint32 max_fd_count_;

	/**
	* @brief 线程id
	*/
	pthread_t thrd_;

	/**
	 * @brief run
	 */
	int32 brun_;
};

#endif


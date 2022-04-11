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
 * @brief �¼�����ص�����
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
	* @brief �¼��ص�������
	*/
	struct event_handle
	{
		/**
		* @brief �����¼��ص�
		*/
		POLL_EVENT_CB input;

		/**
		* @brief ����¼��ص�
		*/
		POLL_EVENT_CB output;

		/**
		* @brief �쳣�¼��ص�
		*/
		POLL_EVENT_CB exception;

		/**
		* @brief �û������¼�ʱ������������
		*/
		void * context_data;
	};

	enum{
		/**
		 * @brief ��Ҫ�����¼�
		 */
		EVENT_INPUT = 0x01,

		/**
		 * @brief ��Ҫ����¼�
		 */
		EVENT_OUTPUT = 0x02,

		/**
		 * @brief ��Ҫ�쳣�¼�
		 */
		EVENT_ERR = 0x04,
	};

	/**
	 * @brief ����
	 */
	CMyEPoll(uint32 max_fd_count = 1024, int32 wait_thrd_count = 10);

	/**
	 * @brief ����
	 */
	~CMyEPoll();

	/**
	 * @brief work ѭ��
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
	 * @brief epoll_wait�߳���,����������鲻�ı�,�������ָ��Ҳ���ı�,������
	 */
	std::vector<epoll_wait_thrd *> v_thrds_;

	/**
	 * @brief 
	 */
	epoll_wait_thrd * ewt_;
};

/**
 * @brief epoll_wait�߳�
 */
class epoll_wait_thrd
{
public:

	/**
	 * @brief ����
	 */
	epoll_wait_thrd(uint32 max_fd_count = 1024);

	/**
	 * @brief ����
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
	 * @brief ���������߳�
	 */
	int32 work();

	/**
	 * @brief �����̺߳���
	 */
	int32 work_loop(int32 timeout = -1);

	/**
	 * @brief ֹͣ�����߳�
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
	/* fd_map_�ı����� */
	pthread_mutex_t fd_map_protect_;

	/**
	* @brief ���ɱ��������¼�����
	*/
	std::vector<struct epoll_event> vevent_;

	/**
	* @brief ��������
	*/
	uint32 max_fd_count_;

	/**
	* @brief �߳�id
	*/
	pthread_t thrd_;

	/**
	 * @brief run
	 */
	int32 brun_;
};

#endif


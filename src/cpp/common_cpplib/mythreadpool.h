/**
* @file mythreadpool.h 
* @brief �̳߳�
* @author linsc
* @blog http://blog.csdn.net/lsccsl
*/
#ifndef __MYTHREADPOLL_H__
#define __MYTHREADPOLL_H__

#include <pthread.h>
#include <vector>
#include <list>
#include <semaphore.h>
#include "type_def.h"

extern "C"
{
	#include "mymsgque.h"
	#include "myevent.h"
	#include "mymutex.h"
}

/**
 * @brief callback
 */
typedef void *(*thread_job_fun_call_back)(const void * context_data, void * user_data);
class mythreadpool;

/**
 * @brief �߳�˽��������Ϣ
 */
class mythreadworker_data
{
public:

	/**
	 * @brief constructor
	 */
	mythreadworker_data(const int32 b_owned_by_thrd_pool = true):b_owned_by_thrd_pool_(b_owned_by_thrd_pool){}

	/**
	 * @brief destructor
	 */
	virtual ~mythreadworker_data(){}

	/**
	 * @brief is owned by thrd pool
	 */
	int32 b_owned_by_thrd_pool(){ return this->b_owned_by_thrd_pool_; }

private:

	/**
	 * @brief ��ǰ�����Ƿ�Ϊ�߳�����,���Ϊ��,����������̳߳ر�����ʱ����delete,����delete,���ɷ�����������delete;
	 */
	const int32 b_owned_by_thrd_pool_;
};

/**
 * @brief �����߳�
 */
class mythreadworker
{
public:

	/**
	 * @brief job��Ϣ
	 */
	struct threadjob_data
	{
		threadjob_data():cb_ignore_(NULL),cb_(NULL),context_data_(NULL),user_data_(NULL){}

		/**
		 * @brief ִ��������������(��ñ�֤��ֵ���������ڲ�����)
		 */
		void * context_data_;

		/**
		 * @brief ִ�������������(����Ǹ�����,Ҳ������malloc�������ڴ����������Ϣ,��cb�����ͷ���Ӧ���ڴ��)
		 */
		void * user_data_;

		/**
		 * @brief ִ�еĺ���
		 */
		thread_job_fun_call_back cb_;
		/* @brief ���̳߳��е��̸߳����ﵽ����,���д˺��� */
		thread_job_fun_call_back cb_ignore_;
	};

public:

	/**
	 * @brief constructor
	 */
	mythreadworker(mythreadpool& thrd_poll, const uint32 thrd_idx);

	/**
	 * @brief destructor
	 */
	~mythreadworker();

	/**
	 * @brief ��ʼ�� 0:�ɹ� ����:ʧ��
	 */
	int init();

	/**
	 * @brief ���ѹ���
	 */
	void wake_up(const void * context_data, void * user_data, thread_job_fun_call_back cb);

	/**
	* @brief �鿴(just for debug, thread unsafe)
	*/
	void view();

private:

	/**
	 * @brief
	 */
	static void * _thread_fun(void * param);

private:

	/**
	 * @brief �߳�ֹͣ�˳�
	 */
	void _exit_thrd();

private:

	/**
	 * @brief ��������
	 */
	threadjob_data thrdjob_data_;

	/**
	 * @brief �����̳߳ض���
	 */
	mythreadpool& thrd_poll_;

	/**
	 * @brief ��һ��job data ���źŴ�������
	 */
	sem_t sem;

	/**
	 * @brief �Ƿ��˳�
	 */
	int32 bexit_;

	/**
	 * @brief �߳�id
	 */
	pthread_t thrd_id_;

	/**
	 * @brief �̱߳��
	 */
	uint32 thrd_idx_;

	/**
	 * @brief �߳�˽������
	 */
	//mythreadworker_data * thrd_data_;
};


/**
 * @brief �̳߳�
 */
class mythreadpool
{
public:

	/**
	 * @brief constructor
	 */
	mythreadpool(const uint32 max_thread_count = 50, const uint32 max_job_delay = 65535,
		int32 b_ignore_when_no_idle = 0);

	/**
	* @brief constructor
	*/
	~mythreadpool();

	/**
	 * @brief ����job
	 */
	void push_job(const void * context_data, void * user_data, thread_job_fun_call_back cb,
		thread_job_fun_call_back cb_ignore = NULL);

	/**
	 * @brief �鿴(just for debug, thread unsafe)
	 */
	void view();

protected:

	friend class mythreadworker;

	/**
	 * @brief ��־�߳�Ϊ��
	 */
	void report_idle(const uint32 thrd_idx);

private:

	/**
	 * @brief �̺߳���
	 */
	static void * _thread_fun(void * param);

private:

	/**
	 * @brief �˳�
	 */
	void _exit_thrd_pool();

private:

	/**
	 * @brief �������л���
	 */
	HMYMSGQUE hjob_que_;

	/**
	 * @brief �߳�����
	 */
	std::vector<mythreadworker *> v_work_;
	uint32 max_thread_count_;

	/**
	 * @brief �����߳�ջ
	 */
	std::list<uint32> s_idle_;
	HMYMUTEX s_idle_protector_;
	sem_t s_idle_signal_;

	/**
	 * @brief �Ƿ��˳�
	 */
	int32 bexit_;

	/**
	* @brief �߳�id
	*/
	pthread_t thrd_id_;

	/**
	* @brief ���̳߳����̸߳����ﵽ��ֵʱ,�Ƿ�����µ���Ϣ
	*/
	int32 b_ignore_when_no_idle_;
};

#endif










/**
* @file mythreadpool.h 
* @brief 线程池
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
 * @brief 线程私有数据信息
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
	 * @brief 当前对象是否为线程所有,如果为真,则该数据在线程池被析构时将被delete,否则不delete,交由分配者来负责delete;
	 */
	const int32 b_owned_by_thrd_pool_;
};

/**
 * @brief 工作线程
 */
class mythreadworker
{
public:

	/**
	 * @brief job信息
	 */
	struct threadjob_data
	{
		threadjob_data():cb_ignore_(NULL),cb_(NULL),context_data_(NULL),user_data_(NULL){}

		/**
		 * @brief 执行所处的上下文(最好保证此值在运行期内不析构)
		 */
		void * context_data_;

		/**
		 * @brief 执行所处理的数据(最好是个索引,也可以用malloc出来的内存块来记载信息,由cb负责释放相应的内存块)
		 */
		void * user_data_;

		/**
		 * @brief 执行的函数
		 */
		thread_job_fun_call_back cb_;
		/* @brief 当线程池中的线程个数达到顶峰,呼叫此函数 */
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
	 * @brief 初始化 0:成功 其它:失败
	 */
	int init();

	/**
	 * @brief 唤醒工作
	 */
	void wake_up(const void * context_data, void * user_data, thread_job_fun_call_back cb);

	/**
	* @brief 查看(just for debug, thread unsafe)
	*/
	void view();

private:

	/**
	 * @brief
	 */
	static void * _thread_fun(void * param);

private:

	/**
	 * @brief 线程停止退出
	 */
	void _exit_thrd();

private:

	/**
	 * @brief 工作数据
	 */
	threadjob_data thrdjob_data_;

	/**
	 * @brief 引用线程池对象
	 */
	mythreadpool& thrd_poll_;

	/**
	 * @brief 用一个job data 与信号传达数据
	 */
	sem_t sem;

	/**
	 * @brief 是否退出
	 */
	int32 bexit_;

	/**
	 * @brief 线程id
	 */
	pthread_t thrd_id_;

	/**
	 * @brief 线程编号
	 */
	uint32 thrd_idx_;

	/**
	 * @brief 线程私有数据
	 */
	//mythreadworker_data * thrd_data_;
};


/**
 * @brief 线程池
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
	 * @brief 发送job
	 */
	void push_job(const void * context_data, void * user_data, thread_job_fun_call_back cb,
		thread_job_fun_call_back cb_ignore = NULL);

	/**
	 * @brief 查看(just for debug, thread unsafe)
	 */
	void view();

protected:

	friend class mythreadworker;

	/**
	 * @brief 标志线程为闲
	 */
	void report_idle(const uint32 thrd_idx);

private:

	/**
	 * @brief 线程函数
	 */
	static void * _thread_fun(void * param);

private:

	/**
	 * @brief 退出
	 */
	void _exit_thrd_pool();

private:

	/**
	 * @brief 工作队列缓存
	 */
	HMYMSGQUE hjob_que_;

	/**
	 * @brief 线程数据
	 */
	std::vector<mythreadworker *> v_work_;
	uint32 max_thread_count_;

	/**
	 * @brief 空闲线程栈
	 */
	std::list<uint32> s_idle_;
	HMYMUTEX s_idle_protector_;
	sem_t s_idle_signal_;

	/**
	 * @brief 是否退出
	 */
	int32 bexit_;

	/**
	* @brief 线程id
	*/
	pthread_t thrd_id_;

	/**
	* @brief 当线程池中线程个数达到峰值时,是否忽略新的消息
	*/
	int32 b_ignore_when_no_idle_;
};

#endif










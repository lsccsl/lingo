/**
* @file mythreadpool.cpp 
* @brief �̳߳�
* @author linsc
*/
#include "mythreadpool.h"
#include "mylogex.h"

/**
* @brief constructor
*/
mythreadworker::mythreadworker(mythreadpool& thrd_poll, const uint32 thrd_idx):bexit_(0),thrd_idx_(thrd_idx),thrd_poll_(thrd_poll)
{
	MYLOG_DEBUG(("mythreadworker::mythreadworker thrd_idx:%d", this->thrd_idx_));

	this->thrdjob_data_.cb_ = NULL;
	this->thrdjob_data_.context_data_ = NULL;
	this->thrdjob_data_.user_data_ = NULL;
}

/**
* @brief destructor
*/
mythreadworker::~mythreadworker()
{
	MYLOG_DEBUG(("mythreadworker::~mythreadworker thrd_idx:%d", this->thrd_idx_));

	this->_exit_thrd();

	sem_destroy(&this->sem);
}

/**
* @brief ��ʼ�� 0:�ɹ� ����:ʧ��
*/
int mythreadworker::init()
{
	MYLOG_DEBUG(("mythreadworker::init thrd_idx:%d", this->thrd_idx_));

	sem_init(&this->sem, 0, 0);

	return pthread_create(&this->thrd_id_, NULL, mythreadworker::_thread_fun, this);
}

/**
* @brief ���ѹ���
*/
void mythreadworker::wake_up(const void * context_data, void * user_data, thread_job_fun_call_back cb)
{
	MYLOG_DEBUG(("mythreadworker::wake_up %d context_data:%x user_data:%x cb:%x", this->thrd_idx_, context_data, user_data, cb));

	this->thrdjob_data_.cb_ = cb;
	this->thrdjob_data_.context_data_ = (void *)context_data;
	this->thrdjob_data_.user_data_ = user_data;

	sem_post(&this->sem);
}

/**
* @brief �߳�ֹͣ�˳�
*/
void mythreadworker::_exit_thrd()
{
	MYLOG_DEBUG(("thread:%d will exit", this->thrd_idx_));

	this->bexit_ = 1;
	this->wake_up(NULL, NULL, NULL);

	pthread_join(this->thrd_id_, NULL);
}

/**
* @brief
*/
void * mythreadworker::_thread_fun(void * param)
{
	MYLOG_DEBUG(("mythreadworker::_thread_fun"));

	mythreadworker * pthis = (mythreadworker *)param;

	while(!pthis->bexit_)
	{
		MYLOG_DEBUG(("thread loop %d", pthis->thrd_idx_));

		/* ����thread poll������� */
		pthis->thrd_poll_.report_idle(pthis->thrd_idx_);

		sem_wait(&pthis->sem);

		if(NULL == pthis->thrdjob_data_.cb_)
		{
			MYLOG_INFO(("cb is null ...%d\r\n", pthis->thrd_idx_));
			continue;
		}

		MYLOG_DEBUG(("begin work context_data_:%x user_data_:%x %d", pthis->thrdjob_data_.context_data_, pthis->thrdjob_data_.user_data_, pthis->thrd_idx_));
		pthis->thrdjob_data_.cb_(pthis->thrdjob_data_.context_data_, pthis->thrdjob_data_.user_data_);
		MYLOG_DEBUG(("end work context_data_:%x user_data_:%x %d\r\n", pthis->thrdjob_data_.context_data_, pthis->thrdjob_data_.user_data_, pthis->thrd_idx_));

		pthis->thrdjob_data_.cb_ = NULL;
		pthis->thrdjob_data_.context_data_ = NULL;
		pthis->thrdjob_data_.user_data_ = NULL;
	}

	MYLOG_DEBUG(("thread:%d exit", pthis->thrd_idx_));

	return NULL;
}

/**
* @brief �鿴(just for debug, thread unsafe)
*/
void mythreadworker::view()
{
	MYLOG_INFOEX(("view", "mythreadpool::view %d context_data_:%x user_data_:%x cb_:%x", this->thrd_idx_,
		this->thrdjob_data_.context_data_, this->thrdjob_data_.user_data_, this->thrdjob_data_.cb_));
}


/**
* @brief constructor
*/
mythreadpool::mythreadpool(const uint32 max_thread_count, const uint32 max_job_delay, int32 b_ignore_when_no_idle)
	:hjob_que_(NULL),bexit_(0),max_thread_count_(max_thread_count),b_ignore_when_no_idle_(b_ignore_when_no_idle)
{
	MYLOG_DEBUG(("mythreadpool::mythreadpool max_thread_count:%d max_job_delay:%d b_ignore_when_no_idle:%d", max_thread_count, max_job_delay, b_ignore_when_no_idle));

	this->hjob_que_ = MyMsgQueConstruct(NULL, max_job_delay);
    
	this->s_idle_protector_ = MyMutexConstruct(NULL);
	sem_init(&this->s_idle_signal_, 0, 0);

	pthread_create(&this->thrd_id_, NULL, mythreadpool::_thread_fun, this);
}

/**
* @brief constructor
*/
mythreadpool::~mythreadpool()
{
	MYLOG_DEBUG(("mythreadpool::~mythreadpool"));

	/* ���ù����߳�ֹͣ���� */
	this->_exit_thrd_pool();

	for(uint32 i = 0; i < this->v_work_.size(); i ++)
	{
		MYLOG_DEBUG(("destruct thrd:%d", i));

		if(NULL == v_work_[i])
			continue;

		delete v_work_[i];
	}

	MyMsgQueDestruct(this->hjob_que_);
	sem_destroy(&this->s_idle_signal_);
	MyMutexDestruct(this->s_idle_protector_);
}

/**
* @brief �˳�
*/
void mythreadpool::_exit_thrd_pool()
{
	MYLOG_DEBUG(("mythreadpool::_exit_thrd_pool"));

	this->bexit_ = 1;
	MyMsgQuePush_block(this->hjob_que_, NULL);
	pthread_join(this->thrd_id_, NULL);
}

/**
* @brief ����job
*/
void mythreadpool::push_job(const void * context_data, void * user_data, thread_job_fun_call_back cb,
	thread_job_fun_call_back cb_ignore)
{
	MYLOG_DEBUG(("mythreadpool::push_job context_data:%x user_data:%x cb:%x", context_data, user_data, cb));

	mythreadworker::threadjob_data * job = new mythreadworker::threadjob_data;
	job->cb_ = cb;
	job->cb_ignore_ = cb_ignore;
	job->context_data_ = (void *)context_data;
	job->user_data_ = user_data;

	MyMsgQuePush_block(this->hjob_que_, job);
}

/**
* @brief ��־�߳�Ϊ��
*/
void mythreadpool::report_idle(const uint32 thrd_idx)
{
	MYLOG_DEBUG(("mythreadpool::report_idle thrd_idx:%d", thrd_idx));

	/* �����ź�,��ջ */
	MyMutexLock(this->s_idle_protector_);

	this->s_idle_.push_front(thrd_idx);

	if(1 == this->s_idle_.size())
	{
		MYLOG_DEBUG(("need wait up"));
		sem_post(&this->s_idle_signal_);
	}

	MyMutexUnLock(this->s_idle_protector_);
}

/**
* @brief �̺߳���
*/
void * mythreadpool::_thread_fun(void * param)
{
	MYLOG_DEBUG(("mythreadpool::_thread_fun"));

	mythreadpool * pthis = (mythreadpool *)param;

	while(!pthis->bexit_)
	{
		MYLOG_DEBUG(("thread poll mgr loop"));

		mythreadworker::threadjob_data * job = (mythreadworker::threadjob_data *)MyMsgQuePop_block(pthis->hjob_que_);

		if(NULL == job)
		{
			MYLOG_INFO(("job is null ...\r\n"));
			continue;
		}

		int32 need_call_ignore = 0;
		while(1)
		{
			MYLOG_DEBUG(("begin find thread"));

			MyMutexLock(pthis->s_idle_protector_);
			if(pthis->s_idle_.empty())
			{
				MYLOG_DEBUG(("idle stack is empty total thread:%d max:%d", pthis->v_work_.size(), pthis->max_thread_count_));

				MyMutexUnLock(pthis->s_idle_protector_);

				if(pthis->v_work_.size() >= pthis->max_thread_count_)
				{
					MYLOG_DEBUG(("reach max wait idle thread signal"));

					if(!pthis->b_ignore_when_no_idle_)
					{
						MYLOG_DEBUG(("not ignore where have no idle"));

						/* �߳������ﵽ����ʱ,�ȴ������ź� */
						sem_wait(&pthis->s_idle_signal_);
						MYLOG_DEBUG(("have signal, find again"));
						continue;
					}
					else
					{
						MYLOG_DEBUG(("ignore where have no idle"));
						need_call_ignore = 1;
						break;
					}
				}
				else
				{
					MYLOG_DEBUG(("not reach max, create new thread and work"));

					/* �������߳� */
					mythreadworker * thrd_worker = new mythreadworker(*pthis, (uint32)pthis->v_work_.size());
					pthis->v_work_.push_back(thrd_worker);

					thrd_worker->init();
					MYLOG_DEBUG(("create new thread end, and wait for it ready"));
					sem_wait(&pthis->s_idle_signal_);
					MYLOG_DEBUG(("have signal, find again"));

					continue;
				}
			}
			else
			{
				MYLOG_DEBUG(("have idle thread"));

				/* ȡһ���߳������� */
				uint32 thrd_idx = *(pthis->s_idle_.begin());
				pthis->s_idle_.pop_front();

				MyMutexUnLock(pthis->s_idle_protector_);

				assert(thrd_idx < pthis->v_work_.size());
				assert(pthis->v_work_[thrd_idx]);

				pthis->v_work_[thrd_idx]->wake_up(job->context_data_, job->user_data_, job->cb_);

				break;
			}
		}

		MYLOG_DEBUG(("push job end\r\n"));

		if(need_call_ignore)
		{
			MYLOG_DEBUG(("need call ignore"));

			if(job->cb_ignore_)
				job->cb_ignore_(job->context_data_, job->user_data_);
			else
				MYLOG_INFO(("no cb ignore, may cause memleak ..."));
		}

		delete job;
	}

	MYLOG_DEBUG(("thread pool mgr thread exit"));

	return NULL;
}

/**
* @brief �鿴(for debug, thread unsafe)
*/
void mythreadpool::view()
{
	MYLOG_INFOEX(("view", "threadpoll view begin max:%d thread:%d=================== ", this->max_thread_count_, this->v_work_.size()));

	/* �鿴����ջ */
	MyMutexLock(this->s_idle_protector_);
	for(std::list<uint32>::iterator it = this->s_idle_.begin(); it != this->s_idle_.end(); it ++)
	{
		MYLOG_INFOEX(("view", "idle:%d", *it));
	}
	MyMutexUnLock(this->s_idle_protector_);

	/* ���µĴ���Ϊ���̲߳���ȫ */
	for(uint32 i = 0; i < this->v_work_.size(); i ++)
	{
		this->v_work_[i]->view();
	}

	MYLOG_INFOEX(("view", "delay job count:%d", MyMsgQueGetCount(this->hjob_que_)));

	MYLOG_INFOEX(("view", "threadpoll view end =================== \r\n"));
}










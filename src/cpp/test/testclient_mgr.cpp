#include "testclient_mgr.h"
#include <windows.h>
#include "mylogex.h"
#include "cfg.h"


void testclient_mgr::init(const int unit_count, const int test_count)
{
	this->test_count_ = test_count;
	this->v_mgr_unit_.resize(unit_count);
}
void testclient_mgr::run_thread()
{
	for (int i = 0; i < this->v_mgr_unit_.size(); i++)
	{
		testclient_mgr_unit& mgr_unit = this->v_mgr_unit_[i];
		mgr_unit.idx = i;
		mgr_unit.thread_ptr_ = new std::thread(&testclient_mgr::thread_func, this, i);
	}
}

void testclient_mgr::join()
{
	for (int i = 0; i < this->v_mgr_unit_.size(); i++)
	{
		testclient_mgr_unit& mgr_unit = this->v_mgr_unit_[i];
		mgr_unit.thread_ptr_->join();
	}
}

void testclient_mgr::thread_func(int idx)
{
	MYLOG_ERR(("idx:%d\r\n", idx));

	testclient_mgr_unit& mgr_unit = this->v_mgr_unit_[idx];
	
	for (auto& it : mgr_unit.map_client_)
		it.second->connect_to_srv(__global_cfg_.remote_ip_, 2003);

	msgpacket::MSG_TEST msg;
	msg.set_str("testabcdefg");

	while (1)
	{
		if (mgr_unit.map_client_.empty())
			Sleep(1000);

		for (auto& it : mgr_unit.map_client_)
		{
			auto cli = it.second;
			msg.set_id(cli->id());
			if (!cli->send_test(msg, mgr_unit.seq, this->test_count_))
			{
				MYLOG_ERR(("send err, re connect to srv :%lld fd:%d------\r\n", it.second->id(), it.second->fd()));
				cli->connect_to_srv(__global_cfg_.remote_ip_, 2003);
			}
		}
		for (auto& it : mgr_unit.map_client_)
		{
			auto cli = it.second;
			msg.set_id(cli->id());
			if (!cli->recv_test(mgr_unit.seq, this->test_count_))
			{
				MYLOG_ERR(("recv err, re connect to srv :%lld fd:%d------\r\n", it.second->id(), it.second->fd()));
				cli->connect_to_srv(__global_cfg_.remote_ip_, 2003);
			}
		}

		mgr_unit.seq += this->test_count_;
	}
}

void testclient_mgr::add_client(const int64 id)
{
	size_t sz = this->v_mgr_unit_.size();

	size_t idx = id % sz;

	this->v_mgr_unit_[idx].map_client_[id] = new testclient(id);
}

void testclient_mgr::dump()
{
	int64 total_diff = 0;
	int64 total_count = 0;
	int64 max_diff = 0;
	int64 min_diff = INT64_MAX;
	int64 max_sendloop_interval = 0;
	int64 min_sendloop_interval = INT64_MAX;

	int64 total_sendloop_interval = 0;
	int64 total_send_loop = 0;

	std::set<int> setTmpNotRun;
	std::set<int> setFailLogin;

	for (auto& it : this->v_mgr_unit())
	{
		for (auto& itM : it.map_client_)
		{
			total_diff += itM.second->tc_static().total_diff;
			total_count += itM.second->tc_static().total_count;

			if (itM.second->tc_static().max_diff > max_diff)
				max_diff = itM.second->tc_static().max_diff;
			if (itM.second->tc_static().min_diff < min_diff)
				min_diff = itM.second->tc_static().min_diff;

			if (itM.second->tc_static().max_sendloop_interval > max_sendloop_interval)
				max_sendloop_interval = itM.second->tc_static().max_sendloop_interval;
			if (itM.second->tc_static().min_sendloop_interval < min_sendloop_interval)
				min_sendloop_interval = itM.second->tc_static().min_sendloop_interval;

			total_sendloop_interval += itM.second->tc_static().t_last_sendloop - itM.second->tc_static().t_first_sendloop;
			total_send_loop += itM.second->tc_static().total_send_loop;

			if (!itM.second->tc_static().b_login_suc)
				setFailLogin.insert(itM.first);
		}

		if (it.lastSample_seq == it.seq)
			setTmpNotRun.insert(it.idx);

		it.lastSample_seq = it.seq;
	}
	if (total_send_loop < 1)
		total_send_loop = 1;
	int64 aver_total_send_loop = total_sendloop_interval / total_send_loop;
	if (total_count <= 0)
		total_count = 1;
	double aver = (total_diff / 1000.f) / total_count;
	printf("total_count:%lld aver rtt:%fs max_diff:%lldms, min_diff:%lldms\r\n"
		"max_sendloop_interval:%lldms min_sendloop_interval:%lldms aver sendloop:%lldms\r\n"
		"not_run_count:%zd\r\n",
		total_count, aver, max_diff, min_diff,
		max_sendloop_interval, min_sendloop_interval, aver_total_send_loop, setTmpNotRun.size());

	if (this->set_not_run_.empty())
	{
		for (auto it : setTmpNotRun)
			this->set_not_run_.insert(it);
	}

	{
		printf("cur not run:\r\n");
		for (auto it : setTmpNotRun) {
			printf("%d ", it);
		}
		printf("\r\n");
	}

	{
		std::list<int> lstDel;
		for (auto it : this->set_not_run_) {
			if (setTmpNotRun.end() == setTmpNotRun.find(it))
				this->set_not_run_.erase(it);
		}
	}

	{
		printf("always not run:\r\n");
		for (auto it : this->set_not_run_) {
			printf("%d ", it);
		}
		printf("\r\n");
	}

	{
		printf("fail login:\r\n");
		for (auto it : setFailLogin) {
			printf("%d ", it);
		}
		printf("\r\n");
	}
}

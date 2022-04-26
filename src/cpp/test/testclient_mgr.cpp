#include "testclient_mgr.h"
#include <windows.h>
#include "mylogex.h"


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
		it.second->connect_to_srv("192.168.2.129", 2003);

	for (auto& it : mgr_unit.map_client_)
		it.second->do_login();

	int64 seq = 0;

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
			if (!cli->send_test(msg, seq, this->test_count_))
			{
				MYLOG_ERR(("send err, re connect to srv :%lld------\r\n", it.second->id()));
				cli->connect_to_srv("192.168.2.129", 2003);
				cli->do_login();
			}
		}
		for (auto& it : mgr_unit.map_client_)
		{
			auto cli = it.second;
			msg.set_id(cli->id());
			if (!cli->recv_test(seq, this->test_count_))
			{
				MYLOG_ERR(("recv err, re connect to srv :%lld------\r\n", it.second->id()));
				cli->connect_to_srv("192.168.2.129", 2003);
				cli->do_login();
			}
		}

		seq += this->test_count_;
	}
}

void testclient_mgr::add_client(const int64 id)
{
	size_t sz = this->v_mgr_unit_.size();

	size_t idx = id % sz;

	this->v_mgr_unit_[idx].map_client_[id] = new testclient(id);
}

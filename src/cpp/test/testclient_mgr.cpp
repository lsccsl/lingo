#include "testclient_mgr.h"
#include <windows.h>

void testclient_mgr::init(const int unit_count)
{
	this->v_mgr_unit_.resize(unit_count);
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
	printf("idx:%d\r\n", idx);

	testclient_mgr_unit& mgr_unit = this->v_mgr_unit_[idx];
	
	for (auto& it : mgr_unit.map_client_)
		it.second->connect_to_srv("192.168.2.129", 2003);

	for (auto& it : mgr_unit.map_client_)
		it.second->do_login();

	int64 seq = 0;
	while (1)
	{
		if (mgr_unit.map_client_.empty())
			Sleep(1000);
		for (auto& it : mgr_unit.map_client_)
		{
			if (!it.second->send_test(seq, 10))
			{
				printf("re connect to srv :%d------\r\n", it.second->id());
				it.second->connect_to_srv("192.168.2.129", 2003);
				it.second->do_login();
			}
		}

		seq++;
	}
}

void testclient_mgr::add_client(const int64 id)
{
	size_t sz = this->v_mgr_unit_.size();

	size_t idx = id % sz;

	this->v_mgr_unit_[idx].map_client_[id] = new testclient(id);
}

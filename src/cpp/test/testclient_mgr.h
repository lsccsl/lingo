#pragma once

#include <map>
#include <thread>
#include "type_def.h"
#include "testclient.h"

struct testclient_mgr_unit
{
	std::thread* thread_ptr_ = nullptr;
	std::map<int64,testclient*> map_client_;
};

class testclient_mgr
{
public:

	testclient_mgr() {}
	~testclient_mgr() {}

	void join();

	void init(const int unit_count);
	void run_thread();

	void add_client(const int64 id);

private:

	void thread_func(int idx);

private:
	
	std::vector<testclient_mgr_unit> v_mgr_unit_;
};

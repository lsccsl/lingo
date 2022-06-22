#pragma once

#include <map>
#include <thread>
#include "type_def.h"
#include "testclient.h"

struct testclient_mgr_unit
{
	std::thread* thread_ptr_ = nullptr;
	std::map<int64,testclient*> map_client_;
	int64 seq = 0;
	mutable int64 lastSample_seq = 0;
	int idx = 0;
};

struct testclient_mgr_static
{
	int64 total_count_last = 0;
	int64 t_last = 0;
};

class testclient_mgr
{
public:

	testclient_mgr() {}
	~testclient_mgr() {}

	void join();

	void init(const int unit_count, const int test_count);
	void run_thread();

	void add_client(const int64 id);

	const std::vector<testclient_mgr_unit>& v_mgr_unit() const{
		return v_mgr_unit_;
	}

	void dump();

private:

	void thread_func(int idx);

private:
	
	std::vector<testclient_mgr_unit> v_mgr_unit_;
	int test_count_ = 10;

	testclient_mgr_static static_;
	std::set<int> set_not_run_;
};

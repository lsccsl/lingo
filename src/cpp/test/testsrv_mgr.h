#pragma once

#include <map>
#include <thread>
#include <vector>
#include "type_def.h"
#include "testsrv.h"


struct testsrv_mgr_unit
{
	std::thread* thread_acpt_ptr_ = nullptr;
	std::thread* thread_dial_ptr_ = nullptr;
	std::map<int64, testsrv*> map_srv_;
	int64 seq = 0;
	mutable int64 lastSample_seq = 0;
	int idx = 0;
};

class testsrv_mgr
{
public:

	testsrv_mgr();
	~testsrv_mgr();

	void init_srv(const int unit_count, const int test_count);
	void run_srv_thread();

	void add_srv(const int64 id,
		const std::string& local_ip, const int local_port,
		const std::string& remote_ip, const int remote_port);

	void dump();

private:

	void thread_acpt_func(int idx);
	void thread_dial_func(int idx);

private:

	std::vector<testsrv_mgr_unit> v_mgr_unit_;
	int test_count_ = 10;

};

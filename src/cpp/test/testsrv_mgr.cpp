#include "testsrv_mgr.h"
#include <windows.h>
#include "mylogex.h"

testsrv_mgr::testsrv_mgr()
{}

testsrv_mgr::~testsrv_mgr()
{}


void testsrv_mgr::init_srv(const int unit_count, const int test_count)
{
	this->test_count_ = test_count;
	this->v_mgr_unit_.resize(unit_count);
}
void testsrv_mgr::run_srv_thread()
{
	for (int i = 0; i < this->v_mgr_unit_.size(); i++)
	{
		testsrv_mgr_unit& mgr_unit = this->v_mgr_unit_[i];
		mgr_unit.idx = i;
		mgr_unit.thread_ptr_ = new std::thread(&testsrv_mgr::thread_srv_func, this, i);
	}
}

void testsrv_mgr::add_srv(const int64 id,
	const std::string& local_ip, const int local_port,
	const std::string& remote_ip, const int remote_port)
{
	size_t sz = this->v_mgr_unit_.size();

	size_t idx = id % sz;

	auto psrv = new testsrv(id,
		local_ip, local_port,
		remote_ip, remote_port);
	this->v_mgr_unit_[idx].map_srv_[id] = psrv;

	psrv->init_listen();
}

void testsrv_mgr::thread_srv_func(int idx)
{
	MYLOG_ERR(("idx:%d\r\n", idx));

	testsrv_mgr_unit& mgr_unit = this->v_mgr_unit_[idx];

	for (auto& it : mgr_unit.map_srv_)
		it.second->connect_to_srv();

	msgpacket::MSG_TEST msg;
	msg.set_str("testabcdefg");

	while (1)
	{
		if (mgr_unit.map_srv_.empty())
			Sleep(1000);

		// acpt
		for (auto& it : mgr_unit.map_srv_)
			it.second->accept_client_no_block();

		// send rpc
		for (auto& it : mgr_unit.map_srv_)
		{
		}

		// recv rpc response
		for (auto& it : mgr_unit.map_srv_)
		{
		}

		// recv rpc request
		for (auto& it : mgr_unit.map_srv_)
		{
		}

		// check dial tcp fail


		mgr_unit.seq += this->test_count_;
	}

}

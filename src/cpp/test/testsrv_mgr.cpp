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
		mgr_unit.thread_acpt_ptr_ = new std::thread(&testsrv_mgr::thread_acpt_func, this, i);
		mgr_unit.thread_dial_ptr_ = new std::thread(&testsrv_mgr::thread_dial_func, this, i);
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

void testsrv_mgr::thread_dial_func(int idx)
{
	MYLOG_ERR(("dial idx:%d\r\n", idx));

	testsrv_mgr_unit& mgr_unit = this->v_mgr_unit_[idx];

	for (auto& it : mgr_unit.map_srv_)
		it.second->connect_to_srv();

	int64 seq = 0;

	while (1)
	{		
		if (mgr_unit.map_srv_.empty())
			Sleep(1000);

		// send rpc
		for (auto& it : mgr_unit.map_srv_)
		{
			testsrv* srv = it.second;
			if (!srv->send_test_rpc(seq, 120000))
			{
				MYLOG_ERR(("send fail, reconnect srv:%lld connect err:%d-%d", srv->srvid(), ::WSAGetLastError(), ::GetLastError()));
				srv->connect_to_srv();
			}
		}

		// recv rpc response
		for (auto& it : mgr_unit.map_srv_)
		{
			testsrv* srv = it.second;
			if (!srv->recv_test_rpc_res(seq))
			{
				MYLOG_ERR(("recv fail, reconnect srv:%lld connect err:%d-%d", srv->srvid(), ::WSAGetLastError(), ::GetLastError()));
				srv->connect_to_srv();
			}
		}

		mgr_unit.seq += this->test_count_;

		//if (seq % 10 == 0)
		//	Sleep(0);
	}
}

void testsrv_mgr::thread_acpt_func(int idx)
{
	MYLOG_ERR(("acpt idx:%d\r\n", idx));

	testsrv_mgr_unit& mgr_unit = this->v_mgr_unit_[idx];

	int64 seq = 0;

	for (auto& it : mgr_unit.map_srv_)
		it.second->http_addsrv();
	Sleep(1000);

	for (auto& it : mgr_unit.map_srv_)
		it.second->accept_client_no_block();

	while (1)
	{
		if (mgr_unit.map_srv_.empty())
			Sleep(1000);

		// recv rpc request
		for (auto& it : mgr_unit.map_srv_)
		{
			testsrv* srv = it.second;
			if (!srv->process_acpt_msg())
				srv->accept_client_no_block();
		}

		mgr_unit.seq += this->test_count_;

		//if (seq % 10 == 0)
		//	Sleep(0);
	}
}

void testsrv_mgr::dump()
{
	//是否连接上 是否成功accept
	std::list<testsrv*> set_no_acpt;
	std::list<testsrv*> set_no_conn;
	for (auto& it : this->v_mgr_unit_)
	{
		const testsrv_mgr_unit& tmu = it;

		for (auto& it : tmu.map_srv_)
		{
			testsrv* srv = it.second;

			if (srv->di().fd_dial_ <= 0)
				set_no_conn.emplace_back(srv);
			if (srv->ai().fd_acpt_ <= 0)
				set_no_acpt.emplace_back(srv);
		}
	}

	printf("no acpt:%zd\r\n", set_no_acpt.size());
	for (auto& it : set_no_acpt)
		printf("%lld:%d ", it->srvid(), it->ai().last_err);

	printf("\r\n\r\nno conn:%zd\r\n", set_no_conn.size());
	for (auto& it : set_no_conn)
		printf("%lld:%d ", it->srvid(), it->di().last_err);
}

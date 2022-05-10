#include <stdio.h>
#include <windows.h>
#include <string>
#include <iostream>
#include <sstream>
#include "curl/curl-7.83.0/include/curl/curl.h"

#include "msg.pb.h"
#include "testclient_mgr.h"
#include "channel.h"
#include "testsrv_mgr.h"
#include "cfg.h"

testclient_mgr * __tm_;
testsrv_mgr* __tsm_;

testcfg __global_cfg_ = {
   "192.168.2.129",
   2003,
   "10.0.14.48"
};


void test_cmd()
{
	std::set<int> setNotRun;
	while (1)
	{
		std::string line;
		std::getline(std::cin, line);

		if ("dump" == line)
		{
			printf("dump:");
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

			for (auto& it : __tm_->v_mgr_unit())
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

			if (setNotRun.empty())
			{
				for (auto it : setTmpNotRun)
					setNotRun.insert(it);
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
				for (auto it : setNotRun) {
					if (setTmpNotRun.end() == setTmpNotRun.find(it))
						setNotRun.erase(it);
				}
			}

			{
				printf("always not run:\r\n");
				for (auto it : setNotRun) {
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
	}
}

int main(int argc, char* argv[])
{
	CURLcode ret = curl_global_init(CURL_GLOBAL_ALL);
	if (CURLE_OK != ret)
	{
		printf("curl init fail");
		return 0;
	}
	//testsrv::httpRequest(100, "10.0.0.1", 8686);

	MYLOG_SET_LOG_DIRECTION(2);
	CChannel::init_sock();
	msgpacket::MSG_RPC_RES msg;
	msg.set_msg_id(123);
	printf("test:%s", msg.DebugString().c_str());
	
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST, msgpacket::_MSG_TEST);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST, msgpacket::_MSG_TEST);
	msgpackhelp::parse_reg(new msgpacket::MSG_LOGIN_RES, msgpacket::_MSG_LOGIN_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RES, msgpacket::_MSG_TEST_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_LOGIN_RES, msgpacket::_MSG_LOGIN_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_RPC, msgpacket::_MSG_RPC);
	msgpackhelp::parse_reg(new msgpacket::MSG_RPC_RES, msgpacket::_MSG_RPC_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_SRV_REPORT, msgpacket::_MSG_SRV_REPORT);
	msgpackhelp::parse_reg(new msgpacket::MSG_SRV_REPORT_RES, msgpacket::_MSG_SRV_REPORT_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_HEARTBEAT, msgpacket::_MSG_HEARTBEAT);
	msgpackhelp::parse_reg(new msgpacket::MSG_HEARTBEAT_RES, msgpacket::_MSG_HEARTBEAT_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TCP_STATIC, msgpacket::_MSG_TCP_STATIC);
	msgpackhelp::parse_reg(new msgpacket::MSG_TCP_STATIC_RES, msgpacket::_MSG_TCP_STATIC_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST, msgpacket::_MSG_TEST);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RES, msgpacket::_MSG_TEST_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_LOGIN, msgpacket::_MSG_LOGIN);
	msgpackhelp::parse_reg(new msgpacket::MSG_LOGIN_RES, msgpacket::_MSG_LOGIN_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RPC, msgpacket::_MSG_TEST_RPC);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RPC_RES, msgpacket::_MSG_TEST_RPC_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RPC_RES, msgpacket::_MSG_TEST_RPC_RES);

	int is_server = 0;
	if (argc >= 2)
		is_server = ::_atoi64(argv[1]);

	if (is_server)
	{
		int64 srvid_base = 100;
		int thread_count = 1;
		int test_count = 1;
		int srv_count = 1;
		if (argc >= 3)
			srvid_base = ::_atoi64(argv[2]);
		if (argc >= 4)
			srv_count = atoi(argv[3]);
		if (argc >= 5)
			test_count = atoi(argv[4]);
		if (argc >= 6)
			thread_count = atoi(argv[5]);

		__tsm_ = new testsrv_mgr;

		__tsm_->init_srv(thread_count, test_count);
		for (int i = 0; i < srv_count; i++)
		{
			__tsm_->add_srv(srvid_base + i,
				__global_cfg_.local_ip_, 10000 + srvid_base + i,
				__global_cfg_.remote_ip_, 2003);
		}
		__tsm_->run_srv_thread();
	}
	else
	{
		int64 id_base = 1000;
		int client_count = 1000;
		int test_count = 10;
		int thread_count = 50;
		if (argc >= 3)
			id_base = ::_atoi64(argv[2]);
		if (argc >= 4)
			client_count = atoi(argv[3]);
		if (argc >= 5)
			test_count = atoi(argv[4]);
		if (argc >= 6)
			thread_count = atoi(argv[5]);

		__tm_ = new testclient_mgr;

		__tm_->init(thread_count, test_count);
		for (int i = 0; i < client_count; i++)
			__tm_->add_client(id_base + i);
		__tm_->run_thread();
	}

	test_cmd();

	if (__tm_)
		__tm_->join();
}

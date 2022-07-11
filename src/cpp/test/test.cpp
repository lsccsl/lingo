#include <stdio.h>
#include <windows.h>
#include <string>
#include <iostream>
#include <sstream>
//#include "curl/curl-7.83.0/include/curl/curl.h"

#include "msg.pb.h"
#include "testclient_mgr.h"
#include "channel.h"
#include "testsrv_mgr.h"
#include "cfg.h"

testclient_mgr * __tm_;
testsrv_mgr* __tsm_;

testcfg __global_cfg_ = {
	"118.190.144.92",
	//"47.104.129.217",
   //"192.168.2.129",
   2003,
   "10.0.14.48"
};


void test_cmd()
{
	while (1)
	{
		std::string line;
		std::getline(std::cin, line);

		if ("dump" == line)
		{
			printf("dump:");
			if (__tm_)
				__tm_->dump();

			if (__tsm_)
				__tsm_->dump();
		}
	}
}

int main(int argc, char* argv[])
{
	//CURLcode ret = curl_global_init(CURL_GLOBAL_ALL);
	//if (CURLE_OK != ret)
	//{
	//	printf("curl init fail");
	//	return 0;
	//}
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
		int thread_count = 50;
		int test_count = 1;
		int srv_count = 100;
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
	else // .\test.exe 0 1000 47.104.129.217 1000 10 50 
	{
		int64 id_base = 1000;
		int client_count = 1000;
		int test_count = 150;
		int thread_count = 50;
		std::string srvIP;
		if (argc >= 3)
			id_base = ::_atoi64(argv[2]);
		if (argc >= 4)
		{
			srvIP = argv[3];
			__global_cfg_.remote_ip_ = srvIP;
		}
		if (argc >= 5)
			client_count = atoi(argv[4]);
		if (argc >= 6)
			test_count = atoi(argv[5]);
		if (argc >= 7)
			thread_count = atoi(argv[6]);

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

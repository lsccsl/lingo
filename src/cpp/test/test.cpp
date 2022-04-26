#include <stdio.h>
#include <windows.h>
#include <string>
#include <iostream>
#include <sstream>

#include "msg.pb.h"
#include "testclient_mgr.h"
#include "channel.h"

testclient_mgr * __tm_;

void test_cmd()
{
	while (1)
	{
		std::string line;
		std::getline(std::cin, line);

		if ("dump" == line)
		{
			printf("dump:");
			int64 total_diff = 0;
			int64 total_count = 0;
			for (auto& it : __tm_->v_mgr_unit())
			{
				for (auto& itM : it.map_client_)
				{
					total_diff += itM.second->tc_static().total_diff_;
					total_count += itM.second->tc_static().total_count_;
				}
			}
			if (total_count <= 0)
				total_count = 1;
			double aver = (total_diff / 1000.f) / total_count;
			printf("total_count:%lld aver rtt:%f\r\n", total_count, aver);
		}
	}
}

int main(int argc, char* argv[])
{
	CChannel::init_sock();
	msgpacket::MSG_RPC_RES msg;
	msg.set_msg_id(123);
	printf("test:%s", msg.DebugString().c_str());
	
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST, msgpacket::_MSG_TEST);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST, msgpacket::_MSG_TEST);
	msgpackhelp::parse_reg(new msgpacket::MSG_LOGIN_RES, msgpacket::_MSG_LOGIN_RES);
	msgpackhelp::parse_reg(new msgpacket::MSG_TEST_RES, msgpacket::_MSG_TEST_RES);

	int client_count = 1;
	int64 id_base = 1000;
	int test_count = 1;
	if (argc >= 2)
		client_count = atoi(argv[1]);
	if (argc >= 3)
		id_base = ::_atoi64(argv[2]);
	if (argc >= 4)
		test_count = atoi(argv[3]);

	__tm_ = new testclient_mgr;

	__tm_->init(50, test_count);
	for (int i = 0; i < client_count; i++)
		__tm_->add_client(id_base + i);
	__tm_->run_thread();

	test_cmd();

	__tm_->join();
}

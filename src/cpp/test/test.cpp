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
			int64 max_diff = 0;
			int64 min_diff = INT64_MAX;
			int64 max_sendloop_interval = 0;
			int64 min_sendloop_interval = INT64_MAX;

			int64 total_sendloop_interval = 0;
			int64 total_send_loop = 0;

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
				}
			}
			if (total_send_loop < 1)
				total_send_loop = 1;
			int64 aver_total_send_loop = total_sendloop_interval / total_send_loop;
			if (total_count <= 0)
				total_count = 1;
			double aver = (total_diff / 1000.f) / total_count;
			printf("total_count:%lld aver rtt:%fs max_diff:%lldms, min_diff:%lldms\r\n"
				"max_sendloop_interval:%lldms min_sendloop_interval:%lldms aver sendloop:%lldms\r\n",
				total_count, aver, max_diff, min_diff,
				max_sendloop_interval, min_sendloop_interval, aver_total_send_loop);
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

	int64 id_base = 1000;
	int client_count = 1000;
	int test_count = 10;
	int thread_count = 50;
	if (argc >= 2)
		id_base = ::_atoi64(argv[1]);
	if (argc >= 3)
		client_count = atoi(argv[2]);
	if (argc >= 4)
		test_count = atoi(argv[3]);
	if (argc >= 5)
		thread_count = atoi(argv[4]);

	__tm_ = new testclient_mgr;

	__tm_->init(thread_count, test_count);
	for (int i = 0; i < client_count; i++)
		__tm_->add_client(id_base + i);
	__tm_->run_thread();

	test_cmd();

	__tm_->join();
}

#include <stdio.h>

#include "msg.pb.h"
#include "testclient_mgr.h"
#include "channel.h"

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
	if (argc >= 2)
		client_count = atoi(argv[1]);
	if (argc >= 3)
		id_base = ::_atoi64(argv[2]);

	testclient_mgr tm;

	tm.init(50);
	for (int i = 0; i < client_count; i++)
		tm.add_client(id_base + i);
	tm.run_thread();

	tm.join();
}

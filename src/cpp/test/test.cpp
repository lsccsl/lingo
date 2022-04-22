#include <stdio.h>

#include "msg.pb.h"

int main()
{
	msgpacket::MSG_RPC_RES msg;
	msg.set_msg_id(123);
	
	printf("test:%s", msg.DebugString().c_str());
}

#pragma once

#include <string>
#include "google/protobuf/message.h"
#include "type_def.h"

#pragma pack(1)
struct msghead
{
	uint32 pack_len;
	uint16 msg_type;
};
#pragma pack()

class msgpackhelp
{
public:

	static void pack_to_bin(std::string& buf_out, int msg_typ, google::protobuf::Message* proto_msg);
	static google::protobuf::Message* parse_from_bin(const void* buf, size_t buf_sz, const int msgtype);

	static void parse_reg(google::protobuf::Message* proto_msg, int msg_type);

private:

	typedef std::map<int/*msg_type*/, google::protobuf::Message*> MAP_MSGTYPE_PROTOMSG;
	static MAP_MSGTYPE_PROTOMSG map_msgtype_protomsg;
};

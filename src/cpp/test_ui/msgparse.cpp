#include "msgparse.h"
#include "type_def.h"
#include "msginter.pb.h"
#include "msgdef.pb.h"


msgpackhelp::MAP_MSGTYPE_PROTOMSG msgpackhelp::map_msgtype_protomsg;
void msgpackhelp::pack_to_bin(std::string& buf_out, int msg_typ, google::protobuf::Message* proto_msg)
{
	std::string msgbin;
	proto_msg->SerializeToString(&msgbin);

	msghead* mh;
	size_t buff_len = sizeof(msghead) + msgbin.size();
	buf_out.resize(buff_len);
	mh = (msghead*)(buf_out.data());
	//mh->pack_len = htonl((uint32)(6 + msgbin.size()));
	//mh->msg_type = htons(msg_typ);
	mh->pack_len = (uint32)(6 + msgbin.size());
	mh->msg_type = msg_typ;


	unsigned char* buf = (unsigned char*)mh + sizeof(msghead);
	memcpy(buf, msgbin.data(), msgbin.size());
}

google::protobuf::Message* msgpackhelp::parse_from_bin(const void* buf, size_t buf_sz, const int msgtype)
{
	auto it = msgpackhelp::map_msgtype_protomsg.find(msgtype);
	if (msgpackhelp::map_msgtype_protomsg.end() == it)
		return NULL;

	auto msg_org = it->second;
	if (!msg_org)
		return NULL;
	auto msg_ret = msg_org->New();
	if (!msg_ret)
		return NULL;

	if (buf_sz > 0)
		msg_ret->ParseFromArray(buf, buf_sz);

	return msg_ret;
}

void msgpackhelp::parse_reg(google::protobuf::Message* proto_msg, int msg_type)
{
	auto it = msgpackhelp::map_msgtype_protomsg.find(msg_type);
	if (msgpackhelp::map_msgtype_protomsg.end() != it)
		delete(it->second);
	msgpackhelp::map_msgtype_protomsg[msg_type] = proto_msg->New();
}

class _init_reg
{
public:

	_init_reg()
	{
		msgpackhelp::parse_reg(new msgpacket::PB_MSG_LOGON_RES, msgpacket::_PB_MSG_LOGON_RES);
	}
};

_init_reg __reg;




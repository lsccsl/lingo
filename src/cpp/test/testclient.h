#pragma once

#include <string>
#include <memory>
#include "channel.h"
#include "msg.pb.h"

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

	static void pack_to_bin(std::string& buf_out, int msg_typ, google::protobuf::Message * proto_msg);
	static google::protobuf::Message* parse_from_bin(const void * buf, size_t buf_sz, const msghead& mh);

	static void parse_reg(google::protobuf::Message* proto_msg, int msg_type);

private:

	typedef std::map<int/*msg_type*/, google::protobuf::Message*> MAP_MSGTYPE_PROTOMSG;
	static MAP_MSGTYPE_PROTOMSG map_msgtype_protomsg;
};

class testclient
{
public:

	testclient(const int64 id):id_(id)
	{}

	bool connect_to_srv(const std::string& srv_ip, int srv_port);
	const int64 id() const{
		return id_;
	}

	bool do_login();
	bool send_test(const int64 seq, int count);

	bool send_msg(int msg_typ, google::protobuf::Message* proto_msg);
	bool recv_one_msg();

private:

	void _reset_client()
	{
		this->lst_msg_recv_.clear();
		this->read_buf_sz_ = 0;
		this->read_buf_.resize(128);
	}

private:

	int32 fd_ = 0;
	int64 id_ = 0;

	std::string read_buf_;
	size_t read_buf_sz_ = 0;
	msghead mh_;

	struct ProtoMsg
	{
		std::shared_ptr<google::protobuf::Message> proto_msg;
		uint16 msg_type = 0;
	};
	std::list<ProtoMsg> lst_msg_recv_;
};

#pragma once

#include <string>
#include <memory>
#include <sys/timeb.h>
#include <stdint.h>
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

	struct testclient_static
	{
		int64 total_diff = 0;
		int64 total_count = 0;

		int64 max_diff = 0;
		int64 min_diff = INT64_MAX;

		int64 t_first_sendloop = 0;
		int64 t_last_sendloop = 0;
		int64 max_sendloop_interval = 0;
		int64 min_sendloop_interval = INT64_MAX;
		int64 total_send_loop = 0;
	};

	inline static int64 get_timestamp_mills()
	{
		timeb now;
		ftime(&now);
		return (now.time * 1000 + now.millitm);
	}

public:

	testclient(const int64 id):id_(id)
	{
		this->tc_static_.t_first_sendloop = this->tc_static_.t_last_sendloop = testclient::get_timestamp_mills();		
	}

	bool connect_to_srv(const std::string& srv_ip, int srv_port);
	inline const int64 id() const{
		return id_;
	}

	bool do_login();
	bool send_test(msgpacket::MSG_TEST& msg, const int64 seq, int count);
	bool recv_test(const int64 seq, int count);

	bool send_msg(int msg_typ, google::protobuf::Message* proto_msg);
	bool recv_one_msg();

public:

	testclient_static& tc_static() {
		return this->tc_static_;
	}

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

	testclient_static tc_static_;
};

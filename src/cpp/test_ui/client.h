#pragma once

#include <string>
#include "type_def.h"
#include "msginter.pb.h"
#include "msgparse.h"

struct client_static
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
	int last_read_err = 0;

	bool b_login_suc = false;
};

class client
{
public:

	client();
	~client();

	int connectToLogon(const std::string& ip, int port);
	int connectToGameSrv();



private:

	int send_msg(int msg_typ, google::protobuf::Message* proto_msg);
	int32 recv_one_msg();

	std::shared_ptr<google::protobuf::Message> get_msg_type(msgpacket::PB_MSG_TYPE msgtype);

private:

	int32 fd_ = -1;
	int64 id_ = 0;
	int32 magic_ = 0;

	int32 fd_gs_ = -1;
	int32 magic_gs_ = 0;
	std::string gs_ip_;
	int gs_port_;

	int64 client_id_ = 123;

	std::string read_buf_;
	size_t read_buf_sz_ = 0;
	msghead mh_;

	struct ProtoMsg
	{
		std::shared_ptr<google::protobuf::Message> proto_msg;
		uint16 msg_type = 0;
	};
	std::list<ProtoMsg> lst_msg_recv_;
	int need_recv_more_ = 0;


	client_static tc_static_;
};

#pragma once

#include <string>
#include <memory>
#include <sys/timeb.h>
#include <stdint.h>
#include <memory>
#include "channel.h"
#include "msg.pb.h"
#include "mylogex.h"
#include "testclient.h"

class testsrv
{
public:
	inline static int64 get_timestamp_mills()
	{
		timeb now;
		ftime(&now);
		return (now.time * 1000 + now.millitm);
	}
	static bool httpRequest(const int64 srvid, const std::string& ip, const int port);

public:

	testsrv(int64 srvid, const std::string& local_ip, const int local_port,
		const std::string& remote_ip, const int remote_port):
		srvid_(srvid),
		local_ip_(local_ip),
		local_port_(local_port),
		remote_ip_(remote_ip),
		remote_port_(remote_port)
	{
	}
	~testsrv();
	const int64 srvid()
	{ return this->srvid_; }

	void init_listen();

	void http_addsrv();
	void accept_client_no_block();
	bool connect_to_srv();

	bool recv_dial();
	template<class T>
	std::shared_ptr<T> recv_dial_msg(const int msgtype);

	bool recv_acpt();
	template<class T>
	std::shared_ptr<T> recv_acpt_msg(const int msgtype);

	bool send_msg_dial(int msg_typ, google::protobuf::Message* proto_msg);
	bool send_msg_acpt(int msg_typ, google::protobuf::Message* proto_msg);

	bool send_test_rpc(const int64 seq, const int64 timeout_wait);
	bool recv_test_rpc_res(const int64 seq);

	bool process_acpt_msg();

private:

	bool do_report();

private:
	int64 srvid_;
	std::string local_ip_;
	int local_port_;
	std::string remote_ip_;
	int remote_port_;

	struct DialInfo
	{
		std::string read_buf_;
		size_t read_buf_sz_ = 0;
		msghead mh_;

		struct ProtoMsg
		{
			std::shared_ptr<google::protobuf::Message> proto_msg;
			uint16 msg_type = 0;
		};
		std::list<ProtoMsg> lst_msg_recv_;

		int32 fd_dial_ = -1;
		int32 magic_dial_ = 0;

		int last_read_err = 0;

		bool b_login_suc = false;

		void _reset()
		{
			this->lst_msg_recv_.clear();
			this->read_buf_sz_ = 0;
			this->read_buf_.resize(128);
			this->magic_dial_ = 0;
			this->fd_dial_ = -1;
		}

	};
	DialInfo di_;

	struct AcptInfo
	{
		std::string read_buf_;
		size_t read_buf_sz_ = 0;
		msghead mh_;

		struct ProtoMsg
		{
			std::shared_ptr<google::protobuf::Message> proto_msg;
			uint16 msg_type = 0;
		};
		std::list<ProtoMsg> lst_msg_recv_;

		int32 fd_acpt_ = -1;
		int32 magic_acpt_ = 0;

		int last_read_err = 0;

		void _reset()
		{
			this->lst_msg_recv_.clear();
			this->read_buf_sz_ = 0;
			this->read_buf_.resize(128);
			this->magic_acpt_ = 0;
			this->fd_acpt_ = -1;
		}
	};
	AcptInfo ai_;


	int32 fd_lsn_ = -1;
};

template<class T>
std::shared_ptr<T> testsrv::recv_dial_msg(const int msgtype)
{
	std::shared_ptr<T> pret;
	if (!this->recv_dial())
	{
		MYLOG_ERR(("srvid:%lld login read err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return nullptr;
	}

	for (auto& it : this->di_.lst_msg_recv_)
	{
		if (it.msg_type == msgtype)
		{
			pret = std::dynamic_pointer_cast<T>(it.proto_msg);
			break;
		}
	}
	this->di_.lst_msg_recv_.clear();

	return pret;
}

template<class T>
std::shared_ptr<T> testsrv::recv_acpt_msg(const int msgtype)
{
	std::shared_ptr<T> pret;
	if (!this->recv_acpt())
	{
		MYLOG_ERR(("srvid:%lld login read err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return nullptr;
	}

	for (auto& it : this->ai_.lst_msg_recv_)
	{
		if (it.msg_type == msgtype)
		{
			pret = std::dynamic_pointer_cast<T>(it.proto_msg);
			break;
		}
	}
	this->ai_.lst_msg_recv_.clear();

	return pret;
}

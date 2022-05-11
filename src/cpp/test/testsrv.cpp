#include "testsrv.h"
#include <winsock2.h>
#include <string>
#include "msg.pb.h"
#include "curl/curl-7.83.0/include/curl/curl.h"

testsrv::~testsrv()
{}

void testsrv::init_listen()
{
	this->fd_lsn_ = CChannel::TcpOpen(local_ip_.c_str(), local_port_, 128);
	if (this->fd_lsn_ < 0)
		MYLOG_ERR(("srv:%lld, listen fail", this->srvid_));

	CChannel::set_no_block(this->fd_lsn_);
}

void testsrv::accept_client_no_block()
{
	testsrv::httpRequest(this->srvid_, this->local_ip_, this->local_port_);
	if (this->ai_.fd_acpt_ > 0)
		CChannel::CloseFd(this->ai_.fd_acpt_);
	this->ai_._reset();

	char actemp[32] = {};
	uint32 port = 0;
	this->ai_.fd_acpt_ = CChannel::TcpAccept(this->fd_lsn_, actemp, sizeof(actemp), &port);
	if (this->ai_.fd_acpt_ < 0)
	{
		MYLOG_ERR(("accept err srv:%lld read head err:%d-%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return;
	}

	auto msgReport = this->recv_acpt_msg<msgpacket::MSG_SRV_REPORT>(msgpacket::_MSG_SRV_REPORT);
	if (!msgReport)
	{
		MYLOG_ERR(("can't recv MSG_SRV_REPORT srv:%lld read head err:%d-%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError()));
	}
	MYLOG_DEBUG(("accept SUC srv:%lld", this->srvid_));
	msgpacket::MSG_SRV_REPORT_RES msgRes;
	msgRes.set_srv_id(this->srvid_);
	msgRes.set_tcp_conn_id(this->ai_.fd_acpt_);
	this->send_msg_acpt(msgpacket::_MSG_SRV_REPORT_RES, &msgRes);

	msgpacket::MSG_TEST_RPC msgTestRPC;
	msgTestRPC.set_rpc_count(100000000);
	this->send_msg_acpt(msgpacket::_MSG_TEST_RPC, &msgTestRPC);
}

bool testsrv::recv_dial()
{
	if (this->di_.read_buf_sz_ >= this->di_.read_buf_.size())
		this->di_.read_buf_.resize(this->di_.read_buf_sz_ * 2);

	void* buf = (void*)(this->di_.read_buf_.data() + this->di_.read_buf_sz_);

	int32 ret = 0;
	int read_sz = sizeof(msghead) - this->di_.read_buf_sz_;
	if (read_sz > 0)
		ret = CChannel::TcpSelectRead(this->di_.fd_dial_, buf, read_sz, 10, 30, &this->di_.last_read_err);
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld read head err:%d-%d read_sz:%d ret:%d magic:%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError(), this->di_.read_buf_sz_, ret, this->di_.magic_dial_));
		return false;
	}
	this->di_.read_buf_sz_ += ret;

	if (ret < sizeof(msghead))
		return true;
	msghead* pmh = (msghead*)this->di_.read_buf_.data();

	this->di_.mh_.pack_len = pmh->pack_len;
	this->di_.mh_.msg_type = pmh->msg_type;

	uint32 body_len = this->di_.mh_.pack_len - sizeof(msghead);
	ret = 0;
	if (body_len >= 0)
	{
		buf = (void*)(this->di_.read_buf_.data() + sizeof(msghead));
		ret = CChannel::TcpSelectRead(this->di_.fd_dial_, buf, body_len, 10, 30, &this->di_.last_read_err);
	}
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld read body err:%d-%d magic:%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError(), this->di_.magic_dial_));
		return false;
	}
	this->di_.read_buf_sz_ += ret;
	if (ret < body_len)
		return true;

	std::shared_ptr<google::protobuf::Message> proto_msg(msgpackhelp::parse_from_bin(this->di_.read_buf_.data() + sizeof(msghead), body_len, this->di_.mh_.msg_type));
	this->di_.lst_msg_recv_.push_back({ proto_msg, this->di_.mh_.msg_type });

	this->di_.read_buf_sz_ = 0;

	return true;
}
bool testsrv::recv_acpt()
{
	if (this->ai_.read_buf_sz_ >= this->ai_.read_buf_.size())
		this->ai_.read_buf_.resize(this->ai_.read_buf_sz_ * 2);

	void* buf = (void*)(this->ai_.read_buf_.data() + this->ai_.read_buf_sz_);

	int32 ret = 0;
	int read_sz = sizeof(msghead) - this->ai_.read_buf_sz_;
	if (read_sz > 0)
		ret = CChannel::TcpSelectRead(this->ai_.fd_acpt_, buf, read_sz, 10, 30, &this->ai_.last_read_err);
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld read head err:%d-%d read_sz:%d ret:%d magic:%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError(), this->ai_.read_buf_sz_, ret, this->ai_.magic_acpt_));
		return false;
	}
	this->ai_.read_buf_sz_ += ret;

	if (ret < sizeof(msghead))
		return true;
	msghead* pmh = (msghead*)this->ai_.read_buf_.data();

	this->ai_.mh_.pack_len = pmh->pack_len;
	this->ai_.mh_.msg_type = pmh->msg_type;

	uint32 body_len = this->ai_.mh_.pack_len - sizeof(msghead);
	ret = 0;
	if (body_len >= 0)
	{
		buf = (void*)(this->ai_.read_buf_.data() + sizeof(msghead));
		ret = CChannel::TcpSelectRead(this->ai_.fd_acpt_, buf, body_len, 10, 30, &this->ai_.last_read_err);
	}
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld read body err:%d-%d magic:%d",
			this->srvid_, ::WSAGetLastError(), ::GetLastError(), this->ai_.magic_acpt_));
		return false;
	}
	this->ai_.read_buf_sz_ += ret;
	if (ret < body_len)
		return true;

	std::shared_ptr<google::protobuf::Message> proto_msg(msgpackhelp::parse_from_bin(this->ai_.read_buf_.data() + sizeof(msghead), body_len, this->ai_.mh_.msg_type));
	this->ai_.lst_msg_recv_.push_back({ proto_msg, this->ai_.mh_.msg_type });

	this->ai_.read_buf_sz_ = 0;

	return true;
}

bool testsrv::send_msg_dial(int msg_typ, google::protobuf::Message* proto_msg)
{
	std::string buf_bin;
	msgpackhelp::pack_to_bin(buf_bin, msg_typ, proto_msg);

	int32 ret = CChannel::TcpSelectWrite(this->di_.fd_dial_, buf_bin.data(), buf_bin.size(), 10, 30);
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld write err:%d-%d ret:%d magic:%d", this->srvid_, ::WSAGetLastError(), ::GetLastError(), ret, this->di_.magic_dial_));
		return false;
	}
	return true;
}
bool testsrv::send_msg_acpt(int msg_typ, google::protobuf::Message* proto_msg)
{
	std::string buf_bin;
	msgpackhelp::pack_to_bin(buf_bin, msg_typ, proto_msg);

	int32 ret = CChannel::TcpSelectWrite(this->ai_.fd_acpt_, buf_bin.data(), buf_bin.size(), 10, 30);
	if (ret < 0)
	{
		MYLOG_ERR(("srv:%lld write err:%d-%d ret:%d magic:%d", this->srvid_, ::WSAGetLastError(), ::GetLastError(), ret, this->ai_.magic_acpt_));
		return false;
	}
	return true;
}

bool testsrv::do_report()
{
	msgpacket::MSG_SRV_REPORT msg;
	msg.set_srv_id(this->srvid_);

	if (!this->send_msg_dial(msgpacket::_MSG_SRV_REPORT, &msg))
	{
		MYLOG_ERR(("srv:%lld login write err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}
	if (!this->recv_dial())
	{
		MYLOG_ERR(("srv:%lld login read err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}

	bool bret = false;
	for (auto& it : this->di_.lst_msg_recv_)
	{
		if (it.msg_type == msgpacket::_MSG_SRV_REPORT_RES)
		{
			msgpacket::MSG_LOGIN_RES* msgLoginRes = dynamic_cast<msgpacket::MSG_LOGIN_RES*>(it.proto_msg.get());
			if (msgLoginRes)
				this->di_.magic_dial_ = msgLoginRes->connect_id();
			bret = true;
			break;
		}
	}
	this->di_.lst_msg_recv_.clear();

	return bret;
}

void testsrv::http_addsrv()
{
	testsrv::httpRequest(this->srvid_, this->local_ip_, this->local_port_);
}

bool testsrv::connect_to_srv()
{
	testsrv::httpRequest(this->srvid_, this->local_ip_, this->local_port_);

	this->di_.b_login_suc = false;
	if (this->di_.fd_dial_ > 0)
		CChannel::CloseFd(this->di_.fd_dial_);
	this->di_._reset();

	this->di_.fd_dial_ = CChannel::TcpConnect(this->remote_ip_.c_str(), this->remote_port_, 10, 30);
	if (this->di_.fd_dial_ < 0)
	{
		MYLOG_ERR(("srv:%lld connect err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}

	CChannel::keep_alive(this->di_.fd_dial_);
	CChannel::set_no_block(this->di_.fd_dial_);
	bool bret = this->do_report();
	if (bret)
	{
		this->di_.b_login_suc = true;
		MYLOG_ERR(("srv:%lld connect suc, fd:%d, magic:%d", this->srvid_, this->di_.fd_dial_, this->di_.magic_dial_));
	}
	else
	{
		MYLOG_ERR(("srv:%lld send login err:%d-%d", this->srvid_, ::WSAGetLastError(), ::GetLastError()));
	}

	return bret;
}

bool testsrv::send_test_rpc(const int64 seq, const int64 timeout_wait)
{
	msgpacket::MSG_RPC msgRPC;
	msgRPC.set_msg_id(seq);
	msgRPC.set_msg_type(msgpacket::_MSG_TEST);
	msgRPC.set_timestamp(testsrv::get_timestamp_mills());
	msgRPC.set_timeout_wait(timeout_wait);

	msgpacket::MSG_TEST msgTest;
	msgTest.set_id(this->srvid_);
	msgTest.set_seq(seq);

	std::string strBin;
	msgRPC.SerializeToString(&strBin);
	msgRPC.set_msg_bin(strBin);

	return send_msg_dial(msgpacket::_MSG_RPC, &msgRPC);
}
bool testsrv::recv_test_rpc_res(const int64 seq)
{
	auto pMsgRes = this->recv_dial_msg<msgpacket::MSG_RPC_RES>(msgpacket::_MSG_RPC_RES);
	if (!pMsgRes)
		return false;
	return true;
}

bool testsrv::process_acpt_msg()
{
	bool bret = this->recv_acpt();
	for (auto& it : this->ai_.lst_msg_recv_)
	{
		switch (it.msg_type)
		{
		case msgpacket::_MSG_RPC:
			{
				//MYLOG_DEBUG(("recv rpc"));
				std::shared_ptr<msgpacket::MSG_RPC> pMsg = std::dynamic_pointer_cast<msgpacket::MSG_RPC>(it.proto_msg);
				if (pMsg)
				{
					msgpacket::MSG_RPC_RES msgRes;
					msgRes.set_msg_id(pMsg->msg_id());
					msgRes.set_msg_type(msgpacket::_MSG_TEST_RES);

					std::shared_ptr<google::protobuf::Message> proto_msg(msgpackhelp::parse_from_bin(pMsg->msg_bin().data(), pMsg->msg_bin().size(), pMsg->msg_type()));
					std::shared_ptr<msgpacket::MSG_TEST> pTest = std::dynamic_pointer_cast<msgpacket::MSG_TEST>(proto_msg);
					if (pTest)
					{
						msgpacket::MSG_TEST_RES msgTestRes;
						msgTestRes.set_id(pTest->id());
						msgTestRes.set_seq(pTest->seq());
						msgTestRes.set_str("````msgTest.Str!!!!");
						std::string strbin;
						msgTestRes.SerializeToString(&strbin);
						msgRes.set_msg_bin(strbin);
					}
					this->send_msg_acpt(msgpacket::_MSG_RPC_RES, &msgRes);
				}
			}
			break;
		case msgpacket::_MSG_HEARTBEAT:
			{
				msgpacket::MSG_HEARTBEAT_RES msgRes;
				std::shared_ptr<msgpacket::MSG_HEARTBEAT> pMsg = std::dynamic_pointer_cast<msgpacket::MSG_HEARTBEAT>(it.proto_msg);
				if (pMsg)
				{
					msgRes.set_id(pMsg->id());
				}
				this->send_msg_acpt(msgpacket::_MSG_HEARTBEAT_RES, &msgRes);
			}
			break;
		case msgpacket::_MSG_SRV_REPORT:
			{
				std::shared_ptr<msgpacket::MSG_SRV_REPORT> pMsg = std::dynamic_pointer_cast<msgpacket::MSG_SRV_REPORT>(it.proto_msg);
				if (pMsg)
				{
					msgpacket::MSG_SRV_REPORT_RES msgRes;
					msgRes.set_srv_id(pMsg->srv_id());
					msgRes.set_tcp_conn_id(this->ai_.fd_acpt_);
					this->send_msg_acpt(msgpacket::_MSG_SRV_REPORT_RES, &msgRes);
				}
			}
			break;
		}
	}
	this->ai_.lst_msg_recv_.clear();

	return bret;
}

size_t req_reply(void* ptr, size_t size, size_t nmemb, void* stream)
{
	std::string* str = (std::string*)stream;
	(*str).append((char*)ptr, size * nmemb);
	return size * nmemb;
}
bool testsrv::httpRequest(const int64 srvid, const std::string& ip, const int port)
{
	CURL* curl = curl_easy_init();
	CURLcode ret = CURLE_OK;

	std::string str;
	str = str + "{\"SrvID\":" + std::to_string(srvid) + ",\"IP\":\"" + ip + "\",\"Port\":" + std::to_string(port) + "}";

	curl_slist* header_list = NULL;
	header_list = curl_slist_append(header_list, "Content-Type: application/text");
	ret = curl_easy_setopt(curl, CURLOPT_HTTPHEADER, header_list);

	curl_easy_setopt(curl, CURLOPT_POST, 1);
	curl_easy_setopt(curl, CURLOPT_URL, "http://192.168.2.129:8803/addserver");
	curl_easy_setopt(curl, CURLOPT_POSTFIELDS, str.c_str());

	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, req_reply);
	std::string response;
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void*)&response);
	ret = curl_easy_perform(curl);
	if (CURLE_OK != ret)
		return false;
	curl_easy_cleanup(curl);
	return true;
}

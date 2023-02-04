#include "client.h"
#include "channel.h"
#include "mylogex.h"
#include "os.h"
#ifdef WIN32
#include <winsock2.h>
#endif

client::client()
{
	this->read_buf_.resize(128);
}

client::~client()
{}

int client::connectToLogon(const std::string& ip, int port)
{
	if (this->fd_ > 0)
		CChannel::CloseFd(this->fd_);
	this->fd_ = -1;
	this->fd_ = CChannel::TcpConnect(ip.c_str(), port, 10, 30);
	if (this->fd_ < 0)
	{
		MYLOG_ERR(("clientid:%lld connect err:%d-%d", this->id_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}

	CChannel::set_no_block(this->fd_);

	msgpacket::PB_MSG_LOGON msgLogon;
	msgLogon.set_client_id(client_id_);
	this->send_msg(msgpacket::_PB_MSG_LOGON, &msgLogon);

	this->recv_one_msg();

	return 0;
}

int32 client::send_msg(int msg_typ, google::protobuf::Message* proto_msg)
{
	std::string buf_bin;
	msgpackhelp::pack_to_bin(buf_bin, msg_typ, proto_msg);

	int32 ret = CChannel::TcpSelectWrite(this->fd_, buf_bin.data(), buf_bin.size(), 10, 30);
	if (ret < 0)
	{
		MYLOG_ERR(("clientid:%lld write err:%d-%d ret:%d magic:%d", this->id_, ::WSAGetLastError(), ::GetLastError(), ret, this->magic_));
		return false;
	}
	return true;
}

int32 client::recv_one_msg()
{
	if (read_buf_sz_ >= this->read_buf_.size())
		this->read_buf_.resize(read_buf_sz_ * 2);

	void* buf = (void*)(this->read_buf_.data() + this->read_buf_sz_);

	int32 ret = 0;
	int read_sz = sizeof(msghead) - this->read_buf_sz_;
	if (read_sz > 0)
		ret = CChannel::TcpSelectRead(this->fd_, buf, read_sz, 10, 30, &this->tc_static_.last_read_err);
	if (ret < 0)
	{
		MYLOG_ERR(("clientid:%lld read head err:%d-%d read_sz:%d ret:%d magic:%d",
			this->id_, ::WSAGetLastError(), ::GetLastError(), this->read_buf_sz_, ret, this->magic_));
		return false;
	}
	this->read_buf_sz_ += ret;

	if (ret < sizeof(msghead))
		return true;
	msghead* pmh = (msghead*)this->read_buf_.data();

	this->mh_.pack_len = pmh->pack_len;
	this->mh_.msg_type = pmh->msg_type;

	uint32 body_len = this->mh_.pack_len - sizeof(msghead);
	ret = 0;
	if (body_len >= 0)
	{
		buf = (void*)(this->read_buf_.data() + sizeof(msghead));
		ret = CChannel::TcpSelectRead(this->fd_, buf, body_len, 10, 30, &this->tc_static_.last_read_err);
	}
	if (ret < 0)
	{
		MYLOG_ERR(("clientid:%lld read body err:%d-%d magic:%d",
			this->id_, ::WSAGetLastError(), ::GetLastError(), this->magic_));
		return false;
	}
	this->read_buf_sz_ += ret;
	if (ret < body_len)
		return true;

	std::shared_ptr<google::protobuf::Message> proto_msg(msgpackhelp::parse_from_bin(this->read_buf_.data() + sizeof(msghead), body_len, this->mh_.msg_type));
	client::ProtoMsg pm;
	pm.proto_msg = proto_msg;
	pm.msg_type = this->mh_.msg_type;
	this->lst_msg_recv_.push_back(pm);

	this->read_buf_sz_ = 0;

	return true;
}

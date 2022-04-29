#include "testclient.h"
#include <string>
#include <winsock2.h>
#include "mylogex.h"

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

google::protobuf::Message* msgpackhelp::parse_from_bin(const void* buf, size_t buf_sz, const msghead& mh)
{
	auto it = msgpackhelp::map_msgtype_protomsg.find(mh.msg_type);
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


bool testclient::connect_to_srv(const std::string& srv_ip, int srv_port)
{
	this->_reset_client();
	if (this->fd_ > 0)
		CChannel::CloseFd(this->fd_);
	this->fd_ = CChannel::TcpConnect(srv_ip.c_str(), srv_port, 10, 30);
	if (this->fd_ < 0)
	{
		MYLOG_ERR(("clientid:%lld connect err:%d-%d", this->id_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}

	CChannel::keep_alive(this->fd_);

	return true;
}

bool testclient::do_login()
{
	msgpacket::MSG_LOGIN msg;
	msg.set_id(this->id_);

	if (!this->send_msg(msgpacket::_MSG_LOGIN, &msg))
	{
		MYLOG_ERR(("clientid:%lld login write err:%d-%d", this->id_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}
	if (!this->recv_one_msg())
	{
		MYLOG_ERR(("clientid:%lld login read err:%d-%d", this->id_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}

	bool bret = false;
	for (auto& it : this->lst_msg_recv_)
	{
		if (it.msg_type == msgpacket::_MSG_LOGIN_RES)
		{
			bret = true;
			break;
		}
	}
	this->lst_msg_recv_.clear();

	return bret;
}
bool testclient::send_test(msgpacket::MSG_TEST& msg, const int64 seq, int count)
{
	this->tc_static_.total_send_loop++;
	int64 tnowMS = testclient::get_timestamp_mills();
	if (this->tc_static_.t_last_sendloop > 0)
	{
		int64 diff = tnowMS - this->tc_static_.t_last_sendloop;

		if (diff > this->tc_static_.max_sendloop_interval)
			this->tc_static_.max_sendloop_interval = diff;
		if (diff < this->tc_static_.min_sendloop_interval)
			this->tc_static_.min_sendloop_interval = diff;
	}
	this->tc_static_.t_last_sendloop = tnowMS;

	for (int i = 0; i < count; i++)
	{
		msg.set_timestamp(testclient::get_timestamp_mills());
		msg.set_seq(seq + i);
		if (!this->send_msg(msgpacket::_MSG_TEST, &msg))
			return false;
	}

	return true;
}

bool testclient::recv_test(const int64 seq, int count)
{
	for (int i = 0; i < count; i++)
	{
		if (!this->recv_one_msg())
			return false;

		if (!this->lst_msg_recv_.empty())
		{
			int64 tnow_mills = testclient::get_timestamp_mills();

			std::list<ProtoMsg>::reverse_iterator it = this->lst_msg_recv_.rbegin();
			if (it->msg_type == msgpacket::_MSG_TEST_RES)
			{
				auto msgRes = std::dynamic_pointer_cast<msgpacket::MSG_TEST_RES>(it->proto_msg);
				if (msgRes->seq() == (seq + i))
				{
					int64 diff = (tnow_mills - msgRes->timestamp());
					this->tc_static_.total_diff += diff;
					this->tc_static_.total_count++;

					if (diff > this->tc_static_.max_diff)
						this->tc_static_.max_diff = diff;
					if (diff < this->tc_static_.min_diff)
						this->tc_static_.min_diff = diff;
				}
			}
		}
	}

	bool bret = false;
	for (std::list<ProtoMsg>::reverse_iterator it = this->lst_msg_recv_.rbegin(); it != this->lst_msg_recv_.rend(); it ++)
	{
		if (it->msg_type == msgpacket::_MSG_TEST_RES)
		{
			auto msgRes = std::dynamic_pointer_cast<msgpacket::MSG_TEST_RES>(it->proto_msg);
			if (msgRes->seq() == (seq + count - 1))
			{
				bret = true;
				break;
			}
		}
	}
	this->lst_msg_recv_.clear();
	if (!bret)
		MYLOG_ERR(("id:%lld recv seq:%lld err", this->id_, seq));

	return true;
}

bool testclient::send_msg(int msg_typ, google::protobuf::Message* proto_msg)
{
	std::string buf_bin;
	msgpackhelp::pack_to_bin(buf_bin, msg_typ, proto_msg);

	int32 ret = CChannel::TcpSelectWrite(this->fd_, buf_bin.data(), buf_bin.size(), 10, 30);
	if (ret < 0)
	{
		MYLOG_ERR(("clientid:%lld write err:%d-%d ret:%d", this->id_, ::WSAGetLastError(), ::GetLastError(), ret));
		return false;
	}
	return true;
}

bool testclient::recv_one_msg()
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
		MYLOG_ERR(("clientid:%lld read head err:%d-%d read_sz:%d ret:%d", this->id_, ::WSAGetLastError(), ::GetLastError(), this->read_buf_sz_, ret));
		return false;
	}
	this->read_buf_sz_ += ret;

	if (ret < sizeof(msghead))
		return true;
	msghead *pmh = (msghead*)this->read_buf_.data();

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
		MYLOG_ERR(("clientid:%lld read body err:%d-%d", this->id_, ::WSAGetLastError(), ::GetLastError()));
		return false;
	}
	this->read_buf_sz_ += ret;
	if (ret < body_len)
		return true;

	std::shared_ptr<google::protobuf::Message> proto_msg(msgpackhelp::parse_from_bin(this->read_buf_.data() + sizeof(msghead), body_len, this->mh_));
	this->lst_msg_recv_.push_back({ proto_msg, this->mh_.msg_type });

	this->read_buf_sz_ = 0;

	return true;
}

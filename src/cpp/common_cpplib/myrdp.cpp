/**
 * @file myrdp.cpp 
 * @brief 可靠udp
 * the udp head like this,
 *
 * request:
 * ver          int8    rdp协议版本号
 * rdp_id       uint64  报文唯一id 32(host_id)+8(reserver:must be 0x00)+24(rdpSeq)
 * rdp_info     uint8   bit1: 1:请求 0:响应
 * block_id     uint16  块的 id
 * block_count  uint16  块的个数(id=0时才有此域)
 *
 * response:
 * ver          int8    rdp协议版本号
 * rdp_id       uint64  报文唯一id
 * rdp_info     uint8   bit1: 1:请求 0:响应
 * block_id     uint16  块id
 *
 * @author linshaochuan
 */
#include "myrdp.h"
#include <string.h>
#include "channel.h"
#include "PacketBase.h"
#include "mylogex.h"
#define MYRDP_LOG_PRE "myrdp"

class rdpPkt : public PacketBase
{
public:

	rdpPkt():ver_(0),rdp_id_(0),rdp_info_(0),block_id_(0),block_count_(0),data_(NULL),data_sz_(0)
	{}

	~rdpPkt()
	{}

	/* @brief 打包 */
	void pack()
	{
		this->Add1Byte(this->ver_);
		this->Add8Byte(this->rdp_id_);
		this->Add1Byte(this->rdp_info_);
		this->Add2Byte(this->block_id_);
		if(this->rdp_info_ & 0x01)
		{
			if(0 == this->block_id_)
				this->Add2Byte(this->block_count_);

			if(data_ && data_sz_)
				this->AddXByte((uint8 *)data_, data_sz_);
		}
	}

	/* @brief 解包 */
	void parse(const void * data, const uint32 data_sz)
	{
		this->EatBuf((uint8 *)data, data_sz);

		this->Get1Byte(this->ver_);
		this->Get8Byte(this->rdp_id_);
		this->Get1Byte(this->rdp_info_);
		this->Get2Byte(this->block_id_);
		if(this->rdp_info_ & 0x01)
		{
			if(0 == this->block_id_)
			{
				this->Get2Byte(this->block_count_);
				assert(data_sz > 14);
				vdata_.resize(data_sz - 14);
			}
			else
			{
				assert(data_sz > 12);
				vdata_.resize(data_sz - 12);
			}

			this->GetXByte((uint8 *)&vdata_[0], vdata_.size());
		}
	}

public:

	/* @brief rdp协议版本号 */
	uint8 ver_;

	/* @brief 报文唯一id */
	uint64 rdp_id_;

	/* @brief bit: 0:请求 1:响应 */
	uint8 rdp_info_;

	/* @brief 块的 id */
	uint16 block_id_;

	/* @brief 块的个数(id=0时才有此域) */
	uint16 block_count_;

	/* @brief 数据 */
	uint8 * data_;
	uint32 data_sz_;

	/* @brief 收到的数据 */
	std::vector<uint8> vdata_;
};

/* @brief 消息类型 */
enum RDP_MSG_TYPE
{
    RDP_MSG_SEND = 0,
};

/* @brief 消息头 */
struct rdp_msg_head
{
	rdp_msg_head(RDP_MSG_TYPE type):type_(type)
	{}

	virtual ~rdp_msg_head(){}

    RDP_MSG_TYPE type_;
};

/* @brief 发送数据请求 */
struct rdp_msg_send : public rdp_msg_head
{
	rdp_msg_send():rdp_msg_head(RDP_MSG_SEND)
	{}

    /* @brief 上下文 */
    unsigned long ctx_;
    
    /* @brief 要发送的数据 */
    std::vector<uint8> data_;

    /* @brief 发送的目标 */
    int8 dst_ip_[16];
    uint16 dst_port_;

    /* @brief 重发次数 */
    uint32 resend_cout_;
    /* @brief 重发间隔 tick */
    uint32 resend_interval_;
};

/**
* @brief constructor
*/
myrdp::myrdp(const char * local_ip, const uint16 port, const uint32 host_id,
	const uint32 mtu, const uint32 max_rcv_time_out):host_id_(host_id),mtu_(mtu),rdp_seq_(0),time_tick_(0),htimer_(NULL)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "open rdp:[%s:%d] host id:%d mtu:%d", local_ip, port, host_id, mtu));

	assert(local_ip);

	this->udp_fd_ = CChannel::UdpOpen(local_ip, port);

	max_recv_timeout_tick_ = max_rcv_time_out / 100;
	if(0 == max_recv_timeout_tick_)
		max_recv_timeout_tick_ = 3;
}

/**
* @brief destructor
*/
myrdp::~myrdp()
{
	MyListernerDestruct(this->lsn_);
	this->lsn_ = NULL;

	CChannel::CloseFd(this->udp_fd_);
	this->udp_fd_ = -1;
}

/**
* @brief init rdp
*/
int32 myrdp::init()
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::init"));

	if(this->udp_fd_ < 0)
	{
		MYLOG_INFO(("can't open udp fd ..."));
		return -1;
	}

	this->lsn_ = MyListernerConstruct(NULL, 65535);
	if(NULL == this->lsn_)
	{
		MYLOG_INFO(("can't create listerner ..."));
		return -1;
	}

	event_handle_t eh = {0};

	eh.input = myrdp::_rdp_handle_input;
	eh.output = myrdp::_rdp_handle_output;
	eh.exception = myrdp::_rdp_handle_exception;
	eh.context_data = (unsigned long)this;

	CChannel::set_no_block(this->udp_fd_);

	if(0 != MyListernerAddFD(this->lsn_, this->udp_fd_, E_FD_READ, &eh))
	{
		MYLOG_INFO(("add udp fd to listerner fail ..."));
		return -1;
	}

	mytimer_node_t n = {0};
	n.context_data = (unsigned long)this;
	n.timer_user_data = (unsigned long)NULL;
	n.timeout_cb = myrdp::_rdp_timeout;
	n.period.tv_sec = 1;
	n.period.tv_usec = 0;
	n.first_expire.tv_sec = 1;
	n.first_expire.tv_usec = 0;

	this->time_resolution_ = 10;
	this->htimer_ = MyListernerAddTimer(this->lsn_, &n);

	MyListernerRun(this->lsn_);

	return 0;
}

/**
* @brief send data
* @param resend_cout:重发次数
* @param resend_interval:重发间隔(毫秒,千分之一秒)
*/
int32 myrdp::send_data(const unsigned long ctx,
	std::vector<uint8>& data,
	const int8 * dst_ip, const uint16 dst_port,
	const uint32 resend_cout, const uint32 resend_interval)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::send_data data_sz:%d resend_cout:%d resend_interval:%d [%s:%d]", data.size(), resend_cout, resend_interval,
		dst_ip ? dst_ip : "null", dst_port));

	rdp_msg_send * m = new rdp_msg_send;

	m->ctx_ = ctx;
    
    m->data_.swap(data);

	strncpy(m->dst_ip_, dst_ip, sizeof(m->dst_ip_));
	m->dst_port_ = dst_port;

    m->resend_cout_ = resend_cout;
    m->resend_interval_ = resend_interval;

	
	MyListernerAddMsg(this->lsn_, m, (unsigned long)this, myrdp::_rdp_handle_msg);
	return 0;
}

/**
* @brief send data
* @param resend_cout:重发次数
* @param resend_interval:重发间隔(毫秒,千分之一秒)
*/
int32 myrdp::_send_data(const unsigned long ctx, 
	const uint8 * data, const uint32 data_sz,
	const int8 * dst_ip, const uint16 dst_port,
	const uint32 resend_cout, const uint32 resend_interval)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::send_data data:%x data_sz:%d resend_cout:%d resend_interval:%d [%s:%d]", data, data_sz, resend_cout, resend_interval,
		dst_ip ? dst_ip : "null", dst_port));

	if(NULL == data || 0 == data_sz || NULL == dst_ip)
	{
		MYLOG_INFO(("err param ... "));
		return -1;
	}

	send_session_t * ss;
	pkt_frame_t * p;
	uint32 pos = 0;
	uint16 block_id = 0;
	uint64 rdp_id;
	uint16 block_count = (data_sz / this->mtu_) + ((data_sz % this->mtu_) ? 1 : 0);

	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "block_count:%d", block_count));

	if(resend_interval < 1000)
	{
		this->_shift_timer_to_100();
	}

	ss = new send_session_t;
	strncpy(ss->dst_ip_, dst_ip, sizeof(ss->dst_ip_));
	ss->dst_port_ = dst_port;
	ss->resend_cout_ = resend_cout;
	ss->resend_interval_tick_ = (resend_interval / 100);
	if(0 == ss->resend_interval_tick_)
		ss->resend_interval_tick_ = 1;

	ss->last_ack_timetick_ = ss->last_send_tick_ = this->time_tick_;

	ss->ctx_ = ctx;

	this->_gen_rdp_id(rdp_id);
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "rdp id:%llx", rdp_id));

	map_send_session_[rdp_id] = ss;

	/* 第一次只发10个报文分片 */
	uint32 first_send_count = 0;

	while(pos < data_sz)
	{
		rdpPkt rp;

		rp.rdp_id_ = rdp_id;

		rp.rdp_info_ = 0x01;
		rp.block_id_ = block_id;
		block_id ++;
		rp.block_count_ = block_count;

		rp.data_ = (uint8 *)(data + pos);

		if(rp.block_id_ < (block_count - 1))
			rp.data_sz_ = this->mtu_;
		else
			rp.data_sz_ = data_sz % this->mtu_;

		MYLOG_DEBUGEX((MYRDP_LOG_PRE, "pos:%d frame sz:%d block id:%d", pos, rp.data_sz_, rp.block_id_));

		pos += rp.data_sz_;

		rp.pack();
		p = new pkt_frame_t;
		p->pkt.swap((std::vector<uint8>&)(rp.GetBufRef()));

		ss->vpkt_frame_.push_back(p);

		MYLOG_DUMP_BIN(&p->pkt[0], p->pkt.size());

		if(first_send_count <= ss->last_send_speed_)
		{
			int32 ret = CChannel::UdpWrite(this->udp_fd_, &p->pkt[0], p->pkt.size(), ss->dst_ip_, ss->dst_port_);
			first_send_count ++;

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "udp send to[ %s:%d] first_send_count:%d ret:%d",
				ss->dst_ip_, ss->dst_port_, first_send_count, ret));
		}
		else
		{
			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "reach send max:%d", first_send_count));
		}
	}

	return 0;
}

/**
* @brief 发送回调
* @param result: 0:成功 其它:失败
*/
void myrdp::send_data_res(const unsigned long ctx, unsigned int result)
{
	MYLOG_INFO(("myrdp::send_data_res ctx:%x result:%d", ctx, result));
}

/**
* @brief 收到数据回调函(不可阻塞)
*/
void myrdp::recv_data(const uint8 * data, const uint32 data_sz, const int8 * src_ip, const uint16 src_port)
{
	MYLOG_INFO(("recv data from[%s:%d] %d byte, over write me please ...", src_ip, src_port, data_sz));
}

/* @brief 生成rdp id */
void myrdp::_gen_rdp_id(uint64& rdp_id)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_gen_rdp_id"));

	rdp_id = this->host_id_;
	rdp_id = (rdp_id << 32) & 0xffffffff00000000LL;

	if(this->rdp_seq_ > 0xffffff)
		this->rdp_seq_ = 0;
	rdp_id += this->rdp_seq_;

	this->rdp_seq_ ++;

	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "rdp_id:%llx", rdp_id));
}

/* @brief 处理有输入事件的回调函数 */
int myrdp::_rdp_handle_input(unsigned long context_data, int fd)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_rdp_handle_input context_data:%x fd:%d", context_data, fd));

	myrdp * rdp = (myrdp *)context_data;

	std::vector<uint8> vbuf(myrdp::_rdp_head_len_ + rdp->mtu_);
	int8 src_ip[16];
	uint16 src_port;

	do
	{
		int32 ret = CChannel::UdpRead(rdp->udp_fd_, &vbuf[0], vbuf.size(), src_ip, sizeof(src_ip), &src_port);

		if(ret <= 0)
			break;

		vbuf.resize(ret);

		MYLOG_DUMP_BIN(&vbuf[0], vbuf.size());

		/* 处理 */
		rdpPkt p;
		p.parse(&vbuf[0], vbuf.size());

		MYLOG_DEBUGEX((MYRDP_LOG_PRE, "block id:%d rdp id:%llx", p.block_id_, p.rdp_id_));
		if(!(p.rdp_info_ & 0x01))
		{
			/* 收到res报文,与发送session里的分片比对,并结束发送session */
			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "res"));

			std::map<uint64, send_session_t *>::iterator it = rdp->map_send_session_.find(p.rdp_id_);
			if(it == rdp->map_send_session_.end())
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "session has been destroy ..."));
				continue;
			}

			send_session_t * ss = it->second;

			if(p.block_id_ > ss->vpkt_frame_.size())
			{
				MYLOG_INFO(("block is wrong block id:%d rdp id:%llx", p.block_id_, p.rdp_id_));
				continue;
			}

			delete ss->vpkt_frame_[p.block_id_];
			ss->vpkt_frame_[p.block_id_] = NULL;

			ss->ack_count_ ++;

			if(ss->ack_count_ < ss->vpkt_frame_.size())
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "the pkt still have sub frame"));
				continue;
			}

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "all sub frame have been send"));

			/* notify up layer app */
			rdp->send_data_res(ss->ctx_, SEND_OK);

			/* release session */
			delete ss;
			rdp->map_send_session_.erase(it);

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "send session count:%d", rdp->map_send_session_.size()));
		}
		else
		{
			/* 收到req报文,启动接收session */
			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "req"));

			recv_session_t * rs = NULL;
			std::map<uint64, recv_session_t *>::iterator it = rdp->map_recv_session_.find(p.rdp_id_);
			if(it == rdp->map_recv_session_.end())
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "start recv session"));
				rs = rdp->map_recv_session_[p.rdp_id_] = new recv_session_t(rdp->time_tick_ + rdp->max_recv_timeout_tick_);
			}
			else
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "recv session exist already"));
				rs = it->second;
			}

			if(rs->vpkt_frame_.size() <= p.block_id_)
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "block id beyond the vector size, resize it"));
				rs->vpkt_frame_.resize(p.block_id_ + 1);
			}

			MYLOG_DUMP_BIN(&(p.vdata_[0]), p.vdata_.size());
			/* 保存报文内容 */
			rs->vpkt_frame_[p.block_id_].pkt.swap(p.vdata_);

			if(0 == p.block_id_)
			{
				rs->b_rcv_first_ = 1;
				rs->total_block_ = p.block_count_;
			}

			rs->cur_rcv_block_ ++;

			/* 同时回复res */
			rdpPkt pr;
			pr.rdp_id_ = p.rdp_id_;
			pr.rdp_info_ = 0x00;
			pr.block_id_ = p.block_id_;

			pr.pack();
			CChannel::UdpWrite(rdp->udp_fd_, &pr.GetBufRef()[0], pr.GetBufRef().size(), src_ip, src_port);

			/* 是否收全 */
			if(!rs->b_rcv_first_)
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "first block not recv yet ..."));
				continue;
			}

			if(rs->cur_rcv_block_ < rs->total_block_)
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "not all block has been recved ..."));
				continue;
			}

			/* 组合数据,通知 */
			for(uint32 i = 1; i < rs->vpkt_frame_.size(); i ++)
			{
				rs->vpkt_frame_[0].pkt.insert(rs->vpkt_frame_[0].pkt.end(),
					rs->vpkt_frame_[i].pkt.begin(), rs->vpkt_frame_[i].pkt.end());
			}

			MYLOG_DUMP_BIN(&(rs->vpkt_frame_[0].pkt[0]), rs->vpkt_frame_[0].pkt.size());

			rdp->recv_data(&(rs->vpkt_frame_[0].pkt[0]), rs->vpkt_frame_[0].pkt.size(), src_ip, src_port);

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "release rdp session %llx", p.rdp_id_));

			rdp->map_recv_session_.erase(p.rdp_id_);
			delete rs;

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "recv session count:%d", rdp->map_recv_session_.size()));
		}
	}while(1);

	return 0;
}

/* @brief 处理有输出事件的回调函数 */
int myrdp::_rdp_handle_output(unsigned long context_data, int fd)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_rdp_handle_output context_data:%x fd:%d", context_data, fd));
	return 0;
}

/* @brief 处理有异常事件的回调函数*/
int myrdp::_rdp_handle_exception(unsigned long context_data, int fd)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_rdp_handle_exception context_data:%x fd:%d", context_data, fd));
	return 0;
}

/* @brief 处理消息队列里的消息 */
int myrdp::_rdp_handle_msg(unsigned long context_data, void * msg)
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_rdp_handle_msg context_data:%x msg:%x", context_data, msg));

	myrdp * rdp = (myrdp *)context_data;

	rdp_msg_head * m = (rdp_msg_head *)msg;

	switch(m->type_)
	{
	case RDP_MSG_SEND:
		{
			rdp_msg_send * ms = (rdp_msg_send *)m;
			rdp->_send_data(ms->ctx_, &ms->data_[0], ms->data_.size(), ms->dst_ip_, ms->dst_port_, ms->resend_cout_, ms->resend_interval_);
		}
		break;

	default:
		break;
	}

	delete m;

	return 0;
}

/* @brief 将定时器调至100ms间隔 */
void myrdp::_shift_timer_to_100()
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_shift_timer_to_100"));

	mytimer_node_t n = {0};
	n.context_data = (unsigned long)this;
	n.timer_user_data = (unsigned long)NULL;
	n.timeout_cb = myrdp::_rdp_timeout;
	n.period.tv_sec = 0;
	n.period.tv_usec = 100 * 1000;
	n.first_expire.tv_sec = 0;
	n.first_expire.tv_usec = 100 * 1000;

	this->time_resolution_ = 1;

	this->htimer_ = MyListernerResetTimer(this->lsn_, this->htimer_, &n);
}
/* @brief 将定时器调至1000ms间隔 */
void myrdp::_shift_timer_to_1000()
{
	MYLOG_DEBUGEX((MYRDP_LOG_PRE, "myrdp::_shift_timer_to_1000"));

	mytimer_node_t n = {0};
	n.context_data = (unsigned long)this;
	n.timer_user_data = (unsigned long)NULL;
	n.timeout_cb = myrdp::_rdp_timeout;
	n.period.tv_sec = 1;
	n.period.tv_usec = 0;
	n.first_expire.tv_sec = 1;
	n.first_expire.tv_usec = 0;

	this->time_resolution_ = 10;

	this->htimer_ = MyListernerResetTimer(this->lsn_, this->htimer_, &n);
}

/* @brief 超时回调 */
int myrdp::_rdp_timeout(unsigned long context_data, unsigned long timer_user_data,HTIMERID timerid)
{
	myrdp * rdp = (myrdp *)context_data;
	rdp->time_tick_ += rdp->time_resolution_;

	{
		/* 扫描所有的未结session,重发未收ack的报文分片 */
		for(std::map<uint64, send_session_t *>::iterator it = rdp->map_send_session_.begin();
			it != rdp->map_send_session_.end() && rdp->map_send_session_.size(); it ++)
		{
			assert(it->second);

			/* 是否超时重发 */
			if(it->second->last_send_tick_ + it->second->resend_interval_tick_ > rdp->time_tick_)
				continue;

			if(0 == it->second->resend_cout_)
			{
				rdp->send_data_res(it->second->ctx_, SEND_TIMER_OUT);

				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "resend count to max, fail ..."));
				std::map<uint64, send_session_t *>::iterator it_temp = it;
				it ++;

				/* 释放session */
				delete it_temp->second;
				rdp->map_send_session_.erase(it_temp);

				if(it == rdp->map_send_session_.end() || 0 == rdp->map_send_session_.size())
				{
					MYLOG_DEBUGEX((MYRDP_LOG_PRE, "have no send session any more, break ..."));
					break;
				}
			}

			/* 如果实际收包速度慢,采用实际速度,如果实际收包速度等于发送速度,尝试加快发送速度 */
			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "ack_count_:%d last_ack_count_:%d last_send_speed_:%d",
				it->second->ack_count_,
				it->second->last_ack_count_,
				it->second->last_ack_count_));

			uint32 recv_speed = it->second->ack_count_ - it->second->last_ack_count_;
			it->second->last_ack_count_ = it->second->ack_count_;
			if(recv_speed < it->second->last_send_speed_)
				it->second->last_ack_count_ = recv_speed;
			else
				it->second->last_ack_count_ = recv_speed + recv_speed/4;

			uint32 send_count_ = 0;
			for(uint32 i = 0; i < it->second->vpkt_frame_.size(); i ++)
			{
				if(!it->second->vpkt_frame_[i])
					continue;

				if(send_count_ <= it->second->last_ack_count_)
				{
					pkt_frame_t * pf = it->second->vpkt_frame_[i];

					MYLOG_DEBUGEX((MYRDP_LOG_PRE, "resend block:%d rdp id:%llx", i, it->first));

					CChannel::UdpWrite(rdp->udp_fd_, &pf->pkt[0], pf->pkt.size(), it->second->dst_ip_, it->second->dst_port_);

					send_count_ ++;
				}
				else
				{
					MYLOG_DEBUGEX((MYRDP_LOG_PRE, "send max:%d", send_count_));
					break;
				}
			}

			it->second->last_send_tick_ = rdp->time_tick_;

			it->second->resend_cout_ --;
		}
	}

	{
		/* 检查是否有接收session超时 */
		for(std::map<uint64, recv_session_t *>::iterator it = rdp->map_recv_session_.begin();
			it != rdp->map_recv_session_.end() && rdp->map_recv_session_.size(); it ++)
		{
			assert(it->second);

			if(it->second->recv_expire_tick_ > rdp->time_tick_)
				continue;

			MYLOG_DEBUGEX((MYRDP_LOG_PRE, "recv session expire %llx", it->first));

			std::map<uint64, recv_session_t *>::iterator it_temp = it;
			it ++;

			/* 释放session */
			delete it_temp->second;
			rdp->map_recv_session_.erase(it_temp);

			if(it == rdp->map_recv_session_.end() || 0 == rdp->map_recv_session_.size())
			{
				MYLOG_DEBUGEX((MYRDP_LOG_PRE, "have no recv session any more, break ..."));
				break;
			}
		}
	}

	return 0;
}







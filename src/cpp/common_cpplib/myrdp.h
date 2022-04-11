/**
 * @file myrdp.h 
 * @brief 可靠udp
 *
 * @author linshaochuan
 * @blog http://blog.csdn.net/lsccsl
 */
#ifndef __MYRDP_H__
#define __MYRDP_H__

#include <string.h>
#include <vector>
#include <map>
#include "type_def.h"
extern "C"
{
    #include "mylisterner.h"
}

class myrdp
{
public:

	enum{
		/* @brief 发送成功 */
		SEND_OK = 0,
		
		/* @brief 发送超时 */
		SEND_TIMER_OUT = -1,
	};

public:

	/**
	* @brief constructor
	*/
	myrdp(const int8 * local_ip, const uint16 port, const uint32 host_id,
		const uint32 mtu = 1024,
		const uint32 max_rcv_time_out = 30000);

	/**
	* @brief destructor
	*/
	virtual ~myrdp();

	/**
	* @brief init rdp
	*/
	int32 init();

	/**
	* @brief send data
	* @param resend_cout:重发次数
	* @param resend_interval:重发间隔(毫秒,千分之一秒)
	*/
	int32 send_data(const unsigned long ctx,
		std::vector<uint8>& data,
		const int8 * dst_ip, const uint16 dst_port,
		const uint32 resend_cout = 3, const uint32 resend_interval = 1000);

	/**
	* @brief 发送回调
	* @param result: 0:成功 其它:失败
	*/
	virtual void send_data_res(const unsigned long ctx, unsigned int result);

	/**
	* @brief 收到数据回调函(不可阻塞)
	*/
	virtual void recv_data(const uint8 * data, const uint32 data_sz, const int8 * src_ip, const uint16 src_port);

private:

	/* @brief 生成rdp id */
	void _gen_rdp_id(uint64& rdp_id);

	/* @brief 将定时器调至100ms间隔 */
	void _shift_timer_to_100();
	/* @brief 将定时器调至1000ms间隔 */
	void _shift_timer_to_1000();

	/**
	* @brief send data
	* @param resend_cout:重发次数
	* @param resend_interval:重发间隔(毫秒,千分之一秒)
	*/
	int32 _send_data(const unsigned long ctx,
		const uint8 * data, const uint32 data_sz,
		const int8 * dst_ip, const uint16 dst_port,
		const uint32 resend_cout, const uint32 resend_interval);

private:

	/* @brief 处理有输入事件的回调函数 */
	static int _rdp_handle_input(unsigned long context_data, int fd);

	/* @brief 处理有输出事件的回调函数 */
	static int _rdp_handle_output(unsigned long context_data, int fd);

	/* @brief 处理有异常事件的回调函数*/
	static int _rdp_handle_exception(unsigned long context_data, int fd);

	/* @brief 处理消息队列里的消息 */
	static int _rdp_handle_msg(unsigned long context_data, void * msg);

	/* @brief 超时回调 */
	static int _rdp_timeout(unsigned long context_data, unsigned long timer_user_data,HTIMERID timerid);

private:

	/* @brief 报头 */
	static const int32 _rdp_head_len_ = 32;

private:

	/* 报文切片缓存 */
	struct pkt_frame_t
	{
		pkt_frame_t()
		{}

		/* @brief 报文 */
		std::vector<uint8> pkt;
	};

	struct send_session_t
	{
		send_session_t():dst_port_(0),resend_cout_(3),resend_interval_tick_(1),ack_count_(0),last_ack_count_(0),last_ack_timetick_(0),last_send_speed_(10)
		{
			memset(dst_ip_, 0, sizeof(dst_ip_));
		}

		~send_session_t()
		{
			for(uint32 i = 0; i < this->vpkt_frame_.size(); i ++)
			{
				if(this->vpkt_frame_[i])
					delete this->vpkt_frame_[i];

				this->vpkt_frame_[i] = NULL;
			}
		}

		int8 dst_ip_[16];
		uint16 dst_port_;
		uint32 resend_cout_;
		uint32 resend_interval_tick_;
		/* @brief 上次发送时时 */
		uint32 last_send_tick_;

		/* @brief 收到ack的个数 */
		uint32 ack_count_;

		/* @brief 要发送的报文分片 */
		std::vector<pkt_frame_t *> vpkt_frame_;

		/* @brief 上下文数据 */
		unsigned long ctx_;

		/* @brief 上次收到ack的个数 */
		uint32 last_ack_count_;
		/* @brief 收ack的时间 */
		uint32 last_ack_timetick_;
		/* @brief 上一次发送了多少报文 */
		uint32 last_send_speed_;
	};
	/* @brief 报文会话 */
	std::map<uint64, send_session_t *> map_send_session_;

	struct recv_session_t
	{
		recv_session_t(const uint32 recv_expire_tick):total_block_(0),cur_rcv_block_(0),b_rcv_first_(0),recv_expire_tick_(recv_expire_tick)
		{}

		/* @brief 是否收到了第一块 */
		int32 b_rcv_first_;

		/* @brief 总共有多少块 */
		uint32 total_block_;
		/* @brief 当前收了多少块 */
		uint32 cur_rcv_block_;

		/* @brief 接收开始的tick */
		uint32 recv_expire_tick_;

		/* @brief 接收的报文分片 */
		std::vector<pkt_frame_t> vpkt_frame_;
	};
	std::map<uint64, recv_session_t *> map_recv_session_;

	/* @brief 监听器 */
	HMYLISTERNER lsn_;

	/* @brief udp fd */
	int32 udp_fd_;

	/* @brief host id */
	uint32 host_id_;

	/* @brief mtu */
	uint32 mtu_;

	/* @brief rdp seq,每发出一个报文递增,至0xffffff,重新回到1 */
	uint32 rdp_seq_;

	/* @brief 扫描定时器 */
	HTIMERID htimer_;
	/* @brief 时间解析度 单位100ms */
	uint32 time_resolution_;
	/* @brief tick,单位100ms */
	uint64 time_tick_;

	/* @brief 接收超时的tick 100ms */
	uint32 max_recv_timeout_tick_;
};

#endif







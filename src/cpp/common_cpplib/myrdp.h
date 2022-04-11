/**
 * @file myrdp.h 
 * @brief �ɿ�udp
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
		/* @brief ���ͳɹ� */
		SEND_OK = 0,
		
		/* @brief ���ͳ�ʱ */
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
	* @param resend_cout:�ط�����
	* @param resend_interval:�ط����(����,ǧ��֮һ��)
	*/
	int32 send_data(const unsigned long ctx,
		std::vector<uint8>& data,
		const int8 * dst_ip, const uint16 dst_port,
		const uint32 resend_cout = 3, const uint32 resend_interval = 1000);

	/**
	* @brief ���ͻص�
	* @param result: 0:�ɹ� ����:ʧ��
	*/
	virtual void send_data_res(const unsigned long ctx, unsigned int result);

	/**
	* @brief �յ����ݻص���(��������)
	*/
	virtual void recv_data(const uint8 * data, const uint32 data_sz, const int8 * src_ip, const uint16 src_port);

private:

	/* @brief ����rdp id */
	void _gen_rdp_id(uint64& rdp_id);

	/* @brief ����ʱ������100ms��� */
	void _shift_timer_to_100();
	/* @brief ����ʱ������1000ms��� */
	void _shift_timer_to_1000();

	/**
	* @brief send data
	* @param resend_cout:�ط�����
	* @param resend_interval:�ط����(����,ǧ��֮һ��)
	*/
	int32 _send_data(const unsigned long ctx,
		const uint8 * data, const uint32 data_sz,
		const int8 * dst_ip, const uint16 dst_port,
		const uint32 resend_cout, const uint32 resend_interval);

private:

	/* @brief �����������¼��Ļص����� */
	static int _rdp_handle_input(unsigned long context_data, int fd);

	/* @brief ����������¼��Ļص����� */
	static int _rdp_handle_output(unsigned long context_data, int fd);

	/* @brief �������쳣�¼��Ļص�����*/
	static int _rdp_handle_exception(unsigned long context_data, int fd);

	/* @brief ������Ϣ���������Ϣ */
	static int _rdp_handle_msg(unsigned long context_data, void * msg);

	/* @brief ��ʱ�ص� */
	static int _rdp_timeout(unsigned long context_data, unsigned long timer_user_data,HTIMERID timerid);

private:

	/* @brief ��ͷ */
	static const int32 _rdp_head_len_ = 32;

private:

	/* ������Ƭ���� */
	struct pkt_frame_t
	{
		pkt_frame_t()
		{}

		/* @brief ���� */
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
		/* @brief �ϴη���ʱʱ */
		uint32 last_send_tick_;

		/* @brief �յ�ack�ĸ��� */
		uint32 ack_count_;

		/* @brief Ҫ���͵ı��ķ�Ƭ */
		std::vector<pkt_frame_t *> vpkt_frame_;

		/* @brief ���������� */
		unsigned long ctx_;

		/* @brief �ϴ��յ�ack�ĸ��� */
		uint32 last_ack_count_;
		/* @brief ��ack��ʱ�� */
		uint32 last_ack_timetick_;
		/* @brief ��һ�η����˶��ٱ��� */
		uint32 last_send_speed_;
	};
	/* @brief ���ĻỰ */
	std::map<uint64, send_session_t *> map_send_session_;

	struct recv_session_t
	{
		recv_session_t(const uint32 recv_expire_tick):total_block_(0),cur_rcv_block_(0),b_rcv_first_(0),recv_expire_tick_(recv_expire_tick)
		{}

		/* @brief �Ƿ��յ��˵�һ�� */
		int32 b_rcv_first_;

		/* @brief �ܹ��ж��ٿ� */
		uint32 total_block_;
		/* @brief ��ǰ���˶��ٿ� */
		uint32 cur_rcv_block_;

		/* @brief ���տ�ʼ��tick */
		uint32 recv_expire_tick_;

		/* @brief ���յı��ķ�Ƭ */
		std::vector<pkt_frame_t> vpkt_frame_;
	};
	std::map<uint64, recv_session_t *> map_recv_session_;

	/* @brief ������ */
	HMYLISTERNER lsn_;

	/* @brief udp fd */
	int32 udp_fd_;

	/* @brief host id */
	uint32 host_id_;

	/* @brief mtu */
	uint32 mtu_;

	/* @brief rdp seq,ÿ����һ�����ĵ���,��0xffffff,���»ص�1 */
	uint32 rdp_seq_;

	/* @brief ɨ�趨ʱ�� */
	HTIMERID htimer_;
	/* @brief ʱ������� ��λ100ms */
	uint32 time_resolution_;
	/* @brief tick,��λ100ms */
	uint64 time_tick_;

	/* @brief ���ճ�ʱ��tick 100ms */
	uint32 max_recv_timeout_tick_;
};

#endif







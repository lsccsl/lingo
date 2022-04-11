/**
 * @file mythrdpoll.h 
 * @brief ����sock��Ϣ�̳߳�
 *
 * @author linshaochuan
 * @blog http://blog.csdn.net/lsccsl
 */
#ifndef __MYTHRD_POLL_H__
#define __MYTHRD_POLL_H__

#pragma   warning(   disable   :   4786)

#include <map>
#include <string>
#include <pthread.h>
#include <vector>
#include <list>
extern "C"
{
	#include "mylisterner.h"
}
#include "type_def.h"
class CMyEPoll;


class ThrdHandle
{
public:

	virtual ~ThrdHandle(){}

	/**
	 * @brief ���ݻص�
	 */
	virtual int32 data_tcp_in(int32 fd, const void * data, const uint32 data_sz, uint32& recved) = 0;

	/**
	 * @brief �ر�fd,�����û���˵,�ر�tcp fd������ζ��һ���Ự������
	 */
	virtual int32 close_fd(int32 fd) = 0;

	/**
	 * @brief ���µ�tcp����
	 */
	virtual int32 tcp_conn(int32 fd, int32 fd_master) = 0;
	virtual int32 accept_have_conn(int32 fd_master, int32 fd_new_conn){ return 0; }

	/**
	 * @brief udp���ݻص�
	 */
	virtual int32 data_udp_in(int32 fd, const void * data, const uint32 data_sz,
		int8 * src_ip, uint16 src_port) = 0;

	/**
	 * @brief ��Ϣ�ص�
	 * @param thrd_from:��Ϣ�����ĸ��߳�,���Ϊ-1,��ʾ������Ϣ�̲߳����ڱ��̳߳�
	 */
	virtual int32 msg_callback(int32 thrd_from, const void * msg) = 0;

	/**
	 * @brief ��ʱ�ص�
	 */
	virtual int32 time_out(uint32 timer_data, HTIMERID timerid) = 0;
};

/**
 * @brief ����sock��Ϣ�̳߳�
 */
class CMyThrdPoll
{
public:

	/**
	 * @brief constructor
	 */
	CMyThrdPoll(const uint32 thread_count = 10, const uint32 max_fd_count = 1024,
		const uint32 max_msg_count = 65535, const uint32 bufsz_reserve = 512,
		const uint32 epoll_thrd_count = 1);

	/**
	 * @brief destructor
	 */
	virtual ~CMyThrdPoll();

	/**
	 * @brief ��ʼ��
	 */
	int32 init();
	/**
	 * @brief ���û���ʼ���߳�����
	 */
	virtual int32 InitThrdHandle(ThrdHandle *& thrd_handle, uint32 thrd_index) = 0;

	/**
	 * @brief ������ѭ��
	 */
	int32 work_loop(int32 timeout = -1);

	/**
	 * @brief ��Ҫ���͵�����ѹ��������,������epoll��wait״̬,���ں��ʵ�ʱ���ٷ�������,�ɶ��߳�ͬʱ���ʴ˺���
	 * @param fd:���
	 * @param data:���ݻ�����
	 * @param data_sz:data�Ĵ�С
	 */
	int32 data_tcp_out(int32 fd, std::vector<uint8>& data);
	/**
	 * @brief ֱ�ӽ����ݷ��ͳ�ȥ,������,�û���֤��ͬһ��fd�Ķ�д�Ǵ��е�(���߼���,���߸�fd��Զֻ��һ���߳�д)
	 */
	int32 data_tcp_out_sync(int32 fd, std::vector<uint8>& data);

	/**
	 * @brief ��ָ�����̷߳���Ϣ
	 */
	int32 push_msg(void * msg, uint32 thrd_to, int32 thrd_from = -1);
	/**
	* @brief ����fd������Ϣ
	*/
	int32 push_msg_by_fd(void * msg, int32 fd_to, int32 thrd_from = -1);

	/**
	 * @brief ��Ӷ�ʱ��
	 */
	int32 add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period = 0);
	/**
	 * @brief ɾ����ʱ��
	 */
	int32 del_time(uint32 thrd_index, HTIMERID timer_id);

	/**
	 * @brief ��һ��udp fd�������
	 */
	int32 add_udp_fd(int32 udp_fd, uint32& thrd_index);

	/**
	 * @brief ��һ��tcp fd�������
	 */
	int32 add_tcp_srv_fd(int32 tcp_fd, uint32& thrd_index, int32 bauto_add = 1);

	/**
	 * @brief ��һ��tcp client fd �������
	 */
	int32 add_tcp_cli_fd(int32 tcp_fd, uint32& thrd_index);

	/**
	 * @brief ȡ����һ������ļ���
	 */
	int32 del_and_close_fd(int32 fd);

	/**
	* @brief close all open tcp connect(listerner tcp sock and udp don't close)
	*/
	void close_tcp_connect();

	/**
	 * @brief view
	 */
	void runtime_view();

	enum
	{
		E_ERR = 0,

		/* epoll��ѭ���߳� -> ������߳� ��������������¼� */
		E_INPUT_TCP,

		/* epoll��ѭ���߳� -> ������߳� ��������������¼� */
		E_OUTPUT_TCP,

		/* epoll��ѭ���߳� -> ������߳� tcp��������������¼� */
		E_ACCEPT_TCP,

		/* �����߳� -> ������߳� �û������¼� */
		E_USER_OUTPUT_TCP,

		/* ������߳� -> ������߳� �����Ӳ����Ļص��¼� */
		E_NEW_CONN,

		/* epoll��ѭ���߳� -> ������߳� udp�����¼� */
		E_INPUT_UDP,

		/* ɾ����� */
		E_DEL_AND_CLOSE_FD,
	};

protected:

	enum SOCKET_FD_TYPE_T
	{
		/* ����accept״̬��tcp fd */
		TCP_ACCEPTOR_FD,
		/* ���յ��ͻ������Ӷ����ɵ�tcp fd */
		TCP_CONNECTOR_FD,

		/* ���ӵ�tcp��������ɵ�fd */
		TCP_CONNECTOR_CLI_FD,

		/* udp fd */
		UDP_FD,
	};

private:

	struct fd_info_t
	{
		fd_info_t():byte_write_(0),byte_read_(0){}

		/**
		 * @brief fd
		 */
		int32 fd_;
		/**
		 * @brief �������
		 */
		SOCKET_FD_TYPE_T type_;
		/**
		 * @brief ��������sock��������Ϣ,����֪��ĳ��sock��Ӧ�ĸ��ն�ip,��������־ʱʹ��
		 */
		std::string remote_ip_;/* ֻ��accept��connect������fd������ */
		uint32 remote_port_;/* ֻ��accept��connect������fd������ */
		std::string local_ip_;
		uint32 local_port_;
		int32 fd_master_;/* TCP_CONNECTOR_FD ֻ��accept������sock������,��¼fd_�����ĸ�����tcp���������� */
		int32 bauto_add_to_listern_;/* TCP_ACCEPTOR_FD���͵�fd�Ƿ�accept������fd�Զ�������� 0:������ 1:���� */

		/**
		 * @brief fd�����Ĵ����߳�
		 */
		HMYLISTERNER hlsn_;
		/**
		 * @brief fd�����̵߳Ĵ���������
		 */
		ThrdHandle * thrd_context_data_;

		/**
		 * @brief �������ݻ���(tcp udp���ݶ���������,��ֻ����tcp,���ܻ���udp)
		 */
		std::vector<uint8> rbuf_;
		/**
		 * @brief recv pos,��ǰ�����ݵ���ʼλ��(��udp��Ч,���ܻ���udp����,������ֱ��֪ͨ�û�)
		 */
		uint32 rpos_;
		/**
		 * @brief read pos,��Ӧ�ö��ߵ�����δλ��(��udp��Ч,���ܻ���udp����,������ֱ��֪ͨ�û�,�������û��Ƿ����˸�����,���´ν��ջᱻ����)
		 */
		uint32 rrpos_;

		/**
		 * @brief Ҫ���͵����ݻ���(��udp��Ч(����û��udp���ͽӿ�),udp����ֱ�ӷ���,����Ҫ����), С����ֱ�������������������
		 */
		std::vector<uint8> sbuf_;
		/* @brief �󱸷��ͻ���,��Դ��ĵ��Ż� */
		std::list<std::vector<uint8> > lst_sbuf_;

		/* @brief ͳ��,д�˶����ֽ� */
		uint32 byte_write_;
		/* @brief ͳ��,���˶����ֽ� */
		uint32 byte_read_;
	};

	struct thrd_msg_user_output
	{
		std::vector<uint8> msg;
	};

	struct thrd_msg_t
	{
		thrd_msg_t():msg_body_(NULL){}

		void free_body(HMYMEMPOOL hm)
		{
			if(msg_body_)
			{
				switch(this->event_)
				{
				case E_USER_OUTPUT_TCP:
					((thrd_msg_user_output*)msg_body_)->~thrd_msg_user_output();
					MyMemPoolFree(hm, msg_body_);
					break;
				}
			}
		}

		/**
		 * @brief �����¼��ľ��
		 */
		int32 fd_;

		/**
		 * @brief �¼�����
		 */
		uint32 event_;

		/**
		 * @brief ��Ϣ��
		 */
		void * msg_body_;
	};


	struct thrd_user_msg_t_
	{
		/* ���ĸ��̷߳��� */
		uint32 thrd_from_;
		/* �����ĸ��߳� */
		uint32 thrd_to_;

		void * msg_body_;
	};

	/**
	 * @brief accept new connection, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoAcceptTcp(int32 fd);
	/**
	 * @brief do input, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoInputTcp(int32 fd);
	/**
	 * @brief do output, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoOutputTcp(int32 fd);
	/**
	 * @brief do user output, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoUserOutputTcp(int32 fd, CMyThrdPoll::thrd_msg_t * msg);
	/**
	 * @brief do new connection in event, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoNewConn(int32 fd);
	/**
	 * @brief do udp input
	 */
	int32 _thrdsDoInputUdp(int32 fd);
	/**
	 * @brief do err, ���̳߳���ĳ���̵߳��������б�����
	 */
	int32 _thrdsDoErr(int32 fd);
	/**
	 * @brief ɾ������
	 */
	int32 _thrdsDoDelAndClose(int32 fd);
	/**
	 * @brief ������Ϣ���������Ϣ, ���̳߳���ĳ���̵߳��������б�����
	 */
	static int32 _thrdsMsgCb(unsigned long context_data, void * msg);
	/**
	 * @brief �û���Ϣ�ص�����
	 */
	static int32 _thrdsUserMsgCb(unsigned long context_data, void * msg);
	/**
	 * @brief ��ʱ�ص�
	 */
	static int32 _thrdsTimeOut(unsigned long context_data, unsigned long timer_user_data, HTIMERID timerid);


	/**
	 * @brief process tcp input,��work_loop���ڵ��߳��б�����
	 */
	static int32 _EpollInputTcp(void * context_data, int32 fd);
	/**
	 * @brief process accept event,��work_loop���ڵ��߳��б�����
	 */
	static int32 _EpollInputAccept(void * context_data, int32 fd);
	/**
	 * @brief process output ��work_loop���ڵ��߳��б�����
	 */
	static int32 _EpollOutputTcp(void * context_data, int32 fd);
	/**
	 * @brief process udp input,��work_loop���ڵ��߳��б�����
	 */
	static int32 _EpollInputUdp(void * context_data, int32 fd);
	/**
	 * @brief process err ��work_loop���ڵ��߳��б�����
	 */
	static int32 _EpollErr(void * context_data, int32 fd);


	/**
	 * @brief ����� map_thrd_ need lock
	 */
	int32 _addToMapThrd(int32 fd, SOCKET_FD_TYPE_T fd_type, uint32& thrd_index, int32 fd_master = -1);
	/**
	 * @brief ɾ�� need lock
	 */
	int32 _delFromMapThrd(int32 fd);
	/**
	 * @brief ��ȡfd������������ need lock
	 */
	int32 _getFdInfoFromMapThrd(int32 fd, fd_info_t*& i);
	int32 _getFdThrd(int32 fd, HMYLISTERNER& hlsn_);

	/**
	 * @brief ���ó��Զ��������Զ��������
	 */
	int32 _set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add);

	/**
	 * @brief �����̳߳��ڲ��߳�֮�����Ϣ
	 */
	int32 _push_inter_msg(HMYLISTERNER hlsn, int32 fd, uint32 ev, void * pb);

private:

	/**
	 * @brief fd - thrd map
	 */
	std::map<int32, fd_info_t *> map_thrd_;
	pthread_mutex_t map_thrd_protector_;

	/**
	 * @brief thrd
	 */
	struct thrd_info_t
	{
		HMYLISTERNER thrd_;
		ThrdHandle * thrd_context_data_;
	};
	std::vector<thrd_info_t> thrds_;

	/**
	 * @brief �̸߳���
	 */
	uint32 real_thrd_count_;
	/**
	 * @brief �����Ϣ����
	 */
	uint32 max_msg_count_;

	/**
	 * @brief epoll api wrapper
	 */
	CMyEPoll * epoll_;

	/*! ���ͻ�����������С,С�ڴ�ֵ,�򽫻�������������������ǰ��,���⻺������������ */
	uint32 bufsz_reserve_;

	/**
	 * @brief �ڴ�ؾ��
	 */
	HMYMEMPOOL hm_;
};

#endif



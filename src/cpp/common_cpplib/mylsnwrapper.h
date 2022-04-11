/**
* @file mylsnwrapper.h
* @brief multi thread select listerner wrapper,
*        wrapper tcp detail, provide network data to up layer
* @author linsc
* @blog http://blog.csdn.net/lsccsl
*/
#ifndef __MYLSNWRAPPER_H__
#define __MYLSNWRAPPER_H__

#include <vector>
#include <map>
#include <string>
#include <pthread.h>
#include "type_def.h"
extern "C"
{
    #include "mylisterner.h"
}

class LsnThrdHandle
{
public:

	virtual ~LsnThrdHandle(){}

	/**
	* @brief tcp data in callback
	*/
	virtual int32 data_tcp_in(int32 fd, const void * data, const uint32 data_sz, uint32& recved) = 0;

	/**
	* @brief fd close callback
	*/
	virtual int32 close_fd(int32 fd) = 0;

	/**
	* @brief have new tcp connect
	*/
	virtual int32 tcp_conn(int32 fd, int32 fd_master) = 0;
	virtual int32 accept_have_conn(int32 fd_master, int32 fd_new_conn){ return 0; }

	/**
	* @brief have udp data in
	*/
	virtual int32 data_udp_in(int32 fd, const void * data, const uint32 data_sz,
		int8 * src_ip, uint16 src_port) = 0;

	/**
	* @brief msg callback
	* @param thrd_from:which thread msg from, if is -1, means the msg come from the thread outof this module
	*/
	virtual int32 msg_callback(int32 thrd_from, const void * msg) = 0;

	/**
	* @brief timer callback
	*/
	virtual int32 time_out(uint32 timer_data, HTIMERID timerid) = 0;
};


class mylsnwrapper
{
public:

	/**
	* @brief constructor
	*/
	mylsnwrapper(uint32 thread_count = 1, uint32 max_fd_count = 1024,
		uint32 max_msg_count = 65535, uint32 bufsz_reserve = 512);

	/**
	* @brief destructor
	*/
	virtual ~mylsnwrapper();

	/**
	* @brief ��ʼ��
	*/
	int32 init();
	/**
	* @brief ���û���ʼ���߳�����
	*/
	virtual int32 InitThrdHandle(LsnThrdHandle *& thrd_handle, uint32 thrd_index) = 0;

	/**
	* @brief add tcp srv to listern
	*/
	int32 add_tcp_srv_fd(int32 tcp_fd, uint32& thrd_index, int32 bauto_add = 1);

	/**
	* @brief add tcp cli/conn
	*/
	int32 add_tcp_cli_fd(int32 tcp_fd, uint32& thrd_index);

	/**
	* @brief add udp fd to listener
	*/
	int32 add_udp_fd(int32 udp_fd, uint32& thrd_index);

	/**
	* @brief close and remove from listener
	*/
	int32 del_and_close_fd(int32 fd);

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
	static int32 data_tcp_out_sync(int32 fd, std::vector<uint8>& data);

	/**
	* @brief ��ָ�����̷߳���Ϣ
	*/
	int32 push_msg(void * msg, uint32 thrd_to, int32 thrd_from = -1);

	/**
	 * @brief ��Ӷ�ʱ��
	 */
	int32 add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period = 0);
	/**
	 * @brief ɾ����ʱ��
	 */
	int32 del_time(uint32 thrd_index, HTIMERID timer_id);

private:

	enum
	{
		/* �����߳� -> ������߳� �û������¼� */
		E_USER_OUTPUT_TCP,

		/* ������߳� -> ������߳� �����Ӳ����Ļص��¼� */
		E_NEW_CONN,

		/* ɾ����� */
		E_DEL_AND_CLOSE_FD,
	};

	enum
	{
		/* acceptor tcp fd */
		TCP_ACCEPTOR_FD,
		/* fd produce by accept */
		TCP_CONNECTOR_FD,

		/* fd produce by connect */
		TCP_CONNECTOR_CLI_FD,

		/* udp fd */
		UDP_FD,
	};

	struct fd_info_t
	{
		/**
		* @brief fd
		*/
		int32 fd_;
		/**
		* @brief sock type, tcp srv, tcp connect(srv / client), or udp
		*/
		uint32 type_;
		/**
		* @brief another sock info
		*/
		std::string remote_ip_;/* meaning for tcp connection sock */
		uint32 remote_port_;/* meaning for tcp connection sock */
		std::string local_ip_;
		uint32 local_port_;
		int32 fd_master_;/* TCP_CONNECTOR_FD, meaning for sock by accept */
		int32 bauto_add_to_listern_;/* if type is TCP_ACCEPTOR_FD, be auto add to listerner, 0:no 1:yes */

		/**
		* @brief which thread does fd belong to
		*/
		HMYLISTERNER hlsn_;
		/**
		* @brief which thread context does fd belong to
		*/
		LsnThrdHandle * thrd_context_data_;

		/**
		* @brief tcp recv data buf
		*/
		std::vector<uint8> rbuf_;
		/**
		* @brief recv pos
		*/
		uint32 rcvpos_;
		/**
		* @brief read pos
		*/
		uint32 readpos_;

		/**
		* @brief the send data buf
		*/
		std::vector<uint8> sbuf_;
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

private:

	/**
	* @brief handle network input event callback
	*/
	static int __handle_tcp_input(unsigned long context_data, int fd);
	static int __handle_tcpsrv_input(unsigned long context_data, int fd);
	/** 
	* @brief handle out event callback
	*/
	static int __handle_tcp_output(unsigned long context_data, int fd);
	/**
	* @brief handle exception callback
	*/
	static int __handle_exception(unsigned long context_data, int fd);
	/**
	* @brief handle timer callback
	*/
	static int __handle_timeout(unsigned long context_data,  unsigned long timer_user_data, HTIMERID timerid);

	/**
	* @brief ������Ϣ���������Ϣ, ���̳߳���ĳ���̵߳��������б�����
	*/
	static int32 _thrdsMsgCb(unsigned long context_data, void * msg);
	/**
	* @brief �û���Ϣ�ص�����
	*/
	static int32 _thrdsUserMsgCb(unsigned long context_data, void * msg);

	/**
	* @brief do new connection in event, ���̳߳���ĳ���̵߳��������б�����
	*/
	int32 _thrdsDoNewConn(int32 fd);
	/**
	* @brief do output, ���̳߳���ĳ���̵߳��������б�����
	*/
	int32 _thrdsDoUserOutputTcp(int32 fd, mylsnwrapper::thrd_msg_t * msg);

	/**
	* @brief �������tcp���ݷ���ȥ
	*/
	int32 _inter_tcp_output(int32 fd);

private:

	/**
	* @brief ����� map_thrd_ need lock
	*/
	int32 _addToMapThrd(int32 fd, uint32 fd_type, uint32& thrd_index, uint32 mask, int32 fd_master = -1);
	/**
	* @brief ��ȡfd������������ need lock
	*/
	int32 _getFdInfoFromMapThrd(int32 fd, fd_info_t*& i);
	int32 _getFdThrd(int32 fd, HMYLISTERNER& hlsn);
	/**
	* @brief ɾ�� need lock
	*/
	int32 _delFromMapThrd(int32 fd);

	/**
	* @brief set auto listen or not
	*/
	int32 _set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add);

	/**
	* @brief �����̳߳��ڲ��߳�֮�����Ϣ
	*/
	int32 _push_inter_msg(HMYLISTERNER hlsn, int32 fd, uint32 ev, void * pb);

private:

	/**
	* @brief thread and thread context
	*/
	struct thrd_cxt_t
	{
		HMYLISTERNER thrd_;
		LsnThrdHandle * thrd_cxt_;
	};
	std::vector<thrd_cxt_t> vthrds_;

	/*! ���ͻ�����������С,С�ڴ�ֵ,�򽫻�������������������ǰ��,���⻺������������ */
	uint32 bufsz_reserve_;

	/**
	* @brief fd - thrd map
	*/
	std::map<int32, fd_info_t *> map_thrd_;
	pthread_mutex_t map_thrd_protector_;

	/**
	* @brief �ڴ�ؾ��
	*/
	HMYMEMPOOL hm_;

	/**
	* @brief �̸߳���
	*/
	uint32 real_thrd_count_;
	/**
	* @brief �����Ϣ����
	*/
	uint32 max_msg_count_;
};

#endif









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
	* @brief 初始化
	*/
	int32 init();
	/**
	* @brief 让用户初始化线程数据
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
	* @brief 将要发送的数据压进缓冲区,并更改epoll的wait状态,等于合适的时候再发送数据,可多线程同时访问此函数
	* @param fd:句柄
	* @param data:数据缓冲区
	* @param data_sz:data的大小
	*/
	int32 data_tcp_out(int32 fd, std::vector<uint8>& data);
	/**
	* @brief 直接将数据发送出去,不缓存,用户保证对同一个fd的读写是串行的(或者加锁,或者该fd永远只被一个线程写)
	*/
	static int32 data_tcp_out_sync(int32 fd, std::vector<uint8>& data);

	/**
	* @brief 给指定的线程发消息
	*/
	int32 push_msg(void * msg, uint32 thrd_to, int32 thrd_from = -1);

	/**
	 * @brief 添加定时器
	 */
	int32 add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period = 0);
	/**
	 * @brief 删除定时器
	 */
	int32 del_time(uint32 thrd_index, HTIMERID timer_id);

private:

	enum
	{
		/* 任意线程 -> 句柄绑定线程 用户输入事件 */
		E_USER_OUTPUT_TCP,

		/* 句柄绑定线程 -> 句柄绑定线程 新连接产生的回调事件 */
		E_NEW_CONN,

		/* 删除句柄 */
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
		* @brief 产生事件的句柄
		*/
		int32 fd_;

		/**
		* @brief 事件参数
		*/
		uint32 event_;

		/**
		* @brief 消息体
		*/
		void * msg_body_;
	};


	struct thrd_user_msg_t_
	{
		/* 从哪个线程发出 */
		uint32 thrd_from_;
		/* 发往哪个线程 */
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
	* @brief 处理消息队列里的消息, 在线程池中某个线程的上下文中被呼叫
	*/
	static int32 _thrdsMsgCb(unsigned long context_data, void * msg);
	/**
	* @brief 用户消息回调函数
	*/
	static int32 _thrdsUserMsgCb(unsigned long context_data, void * msg);

	/**
	* @brief do new connection in event, 在线程池中某个线程的上下文中被呼叫
	*/
	int32 _thrdsDoNewConn(int32 fd);
	/**
	* @brief do output, 在线程池中某个线程的上下文中被呼叫
	*/
	int32 _thrdsDoUserOutputTcp(int32 fd, mylsnwrapper::thrd_msg_t * msg);

	/**
	* @brief 将缓存的tcp内容发出去
	*/
	int32 _inter_tcp_output(int32 fd);

private:

	/**
	* @brief 添加至 map_thrd_ need lock
	*/
	int32 _addToMapThrd(int32 fd, uint32 fd_type, uint32& thrd_index, uint32 mask, int32 fd_master = -1);
	/**
	* @brief 获取fd的上下文数据 need lock
	*/
	int32 _getFdInfoFromMapThrd(int32 fd, fd_info_t*& i);
	int32 _getFdThrd(int32 fd, HMYLISTERNER& hlsn);
	/**
	* @brief 删除 need lock
	*/
	int32 _delFromMapThrd(int32 fd);

	/**
	* @brief set auto listen or not
	*/
	int32 _set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add);

	/**
	* @brief 推入线程池内部线程之间的消息
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

	/*! 发送缓冲区保留大小,小于此值,则将缓冲区拷贝到缓冲区最前面,避免缓冲区恶性增长 */
	uint32 bufsz_reserve_;

	/**
	* @brief fd - thrd map
	*/
	std::map<int32, fd_info_t *> map_thrd_;
	pthread_mutex_t map_thrd_protector_;

	/**
	* @brief 内存池句柄
	*/
	HMYMEMPOOL hm_;

	/**
	* @brief 线程个数
	*/
	uint32 real_thrd_count_;
	/**
	* @brief 最大消息数量
	*/
	uint32 max_msg_count_;
};

#endif









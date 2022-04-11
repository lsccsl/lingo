/**
 * @file mythrdpoll.h 
 * @brief 处理sock消息线程池
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
	 * @brief 数据回调
	 */
	virtual int32 data_tcp_in(int32 fd, const void * data, const uint32 data_sz, uint32& recved) = 0;

	/**
	 * @brief 关闭fd,对于用户来说,关闭tcp fd可能意味着一个会话结束了
	 */
	virtual int32 close_fd(int32 fd) = 0;

	/**
	 * @brief 有新的tcp连接
	 */
	virtual int32 tcp_conn(int32 fd, int32 fd_master) = 0;
	virtual int32 accept_have_conn(int32 fd_master, int32 fd_new_conn){ return 0; }

	/**
	 * @brief udp数据回调
	 */
	virtual int32 data_udp_in(int32 fd, const void * data, const uint32 data_sz,
		int8 * src_ip, uint16 src_port) = 0;

	/**
	 * @brief 消息回调
	 * @param thrd_from:消息来自哪个线程,如果为-1,表示发出消息线程不属于本线程池
	 */
	virtual int32 msg_callback(int32 thrd_from, const void * msg) = 0;

	/**
	 * @brief 超时回调
	 */
	virtual int32 time_out(uint32 timer_data, HTIMERID timerid) = 0;
};

/**
 * @brief 处理sock消息线程池
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
	 * @brief 初始化
	 */
	int32 init();
	/**
	 * @brief 让用户初始化线程数据
	 */
	virtual int32 InitThrdHandle(ThrdHandle *& thrd_handle, uint32 thrd_index) = 0;

	/**
	 * @brief 主工作循环
	 */
	int32 work_loop(int32 timeout = -1);

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
	int32 data_tcp_out_sync(int32 fd, std::vector<uint8>& data);

	/**
	 * @brief 给指定的线程发消息
	 */
	int32 push_msg(void * msg, uint32 thrd_to, int32 thrd_from = -1);
	/**
	* @brief 根据fd来推消息
	*/
	int32 push_msg_by_fd(void * msg, int32 fd_to, int32 thrd_from = -1);

	/**
	 * @brief 添加定时器
	 */
	int32 add_time(uint32 thrd_index, uint32 time_second, uint32 timer_data, HTIMERID& timer_id, int32 period = 0);
	/**
	 * @brief 删除定时器
	 */
	int32 del_time(uint32 thrd_index, HTIMERID timer_id);

	/**
	 * @brief 将一个udp fd加入监听
	 */
	int32 add_udp_fd(int32 udp_fd, uint32& thrd_index);

	/**
	 * @brief 将一个tcp fd加入监听
	 */
	int32 add_tcp_srv_fd(int32 tcp_fd, uint32& thrd_index, int32 bauto_add = 1);

	/**
	 * @brief 将一个tcp client fd 加入监听
	 */
	int32 add_tcp_cli_fd(int32 tcp_fd, uint32& thrd_index);

	/**
	 * @brief 取消对一个句柄的监听
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

		/* epoll主循环线程 -> 句柄绑定线程 句柄产生了输入事件 */
		E_INPUT_TCP,

		/* epoll主循环线程 -> 句柄绑定线程 句柄产生了输入事件 */
		E_OUTPUT_TCP,

		/* epoll主循环线程 -> 句柄绑定线程 tcp监听句柄新连接事件 */
		E_ACCEPT_TCP,

		/* 任意线程 -> 句柄绑定线程 用户输入事件 */
		E_USER_OUTPUT_TCP,

		/* 句柄绑定线程 -> 句柄绑定线程 新连接产生的回调事件 */
		E_NEW_CONN,

		/* epoll主循环线程 -> 句柄绑定线程 udp输入事件 */
		E_INPUT_UDP,

		/* 删除句柄 */
		E_DEL_AND_CLOSE_FD,
	};

protected:

	enum SOCKET_FD_TYPE_T
	{
		/* 处于accept状态的tcp fd */
		TCP_ACCEPTOR_FD,
		/* 接收到客户端连接而生成的tcp fd */
		TCP_CONNECTOR_FD,

		/* 连接到tcp服务端生成的fd */
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
		 * @brief 句柄类型
		 */
		SOCKET_FD_TYPE_T type_;
		/**
		 * @brief 其他关于sock的杂项信息,可以知道某个sock对应哪个终端ip,供分析日志时使用
		 */
		std::string remote_ip_;/* 只对accept与connect产生的fd有意义 */
		uint32 remote_port_;/* 只对accept与connect产生的fd有意义 */
		std::string local_ip_;
		uint32 local_port_;
		int32 fd_master_;/* TCP_CONNECTOR_FD 只对accept产生的sock有意义,记录fd_是由哪个本地tcp监听产生的 */
		int32 bauto_add_to_listern_;/* TCP_ACCEPTOR_FD类型的fd是否将accept产生的fd自动加入监听 0:不加入 1:加入 */

		/**
		 * @brief fd所属的处理线程
		 */
		HMYLISTERNER hlsn_;
		/**
		 * @brief fd所属线程的处理上下文
		 */
		ThrdHandle * thrd_context_data_;

		/**
		 * @brief 接收数据缓冲(tcp udp数据都存在这里,但只缓存tcp,不能缓存udp)
		 */
		std::vector<uint8> rbuf_;
		/**
		 * @brief recv pos,当前收数据的起始位置(对udp无效,不能缓存udp报文,有数据直接通知用户)
		 */
		uint32 rpos_;
		/**
		 * @brief read pos,被应用读走的数据未位置(对udp无效,不能缓存udp报文,有数据直接通知用户,而不管用户是否处理了该数据,在下次接收会被覆盖)
		 */
		uint32 rrpos_;

		/**
		 * @brief 要发送的数据缓冲(对udp无效(本类没有udp发送接口),udp可以直接发送,不需要缓冲), 小报文直接在这个缓冲区里体现
		 */
		std::vector<uint8> sbuf_;
		/* @brief 后备发送缓冲,针对大报文的优化 */
		std::list<std::vector<uint8> > lst_sbuf_;

		/* @brief 统计,写了多少字节 */
		uint32 byte_write_;
		/* @brief 统计,读了多少字节 */
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

	/**
	 * @brief accept new connection, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoAcceptTcp(int32 fd);
	/**
	 * @brief do input, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoInputTcp(int32 fd);
	/**
	 * @brief do output, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoOutputTcp(int32 fd);
	/**
	 * @brief do user output, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoUserOutputTcp(int32 fd, CMyThrdPoll::thrd_msg_t * msg);
	/**
	 * @brief do new connection in event, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoNewConn(int32 fd);
	/**
	 * @brief do udp input
	 */
	int32 _thrdsDoInputUdp(int32 fd);
	/**
	 * @brief do err, 在线程池中某个线程的上下文中被呼叫
	 */
	int32 _thrdsDoErr(int32 fd);
	/**
	 * @brief 删除链接
	 */
	int32 _thrdsDoDelAndClose(int32 fd);
	/**
	 * @brief 处理消息队列里的消息, 在线程池中某个线程的上下文中被呼叫
	 */
	static int32 _thrdsMsgCb(unsigned long context_data, void * msg);
	/**
	 * @brief 用户消息回调函数
	 */
	static int32 _thrdsUserMsgCb(unsigned long context_data, void * msg);
	/**
	 * @brief 超时回调
	 */
	static int32 _thrdsTimeOut(unsigned long context_data, unsigned long timer_user_data, HTIMERID timerid);


	/**
	 * @brief process tcp input,在work_loop所在的线程中被呼叫
	 */
	static int32 _EpollInputTcp(void * context_data, int32 fd);
	/**
	 * @brief process accept event,在work_loop所在的线程中被呼叫
	 */
	static int32 _EpollInputAccept(void * context_data, int32 fd);
	/**
	 * @brief process output 在work_loop所在的线程中被呼叫
	 */
	static int32 _EpollOutputTcp(void * context_data, int32 fd);
	/**
	 * @brief process udp input,在work_loop所在的线程中被呼叫
	 */
	static int32 _EpollInputUdp(void * context_data, int32 fd);
	/**
	 * @brief process err 在work_loop所在的线程中被呼叫
	 */
	static int32 _EpollErr(void * context_data, int32 fd);


	/**
	 * @brief 添加至 map_thrd_ need lock
	 */
	int32 _addToMapThrd(int32 fd, SOCKET_FD_TYPE_T fd_type, uint32& thrd_index, int32 fd_master = -1);
	/**
	 * @brief 删除 need lock
	 */
	int32 _delFromMapThrd(int32 fd);
	/**
	 * @brief 获取fd的上下文数据 need lock
	 */
	int32 _getFdInfoFromMapThrd(int32 fd, fd_info_t*& i);
	int32 _getFdThrd(int32 fd, HMYLISTERNER& hlsn_);

	/**
	 * @brief 设置成自动加入或非自动加入监听
	 */
	int32 _set_accept_auto_add_to_listern_or_not(int32 fd, int32 bauto_add);

	/**
	 * @brief 推入线程池内部线程之间的消息
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
	 * @brief 线程个数
	 */
	uint32 real_thrd_count_;
	/**
	 * @brief 最大消息数量
	 */
	uint32 max_msg_count_;

	/**
	 * @brief epoll api wrapper
	 */
	CMyEPoll * epoll_;

	/*! 发送缓冲区保留大小,小于此值,则将缓冲区拷贝到缓冲区最前面,避免缓冲区恶性增长 */
	uint32 bufsz_reserve_;

	/**
	 * @brief 内存池句柄
	 */
	HMYMEMPOOL hm_;
};

#endif



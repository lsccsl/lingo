/**
* @file tcpclient_pool.h
* @brief tcp client connect pool mgr module
* @author linsc
* @blog http://blog.csdn.net/lsccsl
*/
#ifndef __TCPCLIENT_POOL_H__
#define __TCPCLIENT_POOL_H__

#include <string>
#include <list>
#include <set>
#include <pthread.h>
#include <semaphore.h>
#include "type_def.h"

class tcpclient_send_end_cb
{
public:

	/**
	 * @brief constructor
	 */
	tcpclient_send_end_cb(){}

	/**
	 * @brief destructor
	 */
	virtual ~tcpclient_send_end_cb(){}

	/**
	 * @brief 发送结束时的回调
	 * @param 返回0:成功 返回非0:则该链接不回收,直接关闭
	 */
	virtual const int32 send_end_call_back(int32 fd, const void * context_data, const void * user_data) const = 0;
};

/**
 * @brief 适合于终端发送,服务端响应类型通迅(服务端不会主动发消息)
 */
class tcpclient_pool
{
public:

	/**
	 * @brief constructor
	 * @param srv_ip:tcp服务端的ip
	 * @param srv_port:服务端port
	 * @param max_conn_count:连接上限
	 */
	tcpclient_pool(const int8 * srv_ip, const uint32 srv_port, const uint32 max_conn_count = 10, const uint32 min_conn_count = 5);

	/**
	 * @brief destructor
	 */
	~tcpclient_pool();

	/**
	 * @brief send data 将连接数达到上限时,将阻塞至有空闲连接为止
	 */
	int32 send_data(const void * buf, uint32 buf_sz, const uint32 count = 3, const uint32 timeout_second = 10,
		const tcpclient_send_end_cb * cb = NULL, const void * context_data = NULL, const void * user_data = NULL);

	/**
	* @brief read data,为一个功能函数
	*/
	static int32 read_data(int32 fd, const void * buf, uint32 buf_sz, uint32 count = 3, uint32 timeout_second = 10);

	/**
	 * @brief 获取srvip port
	 */
	const std::string& srv_ip(){ return this->srv_ip_; }
	const uint32 srv_port(){ return this->srv_port_; }

	/**
	* @brief for debug
	*/
	void view();

private:

	/**
	 * @brief 取出或者生成一个tcp链接
	 */
	int32 _get_conn(int32& fd);

	/**
	 * @brief 回取tcp链接
	 */
	int32 _recyc_fd(int32 fd);

private:

	/**
	* @brief 连接相关的信息
	*/
	std::string srv_ip_;
	uint32 srv_port_;
	uint32 max_conn_count_;
	uint32 min_conn_count_;

	/**
	* @brief 连接池
	*/
	std::set<int32> s_fd_;
	std::list<int32> lst_idle_fd_;
	pthread_mutex_t fd_protector_;
	sem_t sem_fd_;
};

#endif






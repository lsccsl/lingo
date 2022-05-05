/**
 * @file channel.cpp 
 * @brief 通道抽象,方便在线程池里回调时,统一接口
 *
 * @author linshaochuan
 */
#include "channel.h"

#ifndef WIN32
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <signal.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <errno.h>
#include <unistd.h>
#include <assert.h>
#include <fcntl.h>

#include <unistd.h>
#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <stdarg.h>
#include <errno.h>
#include <fcntl.h>
#include <pthread.h>
#include <signal.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/types.h>
#include <termios.h>
#include <time.h>
#include <unistd.h>
#include <string.h>
#include <sys/select.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/epoll.h>
#include <sys/timeb.h>
#include <netdb.h>
#include <netinet/tcp.h>
#include <sys/un.h>
#include <sys/time.h>
#else
#include <winsock2.h>
#include <io.h>
#include <mstcpip.h>
#include <ws2ipdef.h>
#define socklen_t int
#define close closesocket
#endif
#include <assert.h>
#include <errno.h>

#include "mylogex.h"

typedef struct sockaddr SA;


/**
 * @brief 获取socket name
 */
int32 CChannel::getSocketName(int32 sock, std::string& ip, uint32& port)
{
	struct sockaddr_in sa = {0};
	socklen_t len = sizeof(sa);
	getsockname(sock, (sockaddr*)&sa, &len);

	int8 * pctemp = inet_ntoa(sa.sin_addr);
	ip = pctemp ? pctemp : "null";
	port = ntohs(sa.sin_port);

	return 0;
}
int32 CChannel::getPeerName(int32 sock, std::string& ip, uint32& port)
{
	struct sockaddr_in sa = {0};
	socklen_t len = sizeof(sa);
	getpeername(sock, (sockaddr*)&sa, &len);

	int8 * pctemp = inet_ntoa(sa.sin_addr);
	ip = pctemp ? pctemp : "null";
	port = ntohs(sa.sin_port);

	return 0;
}

/**
 * @brief 设成非阻塞
 */
int32 CChannel::set_no_block(int32 fd)
{
#ifdef WIN32
	unsigned long iMode = 1;
	ioctlsocket(fd, FIONBIO, &iMode);
#else
    int iOpts = fcntl(fd, F_GETFL);
    if (iOpts < 0)
        return -1;

    iOpts = iOpts | O_NONBLOCK;

    if (fcntl(fd, F_SETFL, iOpts) < 0)
        return -1;
#endif

	return 0;
}
int32 CChannel::set_block(int32 fd)
{
#ifdef WIN32
#else
    int iOpts = fcntl(fd, F_GETFL);
    if (iOpts < 0)
        return -1;

    iOpts = iOpts & (~O_NONBLOCK);

    if (fcntl(fd, F_SETFL, iOpts) < 0)
        return -1;
#endif

	return 0;
}

/**
* @brief 设成立即关闭
*/
int32 CChannel::set_no_linger(int32 fd)
{
	struct linger lng = {1,0};
	if(0 != setsockopt(fd, SOL_SOCKET, SO_LINGER, (char *) &lng, sizeof(lng)))
		return -1;

	return 0;
}

/**
* @brief tcp链接保活
*/
int32 CChannel::keep_alive(int32 fd, uint32 idle, uint32 interval, uint32 retry_count)
{
#ifndef WIN32
	int32 keep_alive = 1;
	int32 ret = setsockopt(fd, SOL_SOCKET, SO_KEEPALIVE, (void*)&keep_alive, sizeof(keep_alive));
	if (ret == -1)
	{
		return -1;
	}

	ret = setsockopt(fd, SOL_TCP, TCP_KEEPIDLE, (void *)&idle, sizeof(idle));
	if (ret == -1)
	{
		return -2;
	}

	ret = setsockopt(fd, SOL_TCP, TCP_KEEPINTVL, (void *)&interval, sizeof(interval));
	if (ret == -1)
	{
		return -3;
	}

	ret = setsockopt(fd, SOL_TCP, TCP_KEEPCNT, (void *)&retry_count, sizeof(retry_count));
	if (ret == -1)
	{
		return -4;
	}
#else

	BOOL bKeepAlive = TRUE;
	int nRet = setsockopt(fd, SOL_SOCKET, SO_KEEPALIVE, (char*)&bKeepAlive, sizeof(bKeepAlive));
	if (nRet == SOCKET_ERROR)
		return -1;

	// set KeepAlive parameter
	tcp_keepalive alive_in;
	tcp_keepalive alive_out;
	alive_in.keepalivetime = idle * 1000;
	alive_in.keepaliveinterval = interval * 1000;
	alive_in.onoff = TRUE;
	unsigned long ulBytesReturn = 0;
	nRet = WSAIoctl(fd, SIO_KEEPALIVE_VALS, &alive_in, sizeof(alive_in), &alive_out, sizeof(alive_out), &ulBytesReturn, NULL, NULL);
	if (nRet == SOCKET_ERROR)
		return -2;

#endif
	return 0;
}

/**
* @brief 获取错误码
*/
int32 CChannel::get_socket_err()
{
#ifdef WIN32
	return WSAGetLastError();
#else
	return errno;
#endif
}

/**
 * @brief open tcp
 */
int32 CChannel::TcpOpen(const int8 * ip, uint32 port, uint32 max_conn)
{
	int32 fd = socket(PF_INET, SOCK_STREAM, 0);
	if (-1 == fd)
	{
		MYLOG_WARN(("socket err:%d", errno));
		return -1;
	}
	
	/* reuse port */
	int32 reuse = 1;
	setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, (int8 *)&reuse, sizeof(reuse));

	struct sockaddr_in addr_srv = {0};
	addr_srv.sin_family = AF_INET;
	if(NULL == ip)
		addr_srv.sin_addr.s_addr = INADDR_ANY;
	else
		addr_srv.sin_addr.s_addr = inet_addr(ip);
	addr_srv.sin_port = htons(port);

	if (-1 == bind(fd, (sockaddr*)&addr_srv, sizeof(addr_srv)))
	{
		MYLOG_WARN(("bind err:%d [%s:%d]", errno, ip ? ip : "NULL", port));
		goto CChannel_TcpOpen_;
	}

	if(-1 == listen(fd, max_conn))
	{
		MYLOG_WARN(("listern err:%d", errno));
		goto CChannel_TcpOpen_;
    }

	MYLOG_DEBUG(("open tcp socket %s:%d suc fd:%d", ip ? ip : "NULL", port, fd));

	return fd;

CChannel_TcpOpen_:

	close(fd);

    return -1;
}

/**
 * @brief tcp accept new connection
 */
int32 CChannel::TcpAccept(int32 fd_srv, int8 * ip, int32 ipsz, uint32 * port)
{
	struct sockaddr_in saddr;
	int sadd_sz = sizeof(saddr);

	int32 fd = ::accept(fd_srv, (struct sockaddr *)&saddr, (socklen_t*)&sadd_sz);

	MYLOG_DEBUG(("accept return:%d", fd));
	
	if(-1 == fd)
		return -1;

#ifdef WIN32
	#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
#endif
	if(inet_ntoa(saddr.sin_addr) && ip && ipsz)
		strncpy(ip, inet_ntoa(saddr.sin_addr), ipsz);

	if(port)
		*port = ntohs(saddr.sin_port);
		
	return fd;
}

/**
 * @brief tcp connect to srv
 */
int32 CChannel::TcpConnect(const int8 * srv_ip, uint32 srv_port, uint32 time_out, uint32 wait_count, int32 b_no_linger)
{
	int32 fd = socket(AF_INET, SOCK_STREAM, 0);

	if(b_no_linger)
		CChannel::set_no_linger(fd);

	int32 reuse = 1;
	setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, (int8 *)&reuse, sizeof(reuse));

	if(fd < 0)
		return -1;

	struct sockaddr_in saddr;
	if(0 == srv_port || NULL == srv_ip)
		return -1;

	memset(&saddr, 0, sizeof(saddr));

	saddr.sin_family = AF_INET;
	saddr.sin_addr.s_addr = inet_addr(srv_ip);
	saddr.sin_port = htons(srv_port);

	if(0 == time_out || 0 == wait_count)
	{
		/*
		* If the connection or binding succeeds, zero is returned.  
		* On error, -1 is returned, and errno is set appropriately.
		*/
		if(0 != connect(fd, (sockaddr *)&saddr, sizeof(saddr)))
		{
			close(fd);
			return -1;
		}

		return fd;
	}

	/* 设成非阻塞 */
	CChannel::set_no_block(fd);

	/* 发起异步连接 */
	if(connect(fd, (sockaddr *)&saddr, sizeof(saddr)))
	{
		while(wait_count)
		{
			wait_count --;
			struct timeval tv;
			fd_set fds;

			tv.tv_sec = time_out;
			tv.tv_usec = 0;

			FD_ZERO(&fds);
			FD_SET(fd, &fds);

			/* 等响应 */
			int32 r = select((int)fd + 1, NULL, &fds, NULL, &tv);
			if(r > 0)
			{
				/* sock可写,说明连接成功了 */
				return fd;
			}
			if(r < 0)
			{
				/* 参考gsoap, i don't know why errno != EINTR means not fail */
				if(errno != EINTR)
				{
					close(fd);
					return -1;
				}
				else
					continue;
			}
			else if(0 == r)
			{
				/* 没事件,再试一次,wait_count-- */
				continue;
			}
		}
	}

	close(fd);
	return -1;
}

/**
 * @brief tcp read
 * @return >0 : 读取到的数据量 0 : 无数据 <0 : 出错 
 */
int32 CChannel::TcpRead(int32 fd, void * buf, uint32 buf_sz)
{
	if(NULL == buf || 0 == buf_sz || fd < 0)
	{
		MYLOG_WARN(("err param or invalid fd:%d", fd));
		return -1;
	}

	int32 ret = recv(fd, (int8 *)buf, (int32)buf_sz, 0);
	if(ret > 0)
		return ret;
	else if(0 == ret)
		return -1;
	else
	{
#ifdef WIN32
		/* fuck vc,why you so... */
		if(WSAEWOULDBLOCK == WSAGetLastError())
			return 0;
#else
		if(EAGAIN == errno)
			return 0;
#endif
		return -1;
	}
}

/**
 * @brief tcp select read
 */
int32 CChannel::TcpSelectRead(int32 fd, void * buf, uint32 buf_sz, uint32 time_out, uint32 count, int * last_err)
{
	if(NULL == buf || 0 == buf_sz ||  fd < 0)
	{
		MYLOG_WARN(("err param or invalid fd:%d", fd));
		return -1;
	}

	int32 current_read = 0;
	int32 ret = 0;

	while((current_read < buf_sz) && count)
	{
		count -= 1;
#if 0
		/*These calls return the number of bytes received,
		*or -1 if an error occurred.
		*The return value will be 0 when the peer has performed an orderly shutdown.*/
		ret = recv(fd, (char*)buf + current_read, buf_sz - current_read, 0);
		if (ret <= 0)
			return -3;
		current_read += ret;
#else
		ret = CChannel::TcpRead(fd, (char*)buf + current_read, buf_sz - current_read);
		if (last_err)
			*last_err = WSAGetLastError();
		if (ret < 0)
		{
			if (last_err)
				*last_err = WSAGetLastError();
			return -1;
		}
		current_read += ret;
		if (current_read >= buf_sz)
			return current_read;
#endif

        fd_set fdWatch;
        struct timeval tvOut;

		FD_ZERO(&fdWatch);
		FD_SET(fd, &fdWatch);

		tvOut.tv_sec = time_out;
		tvOut.tv_usec = 0;

		/*
		*On success, select() and pselect() return the number of file descriptors contained in the three returned descriptor sets (that is, the total number of bits that are set in  readfds,  writefds,  exceptfds) 
		*which may be zero if the timeout expires before anything interesting happens.
		*On error, -1 is returned, and errno is set appropri-ately; the sets and timeout become undefined, so do not rely on their contents after an error.
		*/
		int32 r = select(fd + 1, &fdWatch, NULL, NULL, &tvOut);
		if (r < 0)
		{
			if (last_err)
				*last_err = WSAGetLastError();
			return -2;
		}
		else if (0 == r)
			continue;
		else if(!FD_ISSET(fd, &fdWatch))
			continue;	
	}

	return current_read;
}

/**
 * @brief tcp write
 * @return >0 : 读取到的数据量 0 : 无数据 <0 : 出错 
 */
int32 CChannel::TcpWrite(int32 fd, const void * buf, const uint32 buf_sz)
{
	if(NULL == buf || 0 == buf_sz || fd < 0)
	{
		MYLOG_WARN(("err param or invalid fd:%d", fd));
		return -1;
	}

	int32 ret = send(fd, (int8 *)buf, (int32)buf_sz, 0);
	if(ret >= 0)
		return ret;
	else
	{
#ifdef WIN32
		if (WSAEWOULDBLOCK == WSAGetLastError())
			return 0;
#else
		if (EAGAIN == errno)
			return 0;
#endif
		return -1;
	}
}

/**
 * @brief tcp select write
 * @param time_out : 每次写的超时
 * @param count : 重试次数
 */
int32 CChannel::TcpSelectWrite(int32 fd, const void * buf, const uint32 buf_sz, uint32 time_out, uint32 count)
{
	if(NULL == buf || 0 == buf_sz || fd < 0)
	{
		MYLOG_WARN(("err param or invalid fd:%d", fd));
		return -1;
	}

	int32 current_write = 0;
	int32 ret = 0;

	while((current_write < buf_sz) && (count != 0))
	{
		count -= 1;

		ret = CChannel::TcpWrite(fd, (char*)buf + current_write, buf_sz - current_write);
		if (ret < 0)
			return -2;
		current_write += ret;
		if (current_write >= buf_sz)
			return current_write;

        fd_set fdWatch;
        struct timeval tvOut;

		FD_ZERO(&fdWatch);
		FD_SET(fd, &fdWatch);

		tvOut.tv_sec = time_out;
		tvOut.tv_usec = 0;

		/*
		*On success, select() and pselect() return the number of file descriptors contained in the three returned descriptor sets (that is, the total number of bits that are set in  readfds,  writefds,  exceptfds) 
		*which may be zero if the timeout expires before anything interesting happens.
		*On error, -1 is returned, and errno is set appropri-ately; the sets and timeout become undefined, so do not rely on their contents after an error.
		*/
		int32 r = select(fd + 1, NULL, &fdWatch, NULL, &tvOut);
		if(r < 0)
			return -3;
		else if(0 == r)
			continue;
		else if(!FD_ISSET(fd, &fdWatch))
			continue;
#if 0
		/*
		* On success, these calls return the number of characters sent.
		* On error, -1 is returned, and errno is set appropriately.
		*/
		ret = send(fd, (char *)buf + current_write, buf_sz - current_write, 0);
		if(ret <= 0)
			return -3;

		current_write += ret;
#endif
	}

	return current_write;
}

/**
 * @brief close fd
 */
int32 CChannel::CloseFd(int32 fd)
{
	if(fd >= 0)
		close(fd);
	return 0;
}


/**
 * @brief 打开udp
 */
int32 CChannel::UdpOpen(const int8 * ip, uint32 port)
{
	int32 fd = (int32)socket(AF_INET, SOCK_DGRAM, 0);

	if(-1 == fd)
		return -1;

	int32 reuse = 1;
	setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, (int8 *)&reuse, sizeof(reuse));

	struct sockaddr_in stSockAddr;
	memset(&stSockAddr, 0, sizeof(stSockAddr));

	stSockAddr.sin_family = AF_INET;
	if(NULL == ip)
		stSockAddr.sin_addr.s_addr = INADDR_ANY;
	else
		stSockAddr.sin_addr.s_addr = inet_addr(ip);
	stSockAddr.sin_port = htons(port);

	if(0 == bind(fd, (SA *)&stSockAddr, sizeof(stSockAddr)))
		return fd;

	close(fd);

	return -1;
}

/**
 * @brief udp读
 */
int32 CChannel::UdpRead(int32 fd, void * buf, uint32 buf_sz,
	int8 * src_ip, uint32 src_ip_sz, uint16 * psrc_port)
{

	if(NULL == buf || 0 == buf_sz)
		return -1;

	int32 ret = 0;
	struct sockaddr_in saddr;
	int32 sadd_sz = sizeof(saddr);
	ret = recvfrom(fd, (int8 *)buf, (int32)buf_sz, 0, (SA *)&saddr, (socklen_t*)&sadd_sz);
	if(ret > 0)
	{
		if(psrc_port)
			*psrc_port = ntohs(saddr.sin_port);
		if(src_ip && src_ip_sz)
			strncpy(src_ip, inet_ntoa(saddr.sin_addr), src_ip_sz);
	}

	if(ret >= 0)
		return ret;
	else
	{
		if(EAGAIN == errno)
			return 0;
		return -1;
	}
}

/**
 * @brief udp写
 */
int32 CChannel::UdpWrite(int32 fd, const void * buf, const uint32 buf_sz,
	const int8 * dst_ip, uint16 dst_port)
{
	if(NULL == buf || 0 == buf_sz || NULL == dst_ip)
		return 0;

	if(-1 == fd)
		return -1;

	struct sockaddr_in saddr;
	memset(&saddr, 0, sizeof(saddr));
	saddr.sin_family = AF_INET;
	saddr.sin_addr.s_addr = inet_addr(dst_ip);
	saddr.sin_port = htons(dst_port);

	int32 ret = sendto(fd, (char *)buf, (int32)buf_sz, 0, (SA *)&saddr, sizeof(saddr));
	if(ret >= 0)
		return ret;
	else
	{
		if(EAGAIN == errno)
			return 0;
		return -1;
	}
}

/**
 * @brief
 */
static int ini_flag = 0;
void CChannel::init_sock()
{
#ifdef WIN32
	int err = -1;
	WSADATA wsadata;

	if(ini_flag)
		return;

	WSAStartup(0x0202, &wsadata);
	err = WSAGetLastError();
	if(err != 0)
		return;

	ini_flag = 1;
#endif
}
/**
* @brief 
*/
void CChannel::uninit_sock()
{
	ini_flag = 0;
#ifdef WIN32
	WSACleanup();
#endif
}





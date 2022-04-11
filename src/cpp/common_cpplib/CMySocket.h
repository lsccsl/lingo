

#ifndef __CMYSOCKET_H__
#define __CMYSOCKET_H__


#include <string>
#include "common_def.h"


/*! @class
********************************************************************************
<PRE>
类名称   :  CMyTcpSocket
功能     :  tcp socket wrapper
--------------------------------------------------------------------------------
备注     :  
典型用法 :  
--------------------------------------------------------------------------------
</PRE>
*******************************************************************************/
class CMyTcpSocket
{
public:

    enum
    {
        MAX_LISTEN = 10,
        MAX_EPOLL_SIZE = 10000,
    };

	/*!
	@brief 构造函数
	@param IN pIp:socket绑定的ip地址
	@param IN nPort:监听端口	
	@return
	********************************************************************/
	CMyTcpSocket(IN const char * pIp = NULL, IN const unsigned short nPort = 0);

	/*!
	@brief 构造函数
	@param IN ts 引用另外一个对象的sock fd
	@return
	********************************************************************/
	CMyTcpSocket(IN CMyTcpSocket& ts);

	/*!
	@brief 构造函数
	@param IN nSock 构造一个对象,fd为nSock
	@return
	********************************************************************/
	CMyTcpSocket(IN int nSock);


	/*!
	@brief 析构
	@return
	********************************************************************/
	~CMyTcpSocket();

	/*!
	@brief 打开socket
	@param IN pIp:socket绑定的ip地址
	@param IN nPort:监听端口	
	@return 0:成功 -1:失败
	********************************************************************/
	int Open(IN const char * pIp = NULL, IN const unsigned short nPort = 0);

	/*!
	@brief 打开socket
	@param IN pIp 要连接的host ip
	@param IN nPort 要连接的host port
	@return 0:成功 -1:失败
	********************************************************************/
	int Connect(IN const char * pIp, IN const unsigned short nPort);

	/*!
	@brief 关闭socket
	@return 0:成功 -1:失败
	********************************************************************/
	int Close();	

	/*!
	@brief tcp accept
	@param IN  out_tcp_socket存储accept得到的tcp fd信息
	@param OUT pIp:返回终端的ip	
	@param IN  nIpSz:pIp的大小
	@param OUT pnPort:返回端口
	@return 0:成功 -1:失败
	********************************************************************/
	int Accept(OUT CMyTcpSocket& out_tcp_socket, 
		OUT char * pIp, 
		IN int nIpSz, 
		OUT unsigned short * pnPort);

	/*!
	@brief tcp read
	@param OUT pbuf存储得到的信息
	@param IN  nBufSz pbuf的大小
	@return >0:读到的字节数 -1:失败 0:对方关闭了socket
	********************************************************************/
	int Read(OUT void * pbuf, IN size_t nBufSz);

	/*!
	@brief tcp read
	@param OUT pbuf存储得到的信息
	@param IN  nBufSz pbuf的大小
	@param IN  nTimeOut select read 超时
	@return >0:读到的字节数 >0:失败 0:对方关闭了socket
	********************************************************************/
	enum{
		READ_ERR = -1,
		READ_TIMEOUT = -2,
	};
	int Read(OUT void * pbuf, IN size_t nBufSz, int nTimeOut);

	/*!
	@brief tcp write
	@param IN pbuf要写入的缓冲区
	@param IN  nBufSz pbuf的大小
	@return 0:成功 -1:失败
	********************************************************************/
	int Write(IN const void * pbuf, IN size_t nBufSz);

	/*!
	@brief tcp write
	@return tcp fd句柄
	********************************************************************/
	int GetFd();

	/*!
	@brief 设成非阻塞
	@return 0:成功 -1:失败
	********************************************************************/
	int SetToNoBlock();

protected:

	/*! socket fd  */
	int m_nTcpFd;

	/*! m_nTcpFd 是否归当前实例所有  */
	bool m_bFdIsMine;
};


/*! @class
********************************************************************************
<PRE>
类名称   :  CMyUnixSocket
功能     :  封装unix socket
--------------------------------------------------------------------------------
备注     :  
典型用法 :  
--------------------------------------------------------------------------------
</PRE>
*******************************************************************************/
class CMyUnixSocket
{
public:

	/*!
	@brief 构造
	@param IN pcPath unix path
	@return 无
	********************************************************************/
	CMyUnixSocket(IN const char * pcPath = NULL);

	/*!
	@brief 析构
	@return 无
	********************************************************************/
	~CMyUnixSocket();

	/*!
	@brief 打开unix sock
	@param IN pcPath unix path
	@return 0:成功 -1:失败
	********************************************************************/
	int Open(const char * pcPath);

	/*!
	@brief 关闭sock
	@return 0:成功 -1:失败
	********************************************************************/
	int Close();	

	/*!
	@brief 写unix sock
	@param IN pbuf要写入的缓冲区
	@param IN nBufSz缓冲区的大小
	@param IN pcTargetPath目标的unix path
	@return 0:成功 -1:失败
	********************************************************************/
	int Write(IN const void * pbuf, 
		IN size_t nBufSz, 
		IN const char * pcTargetPath);

	/*!
	@brief 读unix sock
	@param pbuf:缓冲区
	@param nBufSz:缓冲区大小
	@param pcFromPath:记录消息源unix socket的路径
	@param nFromPathSz:from_path的大小
	@return >0:读到的字节数 -1:失败 0:对方关闭了socket
	********************************************************************/
	int Read(OUT void * pbuf, 
		IN size_t nBufSz, 
		char * pcFromPath, 
		size_t nFromPathSz);	

	/*!
	@brief 读unix sock
	@param pbuf:缓冲区
	@param nBufSz:缓冲区大小
	@param pcFromPath:记录消息源unix socket的路径
	@param nFromPathSz:from_path的大小
	@param nTimeOut 读取超时
	@return >0:读到的字节数 -1:失败 0:对方关闭了socket
	********************************************************************/
	enum{
		READ_ERR = -1,
		READ_TIMEOUT = -2,
	};
	int Read(OUT void * pbuf, 
		IN size_t nBufSz, 
		char * pcFromPath, 
		size_t nFromPathSz,
		int nTimeOut);	

	/*!
	@brief 获取unix sock fd
	@return 返回unix sock fd
	********************************************************************/
	int GetFd();

	/*!
	@brief 获取unix path
	@return 返回unix path
	********************************************************************/
	std::string& GetUnixPath();

	/*!
	@brief 获取unix path
	@return 0:成功 -1:失败
	********************************************************************/
	int SetToNoBlock();

protected:
	
	/*! socket fd  */
	int m_nUnixFd;

	/*! unix路径 */
	std::string m_strUnixPath;
};


#endif

#include "CMySocket.h"


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
#else
#include <winsock2.h>
#include <io.h>
#define socklen_t int
#endif

#include "mylogex.h"


#define SA struct sockaddr


/*!
@brief ���캯��
@param IN pIp:socket�󶨵�ip��ַ
@param IN nPort:�����˿�	
@return
********************************************************************/
CMyTcpSocket::CMyTcpSocket(IN const char * pIp /*= NULL*/, 
	IN const unsigned short nPort /*= 0*/):
		m_nTcpFd(-1)
{
	if(NULL != pIp && 0 != nPort)
		this->Open(pIp, nPort);
}

/*!
@brief ���캯��
@param IN ts ��������һ�������sock fd
@return
********************************************************************/
CMyTcpSocket::CMyTcpSocket(CMyTcpSocket& ts):m_nTcpFd(ts.GetFd()),m_bFdIsMine(false)
{}

/*!
@brief ���캯��
@param IN nSock ����һ������,fdΪnSock
@return
********************************************************************/
CMyTcpSocket::CMyTcpSocket(IN int nSock):m_nTcpFd(nSock),m_bFdIsMine(false)
{}

/*!
@brief ����
@return
********************************************************************/
CMyTcpSocket::~CMyTcpSocket()
{
	this->Close();
}

/*!
@brief ��socket
@param IN pIp:socket�󶨵�ip��ַ
@param IN nPort:�����˿�	
@return
********************************************************************/
int CMyTcpSocket::Open(IN const char * pIp /*= NULL*/, 
	IN const unsigned short nPort /*= 0*/)
{
	if(NULL == pIp || 0 == nPort)
	{
		MYLOG_WARN(("CMyTcpSocket::open err param"));
		return -1;
	}
	
	if(this->m_nTcpFd > 0)
	{
		MYLOG_DEBUG(("socket has been open:%d, please close first..", this->m_nTcpFd));
		return -1;
	}

	this->m_nTcpFd = socket(PF_INET, SOCK_STREAM, 0);
	if (-1 == this->m_nTcpFd)
	{
		MYLOG_WARN(("����socket��������(%d)", errno));
		return -1;
	}
	
	/* ���ö˿ڸ��ã������ϴ�û�������ɱ����޷��󶨼����� */
	int iYes = 1;
	setsockopt(this->m_nTcpFd, SOL_SOCKET, SO_REUSEADDR, (char *)&iYes, sizeof(int));

	// ����SOCKET
	struct sockaddr_in addrSvr = {0};
	addrSvr.sin_family = AF_INET;
	addrSvr.sin_addr.s_addr = INADDR_ANY;
	addrSvr.sin_port = htons(nPort);

	if (-1 == bind(this->m_nTcpFd, (sockaddr*)&addrSvr, sizeof(addrSvr)))
	{
		MYLOG_WARN(("����bind��������(%d)", errno));
		goto __CMyTcpSocket_Open_err_;
	}

	if(-1 == listen(this->m_nTcpFd, MAX_LISTEN))
	{
		MYLOG_WARN(("����listen��������(%d)", errno));
		goto __CMyTcpSocket_Open_err_;
    }

	this->m_bFdIsMine = true;

	MYLOG_DEBUG(("open tcp socket %s:%d suc fd:%d %d", 
		pIp, nPort, this->m_nTcpFd, this->m_bFdIsMine));

    return 0;

__CMyTcpSocket_Open_err_:

	this->Close();
	return -1;
}

/*!
@brief ��socket
@param IN pIp Ҫ���ӵ�host ip
@param IN nPort Ҫ���ӵ�host port
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyTcpSocket::Connect(IN const char * pIp, IN const unsigned short nPort)
{
	if(-1 != this->m_nTcpFd)
		this->Close();

	this->m_nTcpFd = socket(AF_INET, SOCK_STREAM, 0);
	this->m_bFdIsMine = true;

	struct sockaddr_in saddr;
	if(0 == nPort || NULL == pIp)
		return -1;

	memset(&saddr, 0, sizeof(saddr));

	saddr.sin_family = AF_INET;
	saddr.sin_addr.s_addr = inet_addr(pIp);
	saddr.sin_port = htons(nPort);

	/*
	* If the connection or binding succeeds, zero is returned.  
	* On error, -1 is returned, and errno is set appropriately.
	*/
	if(0 != connect(this->m_nTcpFd, (sockaddr *)&saddr, sizeof(saddr)))
		return -1;

	return 0;

}

/*!
@brief �ر�socket
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyTcpSocket::Close()
{
	if(!this->m_bFdIsMine)
	{
		this->m_nTcpFd = -1;
		return 0;
	}

	if(-1 == this->m_nTcpFd)
		return -1;
	
	MYLOG_DEBUG(("close tcp sock:%d", this->m_nTcpFd));


#ifdef WIN32
	::closesocket(this->m_nTcpFd);
#else
	::close(this->m_nTcpFd);
#endif
	this->m_nTcpFd = -1;
	return 0;
}

/*!
@brief tcp accept
@param IN  out_tcp_socket�洢accept�õ���tcp fd��Ϣ
@param OUT pIp:�����ն˵�ip	
@param IN  nIpSz:pIp�Ĵ�С
@param OUT pnPort:���ض˿�
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyTcpSocket::Accept(OUT CMyTcpSocket& out_tcp_socket, 
		OUT char * pIp, 
		IN int nIpSz, 
		OUT unsigned short * pnPort)
{
	out_tcp_socket.Close();
	
	struct sockaddr_in saddr;
	int sadd_sz = sizeof(saddr);

	out_tcp_socket.m_nTcpFd = ::accept(this->m_nTcpFd, (struct sockaddr *)&saddr, (socklen_t*)&sadd_sz);
	out_tcp_socket.m_bFdIsMine = true;

	MYLOG_DEBUG(("accept return:%d", out_tcp_socket.m_nTcpFd));
	
	if(-1 == out_tcp_socket.m_nTcpFd)
		return -1;

#ifdef WIN32
	#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
#endif
	if(inet_ntoa(saddr.sin_addr) && pIp && nIpSz)
		strncpy(pIp, inet_ntoa(saddr.sin_addr), nIpSz);

	if(pnPort)
		*pnPort = ntohs(saddr.sin_port);
		
	return 0;
}

/*!
@brief tcp read
@param OUT pbuf�洢�õ�����Ϣ
@param IN  nBufSz pbuf�Ĵ�С
@return >0:�������ֽ��� -1:ʧ�� 0:�Է��ر���socket
********************************************************************/
int CMyTcpSocket::Read(OUT void * pbuf, IN size_t nBufSz)
{
	if(NULL == pbuf || 0 == nBufSz || -1 == this->m_nTcpFd)
	{
		MYLOG_WARN(("err param or invalid fd:%d", this->m_nTcpFd));
		return -1;
	}

	return recv(this->m_nTcpFd, (char *)pbuf, (int)nBufSz, 0);
}

/*!
@brief tcp read
@param OUT pbuf�洢�õ�����Ϣ
@param IN  nBufSz pbuf�Ĵ�С
@param IN  nTimeOut select read ��ʱ
@return >0:�������ֽ��� -1:ʧ�� 0:�Է��ر���socket
********************************************************************/
int CMyTcpSocket::Read(OUT void * pbuf, IN size_t nBufSz, int nTimeOut)
{
	fd_set ssset;
	FD_ZERO(&ssset);		
	FD_SET(this->m_nTcpFd, &ssset);

    struct timeval tv;
    tv.tv_sec = nTimeOut;
    tv.tv_usec = 0;

	/*
	* return the number of file descriptors contained in the three returned descriptor sets (that is, the total number of bits that are set
    *   in  readfds,  writefds,  exceptfds)  
    * which may be zero if the timeout expires before anything interesting happens.  
    * On error, -1 is returned,
    */
	int ret = select(this->m_nTcpFd + 1, &ssset, NULL, NULL, &tv);
	
	if(0 == ret)
		return READ_TIMEOUT;
	if(ret < 0)
		return READ_ERR;

	if(!FD_ISSET(this->m_nTcpFd, &ssset))
		return READ_ERR;

	/*
	* These calls return the number of bytes received, 
	* or -1 if an error occurred.
	* The return value will be 0 when the peer has performed an orderly shutdown.
	*/
	return recv(this->m_nTcpFd, (char *)pbuf, (int)nBufSz, 0);
}

/*!
@brief tcp write
@param IN pbufҪд��Ļ�����
@param IN  nBufSz pbuf�Ĵ�С
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyTcpSocket::Write(IN const void * pbuf, IN size_t nBufSz)
{
	if(NULL == pbuf || 0 == nBufSz || -1 == this->m_nTcpFd)
	{
		MYLOG_WARN(("err param or invalid fd:%d", this->m_nTcpFd));
		return -1;
	}

	return send(this->m_nTcpFd, (char *)pbuf, (int)nBufSz, 0);
}

/*!
@brief tcp write
@return tcp fd���
********************************************************************/
int CMyTcpSocket::GetFd()
{
	return this->m_nTcpFd;
}

/*!
@brief ��ɷ�����
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyTcpSocket::SetToNoBlock()
{
#ifdef WIN32
	unsigned long iMode = 1;
	ioctlsocket(this->m_nTcpFd, FIONBIO, &iMode);
#else
    int iOpts = fcntl(this->m_nTcpFd, F_GETFL);
    if (iOpts < 0)
        return -1;

    iOpts = iOpts | O_NONBLOCK;

    if (fcntl(this->m_nTcpFd, F_SETFL, iOpts) < 0)
        return -1;
#endif

	return 0;
}


/*!
@brief ���캯��
@param IN pIp:socket�󶨵�ip��ַ
@param IN nPort:�����˿�	
@return ��
********************************************************************/
CMyUnixSocket::CMyUnixSocket(IN const char * pcPath /*= NULL*/):m_nUnixFd(-1)
{
	if(NULL != pcPath)
		this->Open(pcPath);
}

/*!
@brief ����
@return ��
********************************************************************/
CMyUnixSocket::~CMyUnixSocket()
{
	this->Close();
}

/*!
@brief ��unix sock
@param IN pcPath unix path
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::Open(const char * pcPath)
{
#ifndef WIN32
	if(NULL == pcPath)
	{
		MYLOG_WARN(("unix socket open err param..."));
		return -1;
	}
	
	this->m_nUnixFd = socket(AF_LOCAL/*AF_UNIX*/, SOCK_DGRAM, 0);
	int rt = 0;

	if (this->m_nUnixFd < 0)
	{
		MYLOG_WARN(("fail get local socket fd"));
		return -1;
	}

	struct sockaddr_un unixaddr;
	memset(&unixaddr, 0, sizeof(unixaddr));
	unixaddr.sun_family = /*AF_UNIX*/AF_LOCAL;
	unlink(pcPath);
	strncpy(&unixaddr.sun_path[0], pcPath, strlen(pcPath) + 1);
	rt = bind(this->m_nUnixFd, (SA *) & unixaddr, sizeof(unixaddr));

	if (rt < 0)
	{
		MYLOG_WARN(("fail bind %s", pcPath));
		return -1;
	}

	MYLOG_DEBUG(("open tcp socket [%s] suc fd:%d", 
		pcPath, this->m_nUnixFd));
	
	this->m_strUnixPath = pcPath;
#endif
	return 0;
}

/*!
@brief �ر�sock
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::Close()
{
#ifndef WIN32
	if(-1 == this->m_nUnixFd)
		return -1;

	MYLOG_DEBUG(("close unix sock:%d", this->m_nUnixFd));

	::close(this->m_nUnixFd);
	this->m_nUnixFd = -1;
	unlink(this->m_strUnixPath.c_str());
#endif
	return 0;
}

/*!
@brief дunix sock
@param IN pbufҪд��Ļ�����
@param IN nBufSz�������Ĵ�С
@param IN pcTargetPathĿ���unix path
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::Write(IN const void * pbuf, 
		IN size_t nBufSz, 
		IN const char * pcTargetPath)
{
#ifdef WIN32
	return -1;
#else
	struct sockaddr_un unixaddr;

	bzero(&unixaddr, sizeof(unixaddr));
	unixaddr.sun_family = AF_LOCAL/*AF_UNIX*/;
	strncpy(&unixaddr.sun_path[0], pcTargetPath, strlen(pcTargetPath));

	return sendto(this->m_nUnixFd, (char *)pbuf, (int)nBufSz, 0, (SA *)&unixaddr, sizeof(unixaddr));
#endif
}

/*!
@brief ��unix sock
@param pbuf:������
@param nBufSz:��������С
@param pcFromPath:��¼��ϢԴunix socket��·��
@param nFromPathSz:from_path�Ĵ�С
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::Read(OUT void * pbuf, 
		IN size_t nBufSz, 
		char * pcFromPath, 
		size_t nFromPathSz)
{
#ifdef WIN32
	return -1;
#else
	struct sockaddr_un caddr;
	memset(&caddr, 0, sizeof(caddr));
	socklen_t clen = sizeof(caddr);

	/*
	* These calls return the number of bytes received, 
	* or -1 if an error occurred. 
	* The return value will be 0 when the peer has performed an orderly shutdown.
	*/
	int ret = recvfrom(this->m_nUnixFd, (char *)pbuf, (int)nBufSz, 0, (SA *)&caddr, (socklen_t*)&clen);

	if(pcFromPath && nFromPathSz)
		strncpy(pcFromPath, caddr.sun_path, nFromPathSz);

	return ret;
#endif
}

/*!
@brief ��unix sock
@param pbuf:������
@param nBufSz:��������С
@param pcFromPath:��¼��ϢԴunix socket��·��
@param nFromPathSz:from_path�Ĵ�С
@param nTimeOut ��ȡ��ʱ
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::Read(OUT void * pbuf, 
	IN size_t nBufSz, 
	char * pcFromPath, 
	size_t nFromPathSz,
	int nTimeOut)
{
	fd_set ssset;
	FD_ZERO(&ssset);		
	FD_SET(this->m_nUnixFd, &ssset);

    struct timeval tv;
    tv.tv_sec = nTimeOut;
    tv.tv_usec = 0;

	/*
	* return the number of file descriptors contained in the three returned descriptor sets (that is, the total number of bits that are set
    *   in  readfds,  writefds,  exceptfds)  
    * which may be zero if the timeout expires before anything interesting happens.  
    * On error, -1 is returned,
    */
	int ret = select(this->m_nUnixFd + 1, &ssset, NULL, NULL, &tv);
	
	if(0 == ret)
		return READ_TIMEOUT;
	if(ret < 0)
		return READ_ERR;

	if(!FD_ISSET(this->m_nUnixFd, &ssset))
		return READ_ERR;

	return this->Read(pbuf, nBufSz, pcFromPath, nFromPathSz);
}

/*!
@brief ��ȡunix sock fd
@return ����unix sock fd
********************************************************************/
int CMyUnixSocket::GetFd()
{
	return this->m_nUnixFd;
}

/*!
@brief ��ɷ�����
@return 0:�ɹ� -1:ʧ��
********************************************************************/
int CMyUnixSocket::SetToNoBlock()
{
#ifndef WIN32
    int iOpts = fcntl(this->m_nUnixFd, F_GETFL);
    if (iOpts < 0)
        return -1;

    iOpts = iOpts | O_NONBLOCK;

    if (fcntl(this->m_nUnixFd, F_SETFL, iOpts) < 0)
        return -1;
#endif
	return 0;
}

/*!
@brief ��ȡunix path
@return ����unix path
********************************************************************/
std::string& CMyUnixSocket::GetUnixPath()
{
	return this->m_strUnixPath;
}



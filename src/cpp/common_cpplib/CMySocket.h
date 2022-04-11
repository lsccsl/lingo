

#ifndef __CMYSOCKET_H__
#define __CMYSOCKET_H__


#include <string>
#include "common_def.h"


/*! @class
********************************************************************************
<PRE>
������   :  CMyTcpSocket
����     :  tcp socket wrapper
--------------------------------------------------------------------------------
��ע     :  
�����÷� :  
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
	@brief ���캯��
	@param IN pIp:socket�󶨵�ip��ַ
	@param IN nPort:�����˿�	
	@return
	********************************************************************/
	CMyTcpSocket(IN const char * pIp = NULL, IN const unsigned short nPort = 0);

	/*!
	@brief ���캯��
	@param IN ts ��������һ�������sock fd
	@return
	********************************************************************/
	CMyTcpSocket(IN CMyTcpSocket& ts);

	/*!
	@brief ���캯��
	@param IN nSock ����һ������,fdΪnSock
	@return
	********************************************************************/
	CMyTcpSocket(IN int nSock);


	/*!
	@brief ����
	@return
	********************************************************************/
	~CMyTcpSocket();

	/*!
	@brief ��socket
	@param IN pIp:socket�󶨵�ip��ַ
	@param IN nPort:�����˿�	
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Open(IN const char * pIp = NULL, IN const unsigned short nPort = 0);

	/*!
	@brief ��socket
	@param IN pIp Ҫ���ӵ�host ip
	@param IN nPort Ҫ���ӵ�host port
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Connect(IN const char * pIp, IN const unsigned short nPort);

	/*!
	@brief �ر�socket
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Close();	

	/*!
	@brief tcp accept
	@param IN  out_tcp_socket�洢accept�õ���tcp fd��Ϣ
	@param OUT pIp:�����ն˵�ip	
	@param IN  nIpSz:pIp�Ĵ�С
	@param OUT pnPort:���ض˿�
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Accept(OUT CMyTcpSocket& out_tcp_socket, 
		OUT char * pIp, 
		IN int nIpSz, 
		OUT unsigned short * pnPort);

	/*!
	@brief tcp read
	@param OUT pbuf�洢�õ�����Ϣ
	@param IN  nBufSz pbuf�Ĵ�С
	@return >0:�������ֽ��� -1:ʧ�� 0:�Է��ر���socket
	********************************************************************/
	int Read(OUT void * pbuf, IN size_t nBufSz);

	/*!
	@brief tcp read
	@param OUT pbuf�洢�õ�����Ϣ
	@param IN  nBufSz pbuf�Ĵ�С
	@param IN  nTimeOut select read ��ʱ
	@return >0:�������ֽ��� >0:ʧ�� 0:�Է��ر���socket
	********************************************************************/
	enum{
		READ_ERR = -1,
		READ_TIMEOUT = -2,
	};
	int Read(OUT void * pbuf, IN size_t nBufSz, int nTimeOut);

	/*!
	@brief tcp write
	@param IN pbufҪд��Ļ�����
	@param IN  nBufSz pbuf�Ĵ�С
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Write(IN const void * pbuf, IN size_t nBufSz);

	/*!
	@brief tcp write
	@return tcp fd���
	********************************************************************/
	int GetFd();

	/*!
	@brief ��ɷ�����
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int SetToNoBlock();

protected:

	/*! socket fd  */
	int m_nTcpFd;

	/*! m_nTcpFd �Ƿ�鵱ǰʵ������  */
	bool m_bFdIsMine;
};


/*! @class
********************************************************************************
<PRE>
������   :  CMyUnixSocket
����     :  ��װunix socket
--------------------------------------------------------------------------------
��ע     :  
�����÷� :  
--------------------------------------------------------------------------------
</PRE>
*******************************************************************************/
class CMyUnixSocket
{
public:

	/*!
	@brief ����
	@param IN pcPath unix path
	@return ��
	********************************************************************/
	CMyUnixSocket(IN const char * pcPath = NULL);

	/*!
	@brief ����
	@return ��
	********************************************************************/
	~CMyUnixSocket();

	/*!
	@brief ��unix sock
	@param IN pcPath unix path
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Open(const char * pcPath);

	/*!
	@brief �ر�sock
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Close();	

	/*!
	@brief дunix sock
	@param IN pbufҪд��Ļ�����
	@param IN nBufSz�������Ĵ�С
	@param IN pcTargetPathĿ���unix path
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int Write(IN const void * pbuf, 
		IN size_t nBufSz, 
		IN const char * pcTargetPath);

	/*!
	@brief ��unix sock
	@param pbuf:������
	@param nBufSz:��������С
	@param pcFromPath:��¼��ϢԴunix socket��·��
	@param nFromPathSz:from_path�Ĵ�С
	@return >0:�������ֽ��� -1:ʧ�� 0:�Է��ر���socket
	********************************************************************/
	int Read(OUT void * pbuf, 
		IN size_t nBufSz, 
		char * pcFromPath, 
		size_t nFromPathSz);	

	/*!
	@brief ��unix sock
	@param pbuf:������
	@param nBufSz:��������С
	@param pcFromPath:��¼��ϢԴunix socket��·��
	@param nFromPathSz:from_path�Ĵ�С
	@param nTimeOut ��ȡ��ʱ
	@return >0:�������ֽ��� -1:ʧ�� 0:�Է��ر���socket
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
	@brief ��ȡunix sock fd
	@return ����unix sock fd
	********************************************************************/
	int GetFd();

	/*!
	@brief ��ȡunix path
	@return ����unix path
	********************************************************************/
	std::string& GetUnixPath();

	/*!
	@brief ��ȡunix path
	@return 0:�ɹ� -1:ʧ��
	********************************************************************/
	int SetToNoBlock();

protected:
	
	/*! socket fd  */
	int m_nUnixFd;

	/*! unix·�� */
	std::string m_strUnixPath;
};


#endif

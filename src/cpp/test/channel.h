/**
 * @file channel.h 
 * @brief ͨ������,�������̳߳���ص�ʱ,ͳһ�ӿ�
 *
 * @author linshaochuan
 * @blog http://blog.csdn.net/lsccsl
 */
#ifndef __CHANNEL_H__
#define __CHANNEL_H__


#include <string>
#include "type_def.h"


/**
 * @brief ͨ������,�������̳߳���ص�ʱ,ͳһ�ӿ�
 */
class CChannel
{
public:

	/**
	 * @brief constructor
	 */
	CChannel(){}

	/**
	 * @brief destructor
	 */
	virtual ~CChannel(){}

	/**
	 * @brief read
	 * @param buf : ���ջ������׵�ַ
	 * @param buf_sz : ���ջ����С
	 * @param peerinfo : �Զ���Ϣ�������׵�ַ
	 * @param peerinfo_sz : peerinfo�Ĵ�С
	 * @return >0 : ��ȡ���������� 0 : ������ <0 : ���� 
	 */
	virtual int32 read(void * buf, uint32 buf_sz, void * peerinfo, uint32 peerinfo_sz) = 0;

	/**
	 * @brief wirte
	 * @param buf : ���ͻ������׵�ַ
	 * @param buf_sz : buf_sz��С
	 * @param peerinfo : �Զ���Ϣ�������׵�ַ
	 * @param peerinfo_sz : peerinfo�Ĵ�С
	 * @return >0 : д�뵽�������� 0 : ������ <0 : ���� 
	 */
	virtual int32 wirte(const void * buf, const uint32 buf_sz, const void * peerinfo, const uint32 peerinfo_sz) = 0;

	/**
	 * @brief get fd
	 */
	virtual int32 fd();

public:

	/**
	 * @brief
	 */
	static void init_sock();
	/**
	* @brief 
	*/
	static void uninit_sock();

	/**
	 * @brief ��ȡsocket name
	 */
	static int32 getSocketName(int32 sock, std::string& ip, uint32& port);
	static int32 getPeerName(int32 sock, std::string& ip, uint32& port);

	/**
	 * @brief ��ɷ�����
	 */
	static int32 set_no_block(int32 fd);
	static int32 set_block(int32 fd);

	/**
	* @brief ��������ر�
	*/
	static int32 set_no_linger(int32 fd);

	/**
	* @brief tcp���ӱ���
	*/
	static int32 keep_alive(int32 fd, uint32 idle = 10, uint32 interval = 10, uint32 retry_count = 10);

	/**
	* @brief ��ȡ������
	*/
	static int32 get_socket_err();

	/**
	 * @brief close fd
	 */
	static int32 CloseFd(int32 fd);

	/**
	 * @brief open tcp
	 * @return >=0 : suc  <0 : err 
	 */
	static int32 TcpOpen(const int8 * ip, uint32 port, uint32 max_conn);

	/**
	 * @brief tcp accept new connection
	 * @return >=0 : suc  <0 : err 
	 */
	static int32 TcpAccept(int32 fd_srv, int8 * ip = NULL, int32 ipsz = 0, uint32 * port = NULL);

	/**
	 * @brief tcp connect to srv
	 * @return >=0 : suc  <0 : err 
	 */
	static int32 TcpConnect(const int8 * srv_ip, uint32 srv_port, uint32 time_out = 10, uint32 wait_count = 3, int32 b_no_linger = 1);

	/**
	 * @brief tcp read
	 * @return >0 : ��ȡ���������� 0 : ������ <0 : ���� 
	 */
	static int32 TcpRead(int32 fd, void * buf, uint32 buf_sz);

	/**
	 * @brief tcp select read
	 * @param time_out : ÿ�ζ��ĳ�ʱ
	 * @param count : ���Դ���
	 */
	static int32 TcpSelectRead(int32 fd, void * buf, uint32 buf_sz, uint32 time_out = 3, uint32 count = 10);

	/**
	 * @brief tcp write
	 * @return >0 : д�뵽�������� 0 : ������ <0 : ���� 
	 */
	static int32 TcpWrite(int32 fd, const void * buf, const uint32 buf_sz);

	/**
	 * @brief tcp select write
	 * @param time_out : ÿ��д�ĳ�ʱ
	 * @param count : ���Դ���
	 */
	static int32 TcpSelectWrite(int32 fd, const void * buf, const uint32 buf_sz, uint32 time_out = 3, uint32 count = 10);


	/**
	 * @brief ��udp
	 */
	static int32 UdpOpen(const int8 * ip, uint32 port);

	/**
	 * @brief udp��
	 * @return >0 : ��ȡ���������� 0 : ������ <0 : ���� 
	 */
	static int32 UdpRead(int32 fd, void * buf, uint32 buf_sz,
		int8 * src_ip, uint32 src_ip_sz, uint16 * psrc_port);

	/**
	 * @brief udpд
	 * @return >0 : д�뵽�������� 0 : ������ <0 : ���� 
	 */
	static int32 UdpWrite(int32 fd, const void * buf, const uint32 buf_sz,
		const int8 * dst_ip, uint16 dst_port);
};


#endif



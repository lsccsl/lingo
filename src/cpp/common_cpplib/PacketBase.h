/**
 * @file PacketBase.h 
 * @brief ˽��Э�鱨�Ļ���
 *
 * @author linshaochuan
 * @blog http://blog.csdn.net/lsccsl
 */
#ifndef __PACKET_BASE_H__
#define __PACKET_BASE_H__


#include <vector>
#include <string>

#include "type_def.h"


/**
 * @brief ˽��Э�鱨�Ļ���
 */
class PacketBase
{
public:
	
	/**
	 * @brief ���������Ͷ���
	 */
	typedef std::vector<uint8> PKT_BUF;

	/**
	 * @brief ����
	 */
	PacketBase();

	/**
	 * @brief ����
	 */
	virtual ~PacketBase();

	/**
	 * @brief ��ȡ������ָ��λ�õĵ�ַ
	 * @param pos:λ��,0��ʾ�׵�ַ
	 */
	const PKT_BUF& GetBufRef()const{return this->vbuf_;}

	/**
	 * @brief ��ȡ���������ݴ�С
	 */
	uint32 GetBufSZ();

	/**
	 * @brief ���û�����������
	 */
	int32 SetBufCap(uint32 sz);

	/**
	 * @brief ��ȡ����������
	 */
	uint32 GetBufCap();

	/**
	 * @brief ��ȡλ��
	 */
	uint32 GetPos(){return this->pos_;}

	/**
	 * @brief ��ȡʣ��δ�������ֽ���
	 */
	uint32 GetLeftByteCount()
	{
		if(this->pos_ > this->buf_sz_)
		{
			return 0;
		}
		else
		{
			return this->buf_sz_ - this->pos_;
		}
	}


	/**
	 * @brief add 1 byte
	 */
	int32 Add1Byte(uint8 p);

	/**
	 * @brief add 2 byte
	 */
	int32 Add2Byte(uint16 p);

	/**
	 * @brief add 4 byte
	 */
	int32 Add4Byte(uint32 p);

	/**
	 * @brief add 8 byte
	 */
	int32 Add8Byte(const uint64& p);

	/**
	 * @brief add x byte
	 */
	int32 AddXByte(uint8 * p, uint32 sz);

	/**
	 * @brief add string
	 */
	int32 AddString(const int8 * p);
	int32 AddString(const std::string& s){
		return this->AddString(s.c_str());
	}


	/**
	 * @brief ����ĳ��λ�õ�ֵ 1Byte
	 */
	int32 Set1Byte(uint32 pos, uint8 p);

	/**
	 * @brief ����ĳ��λ�õ�ֵ 2Byte
	 */
	int32 Set2Byte(uint32 pos, uint16 p);

	/**
	 * @brief ����ĳ��λ�õ�ֵ 4Byte
	 */
	int32 Set4Byte(uint32 pos, uint32 p);

	/**
	 * @brief ����ĳ��λ�õ�ֵ 8Byte
	 */
	int32 Set8Byte(uint32 pos, uint64 p);

public:

	/**
	 * @brief ��ס������
	 */
	int32 EatBuf(const uint8 * buf, uint32 buf_sz);

	/**
	 * @brief get 1 byte
	 */
	int32 Get1Byte(uint8& p);

	/**
	 * @brief get 2 byte
	 */
	int32 Get2Byte(uint16& p);

	/**
	 * @brief get 4 byte
	 */
	int32 Get4Byte(uint32& p);

	/**
	 * @brief get 8 byte
	 */
	int32 Get8Byte(uint64& p);

	/**
	 * @brief get x byte
	 */
	int32 GetXByte(uint8 * p, uint32 sz);

	/**
	 * @brief get string
	 */
	int32 GetString(std::string& str);

	/**
	 * @brief ���û�����
	 */
	int32 ResetBuf();

private:

	/**
	 * @brief ���������
	 */
	mutable PKT_BUF vbuf_;

	/**
	 * @brief ����������
	 */
	uint8 * buf_;
	uint32 buf_sz_;
	/**
	 * @brief ����������������ǰ�Ѿ�δ�������Ļ�����ƫ��
	 */
	uint32 pos_;
};


#endif



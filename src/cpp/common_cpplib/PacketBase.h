/**
 * @file PacketBase.h 
 * @brief 私有协议报文基类
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
 * @brief 私有协议报文基类
 */
class PacketBase
{
public:
	
	/**
	 * @brief 缓冲区类型定义
	 */
	typedef std::vector<uint8> PKT_BUF;

	/**
	 * @brief 构造
	 */
	PacketBase();

	/**
	 * @brief 析构
	 */
	virtual ~PacketBase();

	/**
	 * @brief 获取缓冲区指定位置的地址
	 * @param pos:位置,0表示首地址
	 */
	const PKT_BUF& GetBufRef()const{return this->vbuf_;}

	/**
	 * @brief 获取缓冲区内容大小
	 */
	uint32 GetBufSZ();

	/**
	 * @brief 设置缓冲区的容量
	 */
	int32 SetBufCap(uint32 sz);

	/**
	 * @brief 获取缓冲区容量
	 */
	uint32 GetBufCap();

	/**
	 * @brief 获取位置
	 */
	uint32 GetPos(){return this->pos_;}

	/**
	 * @brief 获取剩余未解析的字节数
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
	 * @brief 设置某个位置的值 1Byte
	 */
	int32 Set1Byte(uint32 pos, uint8 p);

	/**
	 * @brief 设置某个位置的值 2Byte
	 */
	int32 Set2Byte(uint32 pos, uint16 p);

	/**
	 * @brief 设置某个位置的值 4Byte
	 */
	int32 Set4Byte(uint32 pos, uint32 p);

	/**
	 * @brief 设置某个位置的值 8Byte
	 */
	int32 Set8Byte(uint32 pos, uint64 p);

public:

	/**
	 * @brief 含住缓冲区
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
	 * @brief 重置缓冲区
	 */
	int32 ResetBuf();

private:

	/**
	 * @brief 组包缓冲区
	 */
	mutable PKT_BUF vbuf_;

	/**
	 * @brief 解析缓冲区
	 */
	uint8 * buf_;
	uint32 buf_sz_;
	/**
	 * @brief 解析缓冲区索引当前已经未被解析的缓冲区偏移
	 */
	uint32 pos_;
};


#endif



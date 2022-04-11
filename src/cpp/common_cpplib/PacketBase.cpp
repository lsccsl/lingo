/**
 * @file PacketBase.h 
 * @brief 私有协议报文基类
 *
 * @author linshaochuan
 */
#include "PacketBase.h"

#include <assert.h>

#include "mylogex.h"

/**
 * @brief 构造
 */
PacketBase::PacketBase():
	pos_(0),
	buf_(NULL),
	buf_sz_(0)
{}

/**
 * @brief 析构
 */
PacketBase::~PacketBase()
{}

/**
 * @brief 获取缓冲区指定位置的地址
 * @param pos:位置,0表示首地址
 */
//uint8 * PacketBase::GetBuf(uint32 pos)
//{
//	if(pos > this->vbuf_.capacity())
//		return NULL;
//
//	return &this->vbuf_[pos];
//}

/**
 * @brief 获取缓冲区内容大小
 */
uint32 PacketBase::GetBufSZ()
{
	return this->vbuf_.size();
}

/**
 * @brief 设置缓冲区的容量
 */
int32 PacketBase::SetBufCap(uint32 sz)
{
	this->vbuf_.resize(sz);

	return 0;
}

/**
 * @brief 获取缓冲区容量
 */
uint32 PacketBase::GetBufCap()
{
	return this->vbuf_.capacity();
}


/**
 * @brief add 1 byte
 */
int32 PacketBase::Add1Byte(uint8 p)
{
	this->vbuf_.push_back(p);
	return 0;
}

/**
 * @brief add 2 byte
 */
int32 PacketBase::Add2Byte(uint16 p)
{
	this->vbuf_.insert(this->vbuf_.end(), (uint8 *)&p, (uint8 *)(&p + 1));
	return 0;
}

/**
 * @brief add 4 byte
 */
int32 PacketBase::Add4Byte(uint32 p)
{
	this->vbuf_.insert(this->vbuf_.end(), (uint8 *)&p, (uint8 *)(&p + 1));
	return 0;
}

/**
 * @brief add 8 byte
 */
int32 PacketBase::Add8Byte(const uint64& p)
{
	this->vbuf_.insert(this->vbuf_.end(), (uint8 *)&p, (uint8 *)(&p + 1));
	return 0;
}

/**
 * @brief add x byte
 */
int32 PacketBase::AddXByte(uint8 * p, uint32 sz)
{
	this->vbuf_.insert(this->vbuf_.end(), (uint8 *)p, (uint8 *)((uint8 *)p + sz));
	return 0;
}

/**
 * @brief add string
 */
int32 PacketBase::AddString(const int8 * p)
{
	if(NULL == p)
		return -1;

	this->vbuf_.insert(this->vbuf_.end(), (uint8 *)p, (uint8 *)((uint8 *)p + strlen(p)));
	this->vbuf_.push_back('\0');

	return 0;
}


/**
 * @brief 设置某个位置的值 1Byte
 */
int32 PacketBase::Set1Byte(uint32 pos, uint8 p)
{
	if(pos + sizeof(uint8) > this->vbuf_.size())
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 1", __FILE__, __LINE__,
			this->vbuf_.size(), pos));
		return -1;
	}

	this->vbuf_[pos] = p;

	return 0;
}

/**
 * @brief 设置某个位置的值 2Byte
 */
int32 PacketBase::Set2Byte(uint32 pos, uint16 p)
{
	if(pos + sizeof(uint16) > this->vbuf_.size())
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 2", __FILE__, __LINE__,
			this->vbuf_.size(), pos));
		return -1;
	}

	memcpy(&this->vbuf_[pos], (uint8 *)&p, sizeof(p));

	return 0;
}

/**
 * @brief 设置某个位置的值 4Byte
 */
int32 PacketBase::Set4Byte(uint32 pos, uint32 p)
{
	if(pos + sizeof(uint32) > this->vbuf_.size())
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 4", __FILE__, __LINE__,
			this->vbuf_.size(), pos));
		return -1;
	}

	memcpy(&this->vbuf_[pos], (uint8 *)&p, sizeof(p));

	return 0;
}

/**
 * @brief 设置某个位置的值 8Byte
 */
int32 PacketBase::Set8Byte(uint32 pos, uint64 p)
{
	if(pos + sizeof(uint64) > this->vbuf_.size())
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 8", __FILE__, __LINE__,
			this->vbuf_.size(), pos));
		return -1;
	}

	memcpy(&this->vbuf_[pos], (uint8 *)&p, sizeof(p));

	return 0;
}


/**
 * @brief 含住缓冲区
 */
int32 PacketBase::EatBuf(const uint8 * buf, uint32 buf_sz)
{
	this->buf_ = (uint8 *)buf;
	this->buf_sz_ = buf_sz;
	this->pos_ = 0;

	return 0;
}

/**
 * @brief get 1 byte
 */
int32 PacketBase::Get1Byte(uint8& p)
{
	assert(this->buf_ && this->buf_sz_);

	if(this->pos_ + sizeof(uint8) > this->buf_sz_)
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 1", __FILE__, __LINE__,
			this->buf_sz_, this->pos_));
		return -1;
	}

	p = this->buf_[this->pos_];
	this->pos_ += sizeof(uint8);

	return 0;
}

/**
 * @brief get 2 byte
 */
int32 PacketBase::Get2Byte(uint16& p)
{
	assert(this->buf_ && this->buf_sz_);

	if(this->pos_ + sizeof(uint16) > this->buf_sz_)
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 2", __FILE__, __LINE__,
			this->buf_sz_, this->pos_));
		return -1;
	}

	p = *((uint16 *)&(this->buf_[this->pos_]));
	this->pos_ += sizeof(uint16);

	return 0;
}

/**
 * @brief get 4 byte
 */
int32 PacketBase::Get4Byte(uint32& p)
{
	assert(this->buf_ && this->buf_sz_);

	if(this->pos_ + sizeof(uint32) > this->buf_sz_)
	{
		MYLOG_INFO(("msg buf is not big enought size:%d pos:%d want 4",
			this->buf_sz_, this->pos_));
		return -1;
	}

	p = *((uint32 *)&(this->buf_[this->pos_]));
	this->pos_ += sizeof(uint32);

	return 0;
}

/**
 * @brief get 8 byte
 */
int32 PacketBase::Get8Byte(uint64& p)
{
	assert(this->buf_ && this->buf_sz_);

	if(this->pos_ + sizeof(uint64) > this->buf_sz_)
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want 4", __FILE__, __LINE__,
			this->buf_sz_, this->pos_));
		return -1;
	}

	p = *((uint64 *)&(this->buf_[this->pos_]));
	this->pos_ += sizeof(uint64);

	return 0;
}

/**
 * @brief get x byte
 */
int32 PacketBase::GetXByte(uint8 * p, uint32 sz)
{
	assert(this->buf_ && this->buf_sz_);

	if(NULL == p || 0 == sz)
		return -1;

	if((this->pos_ + sz) > this->buf_sz_)
	{
		MYLOG_INFO(("msg buf is not big enought size:%d pos:%d want %d",
			this->buf_sz_, this->pos_, sz));
		return -1;
	}
	
	memcpy(p, &(this->buf_[this->pos_]), sz);
	this->pos_ += sz;

	return 0;
}

/**
 * @brief get string
 */
int32 PacketBase::GetString(std::string& str)
{
	assert(this->buf_ && this->buf_sz_);

	if(this->pos_ >= this->buf_sz_)
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d want string", __FILE__, __LINE__,
			this->buf_sz_, this->pos_));
		return -1;
	}

	str = (char *)&this->buf_[this->pos_];
	if((this->pos_ + str.size()) > this->buf_sz_)
	{
		MYLOG_INFO(("[%s:%d]msg buf is not big enought size:%d pos:%d string len %d", __FILE__, __LINE__,
			this->buf_sz_, this->pos_, str.size()));
	}

	this->pos_ += str.size() + 1;

	return 0;
}


/**
 * @brief 重置缓冲区
 */
int32 PacketBase::ResetBuf()
{
	this->vbuf_.clear();
	this->pos_ = 0;
	this->buf_ = NULL;
	this->buf_sz_ = 0;
	return 0;
}

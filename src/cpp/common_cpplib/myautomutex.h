/**
* @file myautomutex.h
* @brief ×Ô¶¯Ëø
*
* @author linshaochuan
* @blog http://blog.csdn.net/lsccsl
*/
#ifndef __MYAUTOMUTEX_H__
#define __MYAUTOMUTEX_H__

extern "C"
{
	#include "mymutex.h"
}

class myautomutex
{
public:

	/**
	* @brief constructor
	*/
	myautomutex(HMYMUTEX& lock);

	/**
	* @brief destructor
	*/
	~myautomutex();

private:

	HMYMUTEX& lock_;
};

#endif

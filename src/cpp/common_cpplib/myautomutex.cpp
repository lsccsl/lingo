/**
* @file myautomutex.cpp
* @brief ×Ô¶¯Ëø
*
* @author linshaochuan
*/
#include "myautomutex.h"
#include "mylogex.h"

/**
* @brief constructor
*/
myautomutex::myautomutex(HMYMUTEX& lock):lock_(lock)
{
	MYLOG_DEBUG(("myautomutex::myautomutex"));
	MyMutexLock(lock_);
}

/**
* @brief destructor
*/
myautomutex::~myautomutex()
{
	MYLOG_DEBUG(("myautomutex::~myautomutex"));
	MyMutexUnLock(lock_);
}

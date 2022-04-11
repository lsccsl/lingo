
#ifndef __OSINET_H__
#define __OSINET_H__


#include <string>
#include <vector>
#include "type_def.h"


class OsInet
{
public:

	/**
	 * @brief 获取本机ip
	 */
	static int32 GetLocalIP(uint32& ip);
	static void GetLocalIPs(std::vector<uint32>& ip_list);
};


#endif


/**
 * @file os.cpp
 */
#include "myos_cpp.h"

#include <assert.h>
#include <time.h>
#ifdef WIN32
#include <windows.h>
#define sleep(x) Sleep(x * 1000)
#else
#include <unistd.h>
#include <dirent.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <string.h>
#endif
#include <pthread.h>

/**
 * @brief sleep for sec second
 */
void myos::ossleep(uint32 sec)
{
	sleep(sec);
}

/**
 * @brief 
 */
int32 myos::get_all_sub_file(const std::string& parent_path, std::set<std::string>& sfile)
{
	struct dirent ** namelist = NULL;
	int32 n;

#ifndef WIN32
	n = scandir(parent_path.c_str(), &namelist, 0, alphasort);
	if (n < 0)
	{
		perror("scandir");
		return -1;
	}
	else
	{
		while(n--)
		{
			struct stat buf;
			std::string full_path = parent_path;
			full_path.append("/");
			full_path.append(namelist[n]->d_name);
			stat(full_path.c_str(), &buf);
			if(!S_ISDIR(buf.st_mode))
			{
				sfile.insert(namelist[n]->d_name);
			}

			free(namelist[n]);
		}
		free(namelist);
	}
#endif

	return 0;
}
int32 myos::get_all_sub_folder(const std::string& parent_path, std::set<std::string>& sfolder)
{
	struct dirent ** namelist = NULL;
	int32 n;

#ifndef WIN32
	n = scandir(parent_path.c_str(), &namelist, 0, alphasort);
	if (n < 0)
	{
		perror("scandir");
		return -1;
	}
	else
	{
		while(n--)
		{
			if(strncmp(namelist[n]->d_name, ".", 1) == 0)
				continue;
			if(strncmp(namelist[n]->d_name, "..", 2) == 0)
				continue;
			struct stat buf;
			std::string full_path = parent_path;
			full_path.append("/");
			full_path.append(namelist[n]->d_name);
			stat(full_path.c_str(), &buf);
			if(S_ISDIR(buf.st_mode))
			{
				sfolder.insert(namelist[n]->d_name);
			}

			free(namelist[n]);
		}
		free(namelist);
	}
#endif

	return 0;
}

/**
 * @brief 字符串转成整数
 */
void myos::string2uint64(const std::string& s, uint64& n)
{
//#ifndef _MBCS
	std::istrstream iss(s.c_str());
	iss >> n;
//#endif
}
void myos::string2uint32(const std::string& s, uint32& n)
{
	std::istrstream iss(s.c_str());
	iss >> n;
}

/**
 * @brief 整形转成字符串
 */
void myos::uint642string(const uint64 n, std::string& s)
{
	std::ostringstream oss;
	oss << n;
	s = oss.str();
}
void myos::uint322string(const uint32 n, std::string& s)
{
	std::ostringstream oss;
	oss << n;
	s = oss.str();
}

/**
 * @brief ip转正字符串
 */
void myos::ip2string(const uint32 int_ip, std::string& ip)
{
	ip.resize(32);
#ifdef WIN32
	/* why no snprintf for win */
#pragma   warning(   disable   :   4996) /* fuck vc,why warning */ 
	int32 len = sprintf((int8 *)ip.c_str(), "%d.%d.%d.%d",
		int_ip & 0xff,
		(int_ip >> 8) & 0xff,
		(int_ip >> 16) & 0xff,
		(int_ip >> 24) & 0xff);
#else
	int32 len = snprintf((int8 *)ip.c_str(), ip.size(), "%d.%d.%d.%d",
		int_ip & 0xff,
		(int_ip >> 8) & 0xff,
		(int_ip >> 16) & 0xff,
		(int_ip >> 24) & 0xff);
#endif
	ip.resize(len);
}
void myos::string2ip(uint32& uip, const std::string& sip)
{
	uip = inet_addr(sip.c_str());
}

/**
* @brief 将秒数表征的时间转成字符串 如:2011-01-01 01:01:01
*/
void myos::conver_time2string(const time_t time_sec, std::string& output)
{
	struct tm ptm;
	localtime_r(&time_sec, &ptm);
	char actemp[128] = {0};

	sprintf(actemp, "%d-%d-%d %d:%d:%d",
		ptm.tm_year + 1900,
		ptm.tm_mon + 1,
		ptm.tm_mday,
		ptm.tm_hour,
		ptm.tm_min,
		ptm.tm_sec);

	output = actemp;
}
/**
* @brief 解析时间格 2011-01-01 01:01:01
*/
const time_t myos::conver_string2time(const int8 * time_string)
{
	enum{
		PARSE_DATETIME_YEAR,
		PARSE_DATETIME_MONTH,
		PARSE_DATETIME_DAY,
		PARSE_DATETIME_HOUR,
		PARSE_DATETIME_MINUTE,
		PARSE_DATETIME_SEC,
	};

	struct parse_datetime_state_t
	{
		parse_datetime_state_t():state_(PARSE_DATETIME_YEAR){}
		int32 state_;
		int8 * cur_begin_;
		int8 * cur_end_;
	}parse_state;

	struct tm tm_temp = {0};

	if(NULL == time_string)
		return 0;

	int8 actemp[32] = {0};
	strncpy(actemp, time_string, sizeof(actemp) - 1);
	int8 * p = (int8 *)actemp;
	parse_state.cur_begin_ = NULL;
	parse_state.cur_end_ = NULL;
	for(; *p; p ++)
	{
		if(::isdigit(*p))
		{
			if(!parse_state.cur_begin_)
				parse_state.cur_begin_ = p;
			if(0 == *(p+1))
			{
				parse_state.cur_end_ = p + 1;
			}
			else
			{
				continue;
			}
		}
		else
		{
			parse_state.cur_end_ = p;
		}

		/* 根据逻辑,此条件一定成立 */
		assert(parse_state.cur_end_);

		int8 temp = *(parse_state.cur_end_);
		*parse_state.cur_end_ = 0;

		switch(parse_state.state_)
		{
		case PARSE_DATETIME_YEAR:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_year = atoi(parse_state.cur_begin_) - 1900;

				parse_state.state_ = PARSE_DATETIME_MONTH;
			}
			break;

		case PARSE_DATETIME_MONTH:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_mon = atoi(parse_state.cur_begin_) - 1;

				parse_state.state_ = PARSE_DATETIME_DAY;
			}
			break;

		case PARSE_DATETIME_DAY:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_mday = atoi(parse_state.cur_begin_);

				parse_state.state_ = PARSE_DATETIME_HOUR;
			}
			break;

		case PARSE_DATETIME_HOUR:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_hour = atoi(parse_state.cur_begin_);

				parse_state.state_ = PARSE_DATETIME_MINUTE;
			}
			break;

		case PARSE_DATETIME_MINUTE:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_min = atoi(parse_state.cur_begin_);

				parse_state.state_ = PARSE_DATETIME_SEC;
			}
			break;

		case PARSE_DATETIME_SEC:
			{
				if(parse_state.cur_begin_)
					tm_temp.tm_sec = atoi(parse_state.cur_begin_);

				return mktime(&tm_temp);
			}
			break;

		default:
			break;
		}

		*(parse_state.cur_end_) = temp;
		parse_state.cur_begin_ = NULL;
		parse_state.cur_end_ = NULL;
	}

	return mktime(&tm_temp);
}

/**
* @brief 获取cpu使用率,返回值单位是1%
*/
const uint32 myos::get_cpu()
{
#ifndef WIN32
	char actemp[256] = {0};

	FILE * file = popen("cat /proc/stat|grep cpu", "r");
	if(NULL == file)
		return 0;
	fgets(actemp, sizeof(actemp) - 1, file);
	fclose(file);

	int total;
	int user;
	int nice;
	int system;
	int idle;
	char cpu[32];

	sscanf(actemp, "%s %d %d %d %d", cpu, &user, &nice, &system, &idle);

	total = (user + nice + system + idle);

	return ((user + nice + system) * 100) / total;

#else
	return 0;
#endif
}

/**
* @brief 获取内存使用率,返回值单位是1%
*/
const uint32 myos::get_mem()
{
#ifndef WIN32

	char actemp[256] = {0};

	char temp1[32];
	char temp2[32];

	int MemTotal = 0;
	int MemFree = 0;
	int Buffers = 0;
	int Cached = 0;

	FILE * file = popen("cat /proc/meminfo", "r");
	if(NULL == file)
		return 0;

	fgets(actemp, sizeof(actemp) - 1, file);
	sscanf (actemp, "%s %u %s", temp1, &MemTotal, temp2);

	fgets(actemp, sizeof(actemp) - 1, file);
	sscanf (actemp, "%s %u %s", temp1, &MemFree, temp2);

	fgets(actemp, sizeof(actemp) - 1, file);
	sscanf (actemp, "%s %u %s", temp1, &Buffers, temp2);

	fgets(actemp, sizeof(actemp) - 1, file);
	sscanf (actemp, "%s %u %s", temp1, &Cached, temp2);

	fclose(file);

	return ((MemTotal - MemFree - Buffers - Cached) * 100) / MemTotal;

#else
	return 0;
#endif
}

/* @brief 分割pcString 用pcToken */
void myos::StringOpSplitString(const int8 * pcString, 
	const int8 * pcToken, 
	std::vector<std::string>& vOut)
{
	if(NULL == pcString || NULL == pcToken)
		return;

	int32 nLen = strlen(pcString);
	std::string strTemp;
	int nLenTok = strlen(pcToken);
	for(int32 i = 0; i < nLen; i ++)
	{
		//int nLenTok = strlen(pcToken);
		int isTok = 0;
		for(int32 j = 0; j < nLenTok; j ++)
		{
			if(pcToken[j] != pcString[i])
				continue;

			/* 遇到分割符,存入vOut,重新收集 */
			if(strTemp.size())
				vOut.push_back(strTemp);
#ifdef _MBCS
			strTemp = "";
#else
			strTemp.clear();
#endif
			isTok = 1;
			break;
		}
		/* 收集字符串 */
		if(!isTok)
		{
#ifdef _MBCS
			strTemp.append(1, (char)(pcString[i]));
#else
			strTemp.push_back(pcString[i]);
#endif
		}
	}

	/* 还剩一串 */
	if(strTemp.size())
		vOut.push_back(strTemp);
}

/**
* @brief 将含有通配符的语句转化为mysql查询语句
* @param strQuery 输出字符串
* @param pcReg 含有通配符的列目录语句
*/
void myos::StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg)
{
	int iLen = strlen(pcReg);
	for(int i = 0 ; i < iLen; i++)
	{
/*		if( '*' == pcReg[i] || '<' == pcReg[i] ) strQuery.append("%");
		else if( '?' == pcReg[i] || '>' == pcReg[i] ) strQuery.append("_");
		else if( '\"' == pcReg[i] ) strQuery.append(".");
		else */
		if( '%' == pcReg[i] ) strQuery.append("\\%");
		else if( '_' == pcReg[i] ) strQuery.append("\\_");
		else if( '\'' == pcReg[i] ) strQuery.append("\\\'");
		else strQuery.append(pcReg + i, 1);
	}
	return ;
}

















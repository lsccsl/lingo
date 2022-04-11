/**
 * @file os.h
 */
#ifndef __MYOS_H__
#define __MYOS_H__

#include <string>
#include <set>
#include <strstream>
#include <sstream>
#include <vector>

#include "type_def.h"

class myos
{
public:

	/**
	 * @brief sleep for sec second
	 */
	static void ossleep(uint32 sec);

	/**
	 * @brief 获取目录下面的所有文件/目录
	 */
	static int32 get_all_sub_file(const std::string& parent_path, std::set<std::string>& sfile);
	static int32 get_all_sub_folder(const std::string& parent_path, std::set<std::string>& sfolder);

	/**
	 * @brief 字符串转成整数
	 */
	static void string2uint64(const std::string& s, uint64& n);
	static void string2uint32(const std::string& s, uint32& n);

	/**
	 * @brief 整形转成字符串
	 */
	static void uint642string(const uint64 n, std::string& s);
	static void uint322string(const uint32 n, std::string& s);

	/**
	 * @brief ip转正字符串
	 */
	static void ip2string(const uint32 uip, std::string& sip);
	static void string2ip(uint32& uip, const std::string& sip);

	/**
	* @brief 将秒数表征的时间转成字符串 如:2011-01-01 01:01:01
	*/
	static void conver_time2string(const time_t time_sec, std::string& output);
	/**
	* @brief 解析时间格 2011-01-01 01:01:01
	*/
	static const time_t conver_string2time(const int8 * time_string);

	/**
	* @brief 获取cpu使用率,返回值单位是1%
	*/
	static const uint32 get_cpu();

	/**
	* @brief 获取内存使用率,返回值单位是1%
	*/
	static const uint32 get_mem();

	/* @brief 分割pcString 用pcToken */
	static void StringOpSplitString(const char * pcString, 
		const char * pcToken, 
		std::vector<std::string>& vOut);

	/**
	* @brief 将含有通配符的语句转化为mysql查询语句
	* @param strQuery 输出字符串
	* @param pcReg 含有通配符的列目录语句
	*/
	static void StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg);
};

#endif

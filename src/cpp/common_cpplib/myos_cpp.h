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
	 * @brief ��ȡĿ¼����������ļ�/Ŀ¼
	 */
	static int32 get_all_sub_file(const std::string& parent_path, std::set<std::string>& sfile);
	static int32 get_all_sub_folder(const std::string& parent_path, std::set<std::string>& sfolder);

	/**
	 * @brief �ַ���ת������
	 */
	static void string2uint64(const std::string& s, uint64& n);
	static void string2uint32(const std::string& s, uint32& n);

	/**
	 * @brief ����ת���ַ���
	 */
	static void uint642string(const uint64 n, std::string& s);
	static void uint322string(const uint32 n, std::string& s);

	/**
	 * @brief ipת���ַ���
	 */
	static void ip2string(const uint32 uip, std::string& sip);
	static void string2ip(uint32& uip, const std::string& sip);

	/**
	* @brief ������������ʱ��ת���ַ��� ��:2011-01-01 01:01:01
	*/
	static void conver_time2string(const time_t time_sec, std::string& output);
	/**
	* @brief ����ʱ��� 2011-01-01 01:01:01
	*/
	static const time_t conver_string2time(const int8 * time_string);

	/**
	* @brief ��ȡcpuʹ����,����ֵ��λ��1%
	*/
	static const uint32 get_cpu();

	/**
	* @brief ��ȡ�ڴ�ʹ����,����ֵ��λ��1%
	*/
	static const uint32 get_mem();

	/* @brief �ָ�pcString ��pcToken */
	static void StringOpSplitString(const char * pcString, 
		const char * pcToken, 
		std::vector<std::string>& vOut);

	/**
	* @brief ������ͨ��������ת��Ϊmysql��ѯ���
	* @param strQuery ����ַ���
	* @param pcReg ����ͨ�������Ŀ¼���
	*/
	static void StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg);
};

#endif

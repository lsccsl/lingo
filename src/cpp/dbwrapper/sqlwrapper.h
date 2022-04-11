/**
* @file sqlwrapper.h
* @brief db编程接口封装
* @author linsc
*/
#ifndef __SQLWRAPPER_H__
#define __SQLWRAPPER_H__

#include <string>
#include <vector>
#include "type_def.h"

class sqlwrapper
{
public:

	/**
	 * @brief constructor/destructor
	 */
	sqlwrapper(){}
	virtual ~sqlwrapper(){}

	/**
	* @brief 初始化
	*/
	virtual int init() = 0;

	/**
	* @brief 反初始化
	*/
	virtual void uninit() = 0;

	/**
	* @brief 执行sql语句,0:成功 -1:失败
	*/
	virtual int32 query(const int8 * pcsql) = 0;

	/**
	* @brief 获取结果集
	*/
	virtual int32 get_result(std::vector<std::vector<std::string> >& vresult) = 0;

	/**
	* @brief 获取db时间 
	*/
	virtual int32 get_dbtime(uint64& time_db){ return -1; }

	/**
	* @brief 获取db版本信息
	*/
	virtual const int8 * get_dbver(){ return "unknown db ver"; }

	/**
	* @brief 获取字符集
	*/
	virtual const int8 * get_charset(){ return "null"; }

	/**
	* @brief 设置字符集
	*/
	virtual int32 set_charset(const int8 * charset){ return -1; };

	/**
	* @brief 查询表是否存在 0:成功 其它:失败
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist){ return -1; }

	/**
	* @brief 复制表结构
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst){ return -1; }

	/**
	* @brief 获取最后一个插入的自增主键
	*/
	virtual int32 get_last_insert_id(uint64& last_id){ return -1; }

	/**
	* @brief 获取最后一次查询受影响的行数
	*/
	virtual int32 get_affected_rows(uint32& affected_rows){ return -1; }

	/**
	* @brief 获取/设置用户数据
	*/
	virtual void * get_userdata(){ return NULL; }
	virtual void set_user_data(void * data){}

	/**
	* @brief view for debug
	*/
	virtual void view(){};
};

#endif





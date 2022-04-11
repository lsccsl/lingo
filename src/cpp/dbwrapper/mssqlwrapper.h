/**
* @file mssqlwrapper.h
* @brief ms的db编程接口封装
* @author linsc
*/
#ifndef __MSSQLWRAPPER_H__
#define __MSSQLWRAPPER_H__

#include "sqlwrapper.h"
#include "type_def.h"
#include <vector>

struct mssql_handle_inter;

class mssqlwrapper : public sqlwrapper
{
public:

	/**
	* @brief constructor
	*/
	mssqlwrapper(const int8 * ip, uint32 port,
		const int8 * user, const int8 * pwd,
		const int8 * dbname);

	/**
	* @brief constructor
	*/
	virtual ~mssqlwrapper();

	/**
	* @brief 初始化
	*/
	virtual int init();

	/**
	* @brief 反初始化
	*/
	virtual void uninit();

	/**
	* @brief 执行sql语句,0:成功 -1:失败
	*/
	virtual int32 query(const int8 * pcsql);

	/**
	* @brief 获取结果集
	*/
	virtual int32 get_result(std::vector<std::vector<std::string> >& vresult);

	/**
	* @brief 获取db版本信息
	*/
	virtual const int8 * get_dbver();

	/**
	* @brief 查询表是否存在
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist);

	/**
	* @brief 复制表结构
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst);

	/**
	* @brief 获取最后一个插入的自增主键
	*/
	virtual int32 get_last_insert_id(uint64& last_id);

private:

	/**
	 * @brief 数据库连接句柄
	 */
	mssql_handle_inter * h_;

	/**
	* @brief 数据库连接信息
	*/
	std::string db_ip_;
	uint32 db_port_;
	std::string db_user_;
	std::string db_pwd_;
	std::string db_name_;
};

#endif









/**
* @file mysqlwrapper.h
* @brief 
* @author linsc
*/
#ifndef __MYSQLWRAPPER_H__
#define __MYSQLWRAPPER_H__

#include "sqlwrapper.h"

class MysqlWrapper : public sqlwrapper
{
public:

	/**
	 * @brief constructor
	 */
	MysqlWrapper(const int8 * ip, uint32 port,
		const int8 * user, const int8 * pwd,
		const int8 * dbname);

	/**
	* @brief constructor
	*/
	virtual ~MysqlWrapper();

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
	* @brief 查询表是否存在
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist);

	/**
	* @brief 复制表结构
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst);

	/**
	* @brief 获取最后播入的id
	*/
	virtual int32 get_last_insert_id(uint64& last_id);

	/**
	* @brief 获取最后一次查询受影响的行数
	*/
	virtual int32 get_affected_rows(uint32& affected_rows);

	/**
	* @brief 获取db时间 
	*/
	virtual int32 get_dbtime(uint64& time_db);

private:

	/**
	 * @brief 数据库句柄
	 */
	void * hdb_;

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











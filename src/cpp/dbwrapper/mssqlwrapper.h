/**
* @file mssqlwrapper.h
* @brief ms��db��̽ӿڷ�װ
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
	* @brief ��ʼ��
	*/
	virtual int init();

	/**
	* @brief ����ʼ��
	*/
	virtual void uninit();

	/**
	* @brief ִ��sql���,0:�ɹ� -1:ʧ��
	*/
	virtual int32 query(const int8 * pcsql);

	/**
	* @brief ��ȡ�����
	*/
	virtual int32 get_result(std::vector<std::vector<std::string> >& vresult);

	/**
	* @brief ��ȡdb�汾��Ϣ
	*/
	virtual const int8 * get_dbver();

	/**
	* @brief ��ѯ���Ƿ����
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist);

	/**
	* @brief ���Ʊ�ṹ
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst);

	/**
	* @brief ��ȡ���һ���������������
	*/
	virtual int32 get_last_insert_id(uint64& last_id);

private:

	/**
	 * @brief ���ݿ����Ӿ��
	 */
	mssql_handle_inter * h_;

	/**
	* @brief ���ݿ�������Ϣ
	*/
	std::string db_ip_;
	uint32 db_port_;
	std::string db_user_;
	std::string db_pwd_;
	std::string db_name_;
};

#endif









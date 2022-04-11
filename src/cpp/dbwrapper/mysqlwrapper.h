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
	* @brief ��ѯ���Ƿ����
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist);

	/**
	* @brief ���Ʊ�ṹ
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst);

	/**
	* @brief ��ȡ������id
	*/
	virtual int32 get_last_insert_id(uint64& last_id);

	/**
	* @brief ��ȡ���һ�β�ѯ��Ӱ�������
	*/
	virtual int32 get_affected_rows(uint32& affected_rows);

	/**
	* @brief ��ȡdbʱ�� 
	*/
	virtual int32 get_dbtime(uint64& time_db);

private:

	/**
	 * @brief ���ݿ���
	 */
	void * hdb_;

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











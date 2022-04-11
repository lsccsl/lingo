/**
* @file sqlwrapper.h
* @brief db��̽ӿڷ�װ
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
	* @brief ��ʼ��
	*/
	virtual int init() = 0;

	/**
	* @brief ����ʼ��
	*/
	virtual void uninit() = 0;

	/**
	* @brief ִ��sql���,0:�ɹ� -1:ʧ��
	*/
	virtual int32 query(const int8 * pcsql) = 0;

	/**
	* @brief ��ȡ�����
	*/
	virtual int32 get_result(std::vector<std::vector<std::string> >& vresult) = 0;

	/**
	* @brief ��ȡdbʱ�� 
	*/
	virtual int32 get_dbtime(uint64& time_db){ return -1; }

	/**
	* @brief ��ȡdb�汾��Ϣ
	*/
	virtual const int8 * get_dbver(){ return "unknown db ver"; }

	/**
	* @brief ��ȡ�ַ���
	*/
	virtual const int8 * get_charset(){ return "null"; }

	/**
	* @brief �����ַ���
	*/
	virtual int32 set_charset(const int8 * charset){ return -1; };

	/**
	* @brief ��ѯ���Ƿ���� 0:�ɹ� ����:ʧ��
	*/
	virtual int32 table_exist(const int8 * table_name, int32& bexist){ return -1; }

	/**
	* @brief ���Ʊ�ṹ
	*/
	virtual int32 copy_table_struct(const int8 * table_src, const int8 * table_dst){ return -1; }

	/**
	* @brief ��ȡ���һ���������������
	*/
	virtual int32 get_last_insert_id(uint64& last_id){ return -1; }

	/**
	* @brief ��ȡ���һ�β�ѯ��Ӱ�������
	*/
	virtual int32 get_affected_rows(uint32& affected_rows){ return -1; }

	/**
	* @brief ��ȡ/�����û�����
	*/
	virtual void * get_userdata(){ return NULL; }
	virtual void set_user_data(void * data){}

	/**
	* @brief view for debug
	*/
	virtual void view(){};
};

#endif





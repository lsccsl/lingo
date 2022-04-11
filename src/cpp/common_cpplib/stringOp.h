/*!@file
********************************************************************************
<pre>
ģ������       
�ļ�����        stringOp.h
����ļ���      
�ļ�ʵ�ֹ��ܣ�  string operation
���ߣ�          linsc
�汾��          1.0
--------------------------------------------------------------------------------
��ע:
--------------------------------------------------------------------------------
�޸ļ�¼ : 
�� ��        �汾     �޸���            �޸�����
2009/7/8     1.0      linsc             create
</pre>
*******************************************************************************/
#include <vector>
#include <string>

#pragma   warning(   disable   :   4786)

/*!
@brief unicode to utf8
@param wstrIn �����unicode string�����
@param wLen wstrIn�Ĵ�С
@param utf8_string �����utf8 string
@return 
********************************************************************/
extern void StringOpUnicodeToUTF8(const unsigned short *wstrIn,
	int wLen, 
	std::vector<unsigned char>& utf8_string);

/*!
@brief �ָ�string ��pcToken
@param pcString Ҫ�ָ���ַ���
@param pcToken �ָ��
@param vOut �ָ������
@return 
********************************************************************/
extern void StringOpSplitString(const char * pcString, 
	const char * pcToken, 
	std::vector<std::string>& vOut);

/**
 * @brief �ָ��ַ���
 */
extern void StringOpSplitLast(const char * pcString,
	std::string& parent_path,
	std::string& name,
	const char * pcToken);

/* ���������******************************************************************/
/*!
@brief ������ͨ��������ת��Ϊmysql��ѯ���
@param strQuery ����ַ���
@param pcReg ����ͨ�������Ŀ¼���
@return 
*/
extern void StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg);
/****************************************************************************/






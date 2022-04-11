/*!@file
********************************************************************************
<pre>
ģ������       
�ļ�����        stringOp.cpp
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
#include "stringOp.h"
#include "type_def.h"
#include <string.h>


/*!
@brief unicode to utf8
@param wstrIn �����unicode string�����
@param wLen wstrIn�Ĵ�С
@param utf8_string �����utf8 string
@return 
********************************************************************/
void ConVertUnicodeToUTF8(const uint8 *wstrIn,uint32 wLen, std::vector<uint8>& utf8_string)
{
	if(NULL == wstrIn || 0 == wLen)
		return;

#define putchar(a) utf8_string.push_back(a);

	for(uint32 j=0;(uint32)j<wLen;j++)
	{
		uint16 c=wstrIn[j];
		if (c < 0x80)
		{
			putchar (c);
		}
		else if (c < 0x800)
		{
			putchar (0xC0 | c>>6);
			putchar (0x80 | c & 0x3F);
		}
		else if (c < 0x10000)
		{
			putchar (0xE0 | c>>12);
			putchar (0x80 | c>>6 & 0x3F);
			putchar (0x80 | c & 0x3F);
		}
		else if (c < 0x200000)
		{
			putchar (0xF0 | c>>18);
			putchar (0x80 | c>>12 & 0x3F);
			putchar (0x80 | c>>6 & 0x3F);
			putchar (0x80 | c & 0x3F);
		}
	}
#undef putchar
}

/*!
@brief �ָ�string ��pcToken
@param pcString Ҫ�ָ���ַ���
@param pcToken �ָ��
@param vOut �ָ������
@return 
********************************************************************/
void StringOpSplitString(const int8 * pcString, 
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

			/* �����ָ��,����vOut,�����ռ� */
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
		/* �ռ��ַ��� */
		if(!isTok)
		{
#ifdef _MBCS
			strTemp.append(1, (char)(pcString[i]));
#else
			strTemp.push_back(pcString[i]);
#endif
		}
	}

	/* ��ʣһ�� */
	if(strTemp.size())
		vOut.push_back(strTemp);
}

/**
 * @brief �ָ��ַ���
 */
extern void StringOpSplitLast(const int8 * pcString,
	std::string& parent_path,
	std::string& name,
	const int8 * pcToken)
{
	if(NULL == pcString || NULL == pcToken)
		return;

	int32 pos = strlen(pcString) - 1;
	int nLenTok = strlen(pcToken);

	uint32 name_end = 0;
	uint32 name_begin = 0;
	int begin_eat_name = 0;

	for(; pos >= 0; pos --)
	{
		int isTok = 0;
		for(int32 j = 0; j < nLenTok; j ++)
		{
			if(pcString[pos] != pcToken[j])
				continue;

			isTok = 1;

			/*  */
			if(begin_eat_name)
			{
				name_begin = pos + 1;
				break;
			}
		}

		if(isTok)
		{
			if(begin_eat_name)
				break;
			else
				continue;
		}

		if(0 == begin_eat_name)
		{
			name_end = pos;
			begin_eat_name = 1;
		}
	}

	name.insert(0, &pcString[name_begin], name_end - name_begin + 1);
	parent_path.insert(0, pcString, name_begin);
}

/* ���������******************************************************************/
/*!
@brief ������ͨ��������ת��Ϊmysql��ѯ���
@param strQuery ����ַ���
@param pcReg ����ͨ�������Ŀ¼���
@return 
*/
extern void StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg)
{
	int iLen = strlen(pcReg);
	for(int i = 0 ; i < iLen; i++)
	{
		if( '*' == pcReg[i] || '<' == pcReg[i] ) strQuery.append("%");
		else if( '?' == pcReg[i] || '>' == pcReg[i] ) strQuery.append("_");
		else if( '\"' == pcReg[i] ) strQuery.append(".");
		else if( '%' == pcReg[i] ) strQuery.append("\\%");
		else if( '_' == pcReg[i] ) strQuery.append("\\_");
		else if( '\'' == pcReg[i] ) strQuery.append("\\\'");
		else strQuery.append(pcReg + i, 1);
	}
	return ;
}
/****************************************************************************/











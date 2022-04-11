/*!@file
********************************************************************************
<pre>
模块名：       
文件名：        stringOp.h
相关文件：      
文件实现功能：  string operation
作者：          linsc
版本：          1.0
--------------------------------------------------------------------------------
备注:
--------------------------------------------------------------------------------
修改记录 : 
日 期        版本     修改人            修改内容
2009/7/8     1.0      linsc             create
</pre>
*******************************************************************************/
#include <vector>
#include <string>

#pragma   warning(   disable   :   4786)

/*!
@brief unicode to utf8
@param wstrIn 输入的unicode string缓存冲
@param wLen wstrIn的大小
@param utf8_string 输出的utf8 string
@return 
********************************************************************/
extern void StringOpUnicodeToUTF8(const unsigned short *wstrIn,
	int wLen, 
	std::vector<unsigned char>& utf8_string);

/*!
@brief 分割string 用pcToken
@param pcString 要分割的字符串
@param pcToken 分割符
@param vOut 分割后的输出
@return 
********************************************************************/
extern void StringOpSplitString(const char * pcString, 
	const char * pcToken, 
	std::vector<std::string>& vOut);

/**
 * @brief 分割字符串
 */
extern void StringOpSplitLast(const char * pcString,
	std::string& parent_path,
	std::string& name,
	const char * pcToken);

/* 刘冀鹏添加******************************************************************/
/*!
@brief 将含有通配符的语句转化为mysql查询语句
@param strQuery 输出字符串
@param pcReg 含有通配符的列目录语句
@return 
*/
extern void StringOpConvertToMysqlQuery(std::string& strQuery, const char* pcReg);
/****************************************************************************/







/**
 * @file myOsFile.h 封装不同系统的文件操作接口 2008-1-30 00:43
 *
 * @author lin shao chuan (email:lsccsl@tom.com, msn:lsccsl@163.net)
 * @blog http://blog.csdn.net/lsccsl
 *
 * @brief if it works, it was written by lin shao chuan, if not, i don't know who wrote it.
 *        封装不同系统的文件操作接口
 *
 * Permission to use, copy, modify, distribute and sell this software
 * and its documentation for any purpose is hereby granted without fee,
 * provided that the above copyright notice appear in all copies and
 * that both that copyright notice and this permission notice appear
 * in supporting documentation.  lin shao chuan makes no
 * representations about the suitability of this software for any
 * purpose.  It is provided "as is" without express or implied warranty.
 * see the GNU General Public License  for more detail.
 */
#ifndef __OS_FILE_HACB__
#define __OS_FILE_HACB__

#include <stdlib.h>
#include <vector>

#include "type_def.h"


struct __os_file_t_;
typedef struct __os_file_t_ * HMYOSFILE;

/**
 * @brief 打开文件,读写
 */
extern HMYOSFILE myOsFileOpenReadWriteEx(const char * file_name);

/**
 * @brief 以只读的方式打开文件
 */
extern HMYOSFILE myOsFileOpenReadOnlyEx(const char * file_name);

/**
 * @brief 以独占的方式打开文件
 */
extern HMYOSFILE myOsFileOpenExclusiveEx(const char * file_name);

/**
 * @brief 拷贝对象
 */
extern HMYOSFILE myOSFileOpenEx(HMYOSFILE hf);

/**
 * @brief 关闭文件
 */
extern int myOsFileCloseEx(HMYOSFILE hf);

/**
 * @brief 同步文件至辅存
 * @return 0:成功 -1:失败
 */
extern int myOsFileSynEx(HMYOSFILE hf);

/**
 * @brief 写文件
 * @return 0:成功 -1:失败, -2:未写满指定字节
 */
extern int myOsFileWriteEx(HMYOSFILE hf, const void * data, size_t data_size, size_t * write_size);

/**
 * @brief 读文件
 * @return 0:成功 -1:失败
 */
extern int myOsFileReadEx(HMYOSFILE hf, void * data, size_t data_size, size_t * read_size);

/**
 * @brief 移动当前的文件指针至off_set(相对于文件头)
 * @return 0:成功 -1:失败
 */
extern int myOsFileSeekEx(HMYOSFILE hf, int64 off_set);

/**
 * @brief 删除文件
 * @return 0:成功 -1:失败
 */
extern int myOsFileDelEx(const char * file_name);

/**
 * @brief 获取文件的大小
 * @return 0:成功 -1:失败
 */
extern int myOsFileSizeEx(HMYOSFILE hf, int64 * file_size);

/**
 * @brief 获取文件句柄
 */
extern int myOsFileGetFdEx(HMYOSFILE hf);

/**
 * @brief 判断文件是否存在
 * @return 0:文件不存在 非零:文件存在
 */
extern int myOsFileExistsEx(const char * file_name);

/**
 * @brief 栽减文件
 * @return 0:成功, -1:失败
 */
extern int myOsFileTruncateEx(HMYOSFILE hf, int64 nByte);

/**
 * @brief 重命名文件
 */
extern int myOsRenameEx(const char * old_name, const char * new_name);

/**
 * @brief unlink
 */
extern int myOsUnLinkEx(const char * pcPath);

/**
 * @brief rmdir
 */
extern int myOsRmdirEx(const char * pcPath);

/**
 * @brief truncate
 */
extern int myOsTruncateEx(const char * pcPath, int64 sz);

/**
 * @brief mkdir
 */
extern bool myOsCreateDirEx(const char* pszPath, int ch_split = '/');

/**
 * @brief copy file
 */
extern int myOsCopyFile(const char * src_path, const char * dst_path);

/**
 * @brief read hole file content
 */
extern int myOsReadHoleFileEx(const char * path, std::vector<uint8>& vcontent);

#endif

























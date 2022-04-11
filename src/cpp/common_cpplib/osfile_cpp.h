
/**
 * @file myOsFile.h ��װ��ͬϵͳ���ļ������ӿ� 2008-1-30 00:43
 *
 * @author lin shao chuan (email:lsccsl@tom.com, msn:lsccsl@163.net)
 * @blog http://blog.csdn.net/lsccsl
 *
 * @brief if it works, it was written by lin shao chuan, if not, i don't know who wrote it.
 *        ��װ��ͬϵͳ���ļ������ӿ�
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
 * @brief ���ļ�,��д
 */
extern HMYOSFILE myOsFileOpenReadWriteEx(const char * file_name);

/**
 * @brief ��ֻ���ķ�ʽ���ļ�
 */
extern HMYOSFILE myOsFileOpenReadOnlyEx(const char * file_name);

/**
 * @brief �Զ�ռ�ķ�ʽ���ļ�
 */
extern HMYOSFILE myOsFileOpenExclusiveEx(const char * file_name);

/**
 * @brief ��������
 */
extern HMYOSFILE myOSFileOpenEx(HMYOSFILE hf);

/**
 * @brief �ر��ļ�
 */
extern int myOsFileCloseEx(HMYOSFILE hf);

/**
 * @brief ͬ���ļ�������
 * @return 0:�ɹ� -1:ʧ��
 */
extern int myOsFileSynEx(HMYOSFILE hf);

/**
 * @brief д�ļ�
 * @return 0:�ɹ� -1:ʧ��, -2:δд��ָ���ֽ�
 */
extern int myOsFileWriteEx(HMYOSFILE hf, const void * data, size_t data_size, size_t * write_size);

/**
 * @brief ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
extern int myOsFileReadEx(HMYOSFILE hf, void * data, size_t data_size, size_t * read_size);

/**
 * @brief �ƶ���ǰ���ļ�ָ����off_set(������ļ�ͷ)
 * @return 0:�ɹ� -1:ʧ��
 */
extern int myOsFileSeekEx(HMYOSFILE hf, int64 off_set);

/**
 * @brief ɾ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
extern int myOsFileDelEx(const char * file_name);

/**
 * @brief ��ȡ�ļ��Ĵ�С
 * @return 0:�ɹ� -1:ʧ��
 */
extern int myOsFileSizeEx(HMYOSFILE hf, int64 * file_size);

/**
 * @brief ��ȡ�ļ����
 */
extern int myOsFileGetFdEx(HMYOSFILE hf);

/**
 * @brief �ж��ļ��Ƿ����
 * @return 0:�ļ������� ����:�ļ�����
 */
extern int myOsFileExistsEx(const char * file_name);

/**
 * @brief �Լ��ļ�
 * @return 0:�ɹ�, -1:ʧ��
 */
extern int myOsFileTruncateEx(HMYOSFILE hf, int64 nByte);

/**
 * @brief �������ļ�
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

























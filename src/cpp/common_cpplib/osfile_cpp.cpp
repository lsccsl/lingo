/**
 * @file myOsFile.c ��װ��ͬϵͳ���ļ������ӿ� 2008-1-31 00:43
 *
 * @author lin shao chuan (email:lsccsl@tom.com, msn:lsccsl@163.net)
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
#include "osfile_cpp.h"
#include <stdio.h>
#include <errno.h>
#include <string>
#include <string.h>

#ifdef WIN32
#include <assert.h>
#include <windows.h>
#include <direct.h>
#include <stdlib.h>
#include <io.h>

#ifndef INVALID_SET_FILE_POINTER
#define INVALID_SET_FILE_POINTER -1
#endif


typedef struct __os_file_t_
{
	/* �ļ���� */
	HANDLE hfile;

	/* �Ƿ�ռ�������� */
	int bOwner;
}os_file_t;


/**
 * @brief ���ļ�,��д
 */
HMYOSFILE myOsFileOpenReadWriteEx(const char * file_name)
{
	os_file_t * f = NULL;

	if(NULL == file_name)
		return NULL;

	f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

        f->hfile = CreateFileA(file_name,
		GENERIC_READ | GENERIC_WRITE,
		FILE_SHARE_READ | FILE_SHARE_WRITE ,
		NULL,
		OPEN_ALWAYS,
		FILE_ATTRIBUTE_NORMAL | FILE_FLAG_RANDOM_ACCESS,
		NULL);

	if(INVALID_HANDLE_VALUE == f->hfile)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief ��ֻ���ķ�ʽ���ļ�
 */
HMYOSFILE myOsFileOpenReadOnlyEx(const char * file_name)
{
	os_file_t * f = NULL;

	if(NULL == file_name)
		return NULL;

	f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

        f->hfile = CreateFileA(file_name,
		GENERIC_READ,
		FILE_SHARE_READ | FILE_SHARE_WRITE ,
		NULL,
		OPEN_EXISTING,
		FILE_ATTRIBUTE_NORMAL | FILE_FLAG_RANDOM_ACCESS,
		NULL);

	if(INVALID_HANDLE_VALUE == f->hfile)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief �Զ�ռ�ķ�ʽ���ļ�
 */
HMYOSFILE myOsFileOpenExclusiveEx(const char * file_name)
{
	os_file_t * f = NULL;

	if(NULL == file_name)
		return NULL;

	f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

        f->hfile = CreateFileA(file_name,
		GENERIC_READ | GENERIC_WRITE,
		0,
		NULL,
		CREATE_ALWAYS,
		FILE_FLAG_RANDOM_ACCESS,
		NULL);

	if(INVALID_HANDLE_VALUE == f->hfile)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief ��������
 */
HMYOSFILE myOSFileOpenEx(HMYOSFILE hf)
{
	os_file_t * f = NULL;

	if(NULL == hf)
		return NULL;

	f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

	f->bOwner = 0;
	f->hfile = hf->hfile;
	return f;
}

/**
 * @brief �ر��ļ�
 */
int myOsFileCloseEx(HMYOSFILE hf)
{
	if(NULL == hf)
		return -1;

	if(hf->bOwner)
	{
		if(0 == CloseHandle(hf->hfile))
			return -1;
	}

	free(hf);

	return 0;
}

/**
 * @brief ͬ���ļ�������
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSynEx(HMYOSFILE hf)
{
	if(NULL == hf || INVALID_HANDLE_VALUE == hf->hfile)
		return -1;

	if(FlushFileBuffers(hf->hfile))
		return 0;

	return -1;
}

/**
 * @brief д�ļ�
 * @return 0:�ɹ� -1:ʧ��, -2:δд��ָ���ֽ�
 */
int myOsFileWriteEx(HMYOSFILE hf, const void * data, size_t data_size, size_t * write_size)
{
	int rc = 0;
	DWORD wrote = 0;
	DWORD total_wrote = 0;

	if(NULL == hf || INVALID_HANDLE_VALUE == hf->hfile)
		return -1;

	assert(data_size > 0);

	if(write_size)
		*write_size = 0;

	while(data_size 
		&& (rc = WriteFile(hf->hfile, data, (DWORD)data_size, &wrote, 0))
		&& wrote > 0 
		&& wrote <= data_size)
	{
		data_size -= wrote;

		total_wrote += wrote;

		data = &((char*)data)[wrote];
	}

	if(!total_wrote)
		return -1;

	if(write_size)
		*write_size = total_wrote;

	if(!rc || data_size > wrote)
		return -2;

	return 0;
}

/**
 * @brief ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileReadEx(HMYOSFILE hf, void * data, size_t data_size, size_t * read_size)
{
	DWORD got;
	if(NULL == hf || INVALID_HANDLE_VALUE == hf->hfile)
		return -1;

	if(read_size)
		*read_size = 0;

	if(!ReadFile(hf->hfile, data, (DWORD)data_size, &got, 0))
	{
		got = 0;
		return -1;
	}

	if(read_size)
		*read_size = got;

	return 0;
}

/**
 * @brief �ƶ���ǰ���ļ�ָ����off_set(������ļ�ͷ)
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSeekEx(HMYOSFILE hf, int64 off_set)
{
	DWORD rc;
	LONG upperBits = (LONG)(off_set>>32);
	LONG lowerBits = (LONG)(off_set & 0xffffffff);

	if(NULL == hf || INVALID_HANDLE_VALUE == hf->hfile)
		return -1;

	rc = SetFilePointer(hf->hfile, lowerBits, &upperBits, FILE_BEGIN);

	if(rc==INVALID_SET_FILE_POINTER && GetLastError()!=NO_ERROR)
		return -1;

	return 0;
}

/**
 * @brief ɾ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileDelEx(const char * file_name)
{
	if(NULL == file_name)
		return -1;

        if(0 == DeleteFileA(file_name))
		return -1;

	return 0;
}

/**
 * @brief ��ȡ�ļ��Ĵ�С
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSizeEx(HMYOSFILE hf, int64 * file_size)
{
	DWORD upperBits, lowerBits;

	if(NULL == hf || NULL == file_size)
		return -1;

	lowerBits = GetFileSize(hf->hfile, &upperBits);
	*file_size = (((int64)upperBits)<<32) + lowerBits;

	return 0;
}

/**
 * @brief �ж��ļ��Ƿ����
 * @return 0:�ļ������� ����:�ļ�����
 */
int myOsFileExists(const char * file_name)
{
        return GetFileAttributesA(file_name) != 0xffffffff;
}

/**
 * @brief �Լ��ļ�
 * @return 0:�ɹ�, -1:ʧ��
 */
int myOsFileTruncateEx(HMYOSFILE hf, int64 nByte)
{
	LONG upperBits = (LONG)(nByte>>32);

	if(NULL == hf)
		return -1;

	SetFilePointer(hf->hfile, (LONG)nByte, &upperBits, FILE_BEGIN);
	SetEndOfFile(hf->hfile);

	return 0;
}
#else


#include <assert.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

#ifndef O_BINARY
# define O_BINARY 0
#endif


typedef struct __os_file_t_
{
	/* �ļ���� */
	int fd;

	/* ��ǰ�ļ��Ķ�дƫ�� */
	int64 off_set;

	/* �Ƿ�ռ�������� */
	int bOwner;
}os_file_t;


/**
 * @brief ���ļ�,��д
 */
HMYOSFILE myOsFileOpenReadWriteEx(const char * file_name)
{
	os_file_t * f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

	f->off_set = 0;
	f->fd = open(file_name, O_RDWR | O_CREAT | O_BINARY, 0644);
	if(f->fd < 0)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief ��ֻ���ķ�ʽ���ļ�
 */
HMYOSFILE myOsFileOpenReadOnlyEx(const char * file_name)
{
	os_file_t * f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

	f->off_set = 0;
	f->fd = open(file_name, O_RDONLY | O_BINARY);
	if(f->fd < 0)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief �Զ�ռ�ķ�ʽ���ļ�
 */
HMYOSFILE myOsFileOpenExclusiveEx(const char * file_name)
{
	os_file_t * f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

	f->off_set = 0;
	f->fd = open(file_name, O_EXCL | O_RDWR | O_CREAT | O_BINARY, 0644);
	if(f->fd < 0)
	{
		free(f);
		return NULL;
	}

	f->bOwner = 1;
	return f;
}

/**
 * @brief ��������
 */
HMYOSFILE myOSFileOpenEx(HMYOSFILE hf)
{
	os_file_t * f = NULL;

	if(NULL == hf)
		return NULL;

	f = (os_file_t *)malloc(sizeof(*f));
	if(NULL == f)
		return NULL;

	f->bOwner = 0;
	f->fd = hf->fd;
	f->off_set = hf->off_set;
	return f;
}

/**
 * @brief �ر��ļ�
 */
int myOsFileCloseEx(HMYOSFILE hf)
{
	if(NULL == hf)
		return -1;

	if(hf->bOwner)
	{
		if(0 != close(hf->fd))
			return -1;
	}

	free(hf);

	return 0;
}

/**
 * @brief ͬ���ļ�������
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSynEx(HMYOSFILE hf)
{
	/*
	* fdatasync:ֻͬ���ļ�������,������ĳЩϵͳ֧��,���freebsd, mac os x10.3
	* fsync:������fdatasync����,����ͬ���ļ�����,����ͬ���ļ�������(����ļ��޸�ʱ��֮���)
	* fcntl(fd, F_FULLFSYNC, 0): �ƺ�ֻ��mac os x֧��
	*/
	if(NULL == hf || hf->fd < 0)
		return -1;

	if(0 != fsync(hf->fd))
		return -1;

	return 0;
}


/**
 * @brief ͬ���ļ�������
 * @return 0:�ɹ� -1:ʧ��
 */
static int seek_and_write(os_file_t * f, const void * data, size_t data_size)
{
	int wrote = 0;

	assert(f && f->fd > 0);

	lseek(f->fd, f->off_set, SEEK_SET);

	wrote = write(f->fd, data, data_size);
	if(wrote > 0)
		f->off_set += wrote;

	return wrote;
}

/**
 * @brief д�ļ�
 * @return 0:�ɹ� -1:ʧ��, -2:δд��ָ���ֽ�
 */
int myOsFileWriteEx(HMYOSFILE hf, const void * data, size_t data_size, size_t * write_size)
{
	int wrote = 0;
	size_t total_wrote = 0;

	if(NULL == hf || hf->fd < 0)
		return -1;

	if(NULL == data || 0 == data_size)
		return -1;

	if(write_size)
		*write_size = 0;

	while(data_size > 0 && (wrote = seek_and_write(hf, data, data_size)) > 0)
	{
		total_wrote += wrote;

		data_size -= wrote;
		data = &((unsigned char*)data)[wrote];
	}

	if(write_size)
		*write_size = total_wrote;

	if(data_size > 0)
	{
		if(wrote < 0)
			return -1;
		else
			return -2;
	}

	return 0;
}


/**
 * @brief ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
static int seek_and_read(os_file_t * f, void *pBuf, int cnt)
{
	int got;

	assert(f && f->fd > 0);

	lseek64(f->fd, f->off_set, SEEK_SET);

	got = read(f->fd, pBuf, cnt);
	if(got > 0)
		f->off_set += got;

	return got;
}

/**
 * @brief ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileReadEx(HMYOSFILE hf, void * data, size_t data_size, size_t * read_size)
{
	int got;

	if(NULL == hf || hf->fd < 0)
		return -1;

	if(NULL == data || 0 == data_size)
		return -1;

	if(read_size)
		*read_size = 0;

	got = seek_and_read(hf, data, data_size);

	if(got > 0 && read_size)
		*read_size = got;

	if(got >= 0)
		return 0;
	else
		return -1;
}

/**
 * @brief �ƶ���ǰ���ļ�ָ����off_set(������ļ�ͷ)
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSeekEx(HMYOSFILE hf, int64 off_set)
{
	if(NULL == hf || hf->fd < 0)
		return -1;

	hf->off_set = off_set;

	return 0;
}

/**
 * @brief ɾ���ļ�
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileDelEx(const char * file_name)
{
	if(0 != unlink(file_name))
		return -1;

	return 0;
}

/**
 * @brief ��ȡ�ļ��Ĵ�С
 * @return 0:�ɹ� -1:ʧ��
 */
int myOsFileSizeEx(HMYOSFILE hf, int64 * file_size)
{
	int rc;
	struct stat buf;

	if(NULL == hf || hf->fd < 0)
		return -1;

	rc = fstat(hf->fd, &buf);

	if(rc != 0)
		return -1;

	*file_size = buf.st_size;

	return 0;
}

/**
 * @brief �Լ��ļ�
 * @return 0:�ɹ�, -1:ʧ��
 */
int myOsFileTruncateEx(HMYOSFILE hf, int64 nByte)
{
	if(NULL == hf || hf->fd < 0)
		return -1;

	if(0 != ftruncate(hf->fd, nByte))
		return -1;

	return 0;
}

#endif

/**
 * @brief �������ļ�
 */
int myOsRenameEx(const char * old_name, const char * new_name)
{
	if(NULL == old_name || NULL == new_name)
		return -1;

	return rename(old_name, new_name);
}

/**
 * @brief unlink
 */
int myOsUnLinkEx(const char * pcPath)
{
	if(NULL == pcPath)
		return -1;

#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
	return unlink(pcPath);
}

/**
 * @brief rmdir
 */
int myOsRmdirEx(const char * pcPath)
{
	if(NULL == pcPath)
		return -1;

#ifdef WIN32
	return rmdir(pcPath);
#else
	std::string cmd = "rm -fr ";

	cmd.append(pcPath);

	int ret = system(cmd.c_str());
	if(0 != ret)
		return -1;
	return 0;
#endif
}

/**
 * @brief truncate
 */
int myOsTruncateEx(const char * pcPath, int64 sz)
{
#ifdef WIN32

	HMYOSFILE hf1 = myOsFileOpenReadWriteEx(pcPath);
	myOsFileTruncateEx(hf1, sz);
	myOsFileCloseEx(hf1);
	return 0;

#else

	if(NULL == pcPath)
		return -1;

	return truncate(pcPath, sz);

#endif
}

/**
 * @brief ��ȡ�ļ����
 */
int myOsFileGetFdEx(HMYOSFILE hf)
{
	if(NULL == hf)
		return -1;

#ifdef WIN32
	return (int)hf->hfile;
#else
	return hf->fd;
#endif
}

/**
 * @brief mkdir
 */
bool myOsCreateDirEx(const char* pszPath, int ch_split)
{
#ifdef linux
    int iRet = mkdir(pszPath, 777);
#else
    int iRet = mkdir(pszPath);
#endif
    if (0 == iRet || errno == EEXIST)
    {
        return true;
    }

    const char* p = strrchr(pszPath, ch_split);
	if (NULL != p && myOsCreateDirEx(std::string(pszPath, p - pszPath).c_str(), ch_split))
    {
#ifdef linux
        return 0 == mkdir(pszPath, 777);
#else
        return 0 == mkdir(pszPath);
#endif
    }

    return false;
}

/**
 * @brief copy file
 */
int myOsCopyFile(const char * src_path, const char * dst_path)
{
#ifndef WIN32
	std::string cmd = "cp ";
#else
	std::string cmd = "copy ";
#endif

	cmd.append(src_path);
	cmd.append(" ");
	cmd.append(dst_path);

	int ret = system(cmd.c_str());
	if(0 != ret)
		return -1;
	return 0;
}

/**
 * @brief �ж��ļ��Ƿ����
 * @return 0:�ļ������� ����:�ļ�����
 */
int myOsFileExistsEx(const char * file_name)
{
	return access(file_name, 0)==0;
}

/**
 * @brief read hole file content
 */
int myOsReadHoleFileEx(const char * path, std::vector<uint8>& vcontent)
{
	/* �����ļ� */
	HMYOSFILE hfile = myOsFileOpenReadWriteEx(path);
	if(NULL == hfile)
	{
		return -1;
	}

	int64 sz = 0;
	if(0 != myOsFileSizeEx(hfile, &sz))
	{
		return -1;
	}

	if(sz > 0)
	{
		vcontent.resize(sz);
		size_t read_sz = 0;
		myOsFileReadEx(hfile, &vcontent[0], vcontent.size(), &read_sz);
		if(sz != vcontent.size())
		{
			return -1;
		}
	}

	return 0;
}
























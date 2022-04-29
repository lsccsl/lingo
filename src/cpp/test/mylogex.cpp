#include "mylogex.h"
//#include "channel.h"
extern "C"
{
	//#include "mythread.h"
	//#include "mymsgque.h"
	#include "gettimeofday.h"
};

#include <ctype.h>
#include <time.h>
#include <sys/stat.h>
#ifndef WIN32
#include <syslog.h>
#include <iconv.h>
#include <pthread.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/time.h>
#endif

#pragma   warning(   disable   :   4996)
#pragma   warning(   disable   :   4005)

#ifdef WIN32
#include <windows.h>
#define vsnprintf _vsnprintf
#include <assert.h>
#endif

static char g_acPre[128] = "log";
void MYLOG_PRE(const char * pcPre)
{
	strncpy(g_acPre, pcPre, sizeof(g_acPre) - 1);
	g_acPre[sizeof(g_acPre) - 1] = 0;
}

static int bsyn_ = 1;

int __g_level_ = 0;
#define __LOG_TO_FILE_ (0x01)
#define __LOG_TO_SCREEN_ (0x02)
#define __LOG_TO_SRV_ (0x04)
unsigned int __g_direction_ = __LOG_TO_SCREEN_;
//static HMYTHREAD hlog_thread_ = NULL;
//static HMYMSGQUE hlog_que_ = NULL;

struct __log_client_t_
{
	char ac_srv_ip[32];
	int port;
	volatile int fd;
};
static __log_client_t_ __g_lci_ = {
	"192.168.21.127",
	3353,
	-1,
};

/**
 * @brief 写往日志服务器
 */
struct log_msg_t
{
	unsigned int log_flag;
	unsigned int thrd;
	char pre[32];
	char log_con[1];
};
void write_to_log_srv(const char * log, const char * pre = NULL)
{
//	if(__g_lci_.fd < 0)
//		__g_lci_.fd = CChannel::UdpOpen("0.0.0.0", __g_lci_.port + 1);
//
//	if(NULL == pre)
//		pre = g_acPre;
//
//	log_msg_t * mg = (log_msg_t *)malloc(strlen(log) + 1 + sizeof(*mg));
//#ifndef WIN32
//	mg->thrd = pthread_self();
//#endif
//	memcpy(mg->log_con, log, strlen(log));
//	memcpy(mg->pre, pre, sizeof(mg->pre) - 1);
//	mg->pre[3] = 0;
//
//	CChannel::UdpWrite(__g_lci_.fd, mg/*log*/, strlen(log) + sizeof(mg->thrd) + sizeof(mg->pre)/*strlen(log)*/,
//		__g_lci_.ac_srv_ip, __g_lci_.port);
//
//	free(mg);
}

/**
 * @brief 设置日志服务器的ip:port
 */
void MYLOG_SET_LOG_SRV(const char * ip, int port)
{
	if(ip)
		strncpy(__g_lci_.ac_srv_ip, ip, sizeof(__g_lci_.ac_srv_ip) - 1);
	__g_lci_.port = port;
}

/**
 * @brief 设置日志级别
 */
void MYLOG_LEVEL(int level)
{
	__g_level_ = level;
}

/**
 * @brief 设置日志文件大小
 */
static int unsigned long __g_fileSz_ = 10 * 1024 * 1024;
extern void MYLOG_SETLOGFILESZ(unsigned long nFileSZ_In_M)
{
	__g_fileSz_ = nFileSZ_In_M * 1024 * 1024;
	if(__g_fileSz_ < 1024 * 1024)
		__g_fileSz_ = 1024 * 1024;
}

/**
 * @brief 设置日志文件的输出方向
 */
extern void MYLOG_SET_LOG_DIRECTION(unsigned int d)
{
	__g_direction_ = d;
}

/**
 * @brief 设置日志前缀
 */
void MYLOG_GET_PRE(char * pre, int pre_sz)
{
	strncpy(pre, g_acPre, pre_sz - 1);
	pre[pre_sz - 1] = 0;
}

int MYLOG_GETLEVEL()
{
	return __g_level_;
}

int MYLOG_GET_LOG_DIRECTION()
{
	return __g_direction_;
}



static unsigned long get_file_size(const char *filename)
{
	struct stat buf;
	if(stat(filename, &buf)<0)
	{
		return 0;
	}
	return (unsigned long)buf.st_size;
}
static void backup_log(const char * filename)
{
#ifndef WIN32
	char acfile_name[128+4] = {0};
	snprintf(acfile_name, sizeof(acfile_name) - 1, "%s.bak", filename);

	char acCmd[256] = {0};
	snprintf(acCmd, sizeof(acCmd) - 1, "rm -fr %s", acfile_name);
	system(acCmd);	
	snprintf(acCmd, sizeof(acCmd) - 1, "mv %s %s", filename, acfile_name);
	system(acCmd);
#else
	char acfile_name[128+4] = {0};
	_snprintf(acfile_name, sizeof(acfile_name) - 1, "%s.bak", filename);

	rename(filename, acfile_name);
#endif
}

void WRITE_FILE(const char * x, const char * pre = NULL, unsigned int thrdid = 0)
{
#ifndef WIN32
	if(NULL == pre)
		pre = g_acPre;

	char acfile_name[128] = {0};
	snprintf(acfile_name, sizeof(acfile_name) - 1, "%s_%x.log", pre, thrdid ? thrdid : pthread_self());
	if(get_file_size(acfile_name) > __g_fileSz_)
	{backup_log(acfile_name);}

	FILE * pfile = fopen(acfile_name, "aw");
	if(pfile)
	{
		fprintf(pfile, x);
		fclose(pfile);
	}
#else
	if(NULL == pre)
		pre = g_acPre;

	char acfile_name[128] = {0};
	_snprintf(acfile_name, sizeof(acfile_name) - 1, "%s_%x.log", pre, thrdid ? thrdid : GetCurrentThreadId());
	if(get_file_size(acfile_name) > __g_fileSz_)
	{backup_log(acfile_name);}

        FILE * pfile = fopen(acfile_name, "a+");
	if(pfile)
	{
		fprintf(pfile, x);
                OutputDebugStringA(x);
		fclose(pfile);
	}
#endif
}


void WRITE_FILE_BIN(const void * x, size_t l, const char * pre = NULL,  unsigned int thrdid = 0)
{
#ifndef WIN32
	if(NULL == pre)
		pre = g_acPre;

	char acfile_name[128] = {0};
	snprintf(acfile_name, sizeof(acfile_name) - 1, "%s_%x.log", pre, thrdid ? thrdid : pthread_self());
	FILE * pfile = fopen(acfile_name, "aw");
	if(pfile)
	{
		fwrite(x, 1, l, pfile);
		fclose(pfile);
	}
#else
	if(NULL == pre)
		pre = g_acPre;

	char acfile_name[128] = {0};
	_snprintf(acfile_name, sizeof(acfile_name) - 1, "%s_%x.log", pre, thrdid ? thrdid : GetCurrentThreadId());
        FILE * pfile = fopen(acfile_name, "a+");
	if(pfile)
	{
		fwrite(x, 1, l, pfile);
		fclose(pfile);
	}
#endif
}


void push_to_log_que(int log_flag, const char * pre, const char * con)
{
	assert(pre && con);
	log_msg_t * msg = (log_msg_t *)malloc(sizeof(*msg) + strlen(con) + 1);

	msg->log_flag = log_flag;
	strncpy(msg->pre, pre, sizeof(msg->pre));
#ifndef WIN32
	msg->thrd = 1;
#else
	msg->thrd = 1;
#endif
	strncpy(msg->log_con, con, strlen(con) + 1);

	//MyMsgQuePush_block(hlog_que_, msg);
}

void mylog_err(char * fmt, ...)
{
	va_list var;

	char actemp[1024] = {0};	
	va_start(var, fmt);
	vsnprintf(actemp, sizeof(actemp) - 1, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;31m");
#endif
			printf(actemp);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_ERR, g_acPre, actemp);
	}
}

void mylog_warn(char * fmt, ...)
{
	va_list var;

	char actemp[1024] = {0};
	va_start(var, fmt);
	vsnprintf(actemp, sizeof(actemp) - 1, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;33m");
#endif
			printf(actemp);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_WARN, g_acPre, actemp);
	}
}

void mylog_info(char * fmt, ...)
{
	va_list var;

	char actemp[1024] = {0};
	va_start(var, fmt);
	vsnprintf(actemp, sizeof(actemp) - 1, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;32m");
#endif
			printf(actemp);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_INFO, g_acPre, actemp);
	}
}

void mylog_trace(char * fmt, ...)
{
	va_list var;

	char actemp[1024] = {0};
	va_start(var, fmt);
	vsnprintf(actemp, sizeof(actemp) - 1, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
			printf(actemp);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_TRACE, g_acPre, actemp);
	}
}

void mylog_debug(char * fmt, ...)
{
	va_list var;

	char actemp[1024] = {0};
	va_start(var, fmt);
	vsnprintf(actemp, sizeof(actemp) - 1, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;34m");
#endif
			printf(actemp);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_DEBUG, g_acPre, actemp);
	}
}

void mylog_debug_ex(const char * pcpre, char * log)
{
	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;34m");
#endif
			printf(log);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(log, pcpre);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(log, pcpre);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_DEBUG, pcpre, log);
	}
}
void mylog_trace_ex(const char * pcpre, char * log)
{
	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
			printf(log);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(log, pcpre);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(log, pcpre);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_TRACE, pcpre, log);
	}
}
void mylog_info_ex(const char * pcpre, char * log)
{
	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;32m");
#endif
			printf(log);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(log, pcpre);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(log, pcpre);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_INFO, pcpre, log);
	}
}
void mylog_warn_ex(const char * pcpre, char * log)
{
	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;33m");
#endif
			printf(log);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(log, pcpre);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(log, pcpre);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_WARN, pcpre, log);
	}
}
void mylog_err_ex(const char * pcpre, char * log)
{
	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;31m");
#endif
			printf(log);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(log, pcpre);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(log, pcpre);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_ERR, pcpre, log);
	}
}

void mylog_dump_bin(const char * file, const int line, const void * buf, int buf_sz)
{
	if(__g_level_ > 0)
		return;

	int i = 0;
	char actemp[8192] = {0};
	int pos = 0;

	pos += sprintf(actemp + pos, "\r\n[%s:%d] len:%d", file, line, buf_sz);
	if(NULL == buf)
	{
		if(bsyn_)
		{
			if(__g_direction_ & __LOG_TO_SCREEN_)
				printf(actemp);
			if(__g_direction_ & __LOG_TO_FILE_)
				WRITE_FILE_BIN(actemp, pos);
			if(__g_direction_ & __LOG_TO_SRV_)
				write_to_log_srv(actemp);
			return;
		}
		else
		{
			push_to_log_que(MYLOG_FLAG_DEBUG, g_acPre, actemp);
		}
	}

	if(buf_sz > 1024)
		buf_sz = 1024;

	for(; i < buf_sz; i ++)
	{
		if(0 == i % 16)
			pos += sprintf(actemp + pos, "\r\n");

		pos += sprintf(actemp + pos, "%2x ", *((unsigned char *)buf + i));
	}

	pos += sprintf(actemp + pos, "\r\n");

	for(i = 0; i < buf_sz; i ++)
	{
		if(0 == i % 16)
			pos += sprintf(actemp + pos, "\r\n");

		if(isdigit(*((unsigned char *)buf + i)))
			pos += sprintf(actemp + pos, "%c", *((unsigned char *)buf + i));
		else if(isalpha(*((unsigned char *)buf + i)))
			pos += sprintf(actemp + pos, "%c", *((unsigned char *)buf + i));
		//else if(*((unsigned char *)buf + i) == '\r' || *((unsigned char *)buf + i) == '\n' || *((unsigned char *)buf + i) == '\t' || *((unsigned char *)buf + i) == ' ')
		//	pos += sprintf(actemp + pos, ".");
		//else if(isgraph(*((unsigned char *)buf + i)))
		//	pos += sprintf(actemp + pos, "%c", *((unsigned char *)buf + i));
		else
			pos += sprintf(actemp + pos, ".");
	}

	pos += sprintf(actemp + pos, "\r\n\r\n");

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
			printf(actemp);
		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE_BIN(actemp, pos);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_DEBUG, g_acPre, actemp);
	}
}

void get_log_time(char * pctemp, int& len)
{
	#ifdef WIN32
	timeval tv;
	gettimeofday(&tv, NULL);
        time_t temp = tv.tv_sec;
        struct tm * ptm = localtime(&temp);

	len = sprintf(pctemp, "%d-%d-%d %d:%d:%d %6d",
		ptm->tm_year + 1900,
		ptm->tm_mon + 1,
		ptm->tm_mday,
		ptm->tm_hour,
		ptm->tm_min,
		ptm->tm_sec,
		tv.tv_usec);
	#else
	timeval tv;
	gettimeofday(&tv, NULL);
	struct tm * ptm = localtime(&tv.tv_sec);
	len =  snprintf(pctemp, len - 1, "%d-%d-%d %d:%d:%d %6d ",
		ptm->tm_year + 1900,
		ptm->tm_mon + 1,
		ptm->tm_mday,
		ptm->tm_hour,
		ptm->tm_min,
		ptm->tm_sec,
		tv.tv_usec);
	#endif
}

void MYLOG_X(int level, const char * fmt, ...)
{
	if(__g_level_ > level)
		return;

	va_list var;

	char actemp[1024] = {"\r\n"};
	va_start(var, fmt);
	int n = vsnprintf(actemp + 2, sizeof(actemp) - 3, fmt, var);
	va_end(var);

	if(bsyn_)
	{
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			printf("\033[1;35m");
#endif
			vprintf(fmt, var);
			printf("\r\n");
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(actemp);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(actemp);
	}
	else
	{
		push_to_log_que(MYLOG_FLAG_ERR, g_acPre, actemp);
	}
}

void * log_thrd_fun(void * param)
{
#if 0
	log_msg_t * msg;
	while(1)
	{
		msg = (log_msg_t *)MyMsgQuePop_block(hlog_que_);
		if(__g_direction_ & __LOG_TO_SCREEN_)
		{
#ifndef WIN32
			switch(msg->log_flag)
			{
			case MYLOG_FLAG_DEBUG:
				printf("\033[1;34m");
				break;

			case MYLOG_FLAG_TRACE:
				break;

			case MYLOG_FLAG_INFO:
				printf("\033[1;32m");
				break;

			case MYLOG_FLAG_WARN:
				printf("\033[1;33m");
				break;

			case MYLOG_FLAG_ERR:
				printf("\033[1;31m");
				break;

			default:
				break;
			}
#endif
			printf(msg->log_con);
#ifndef WIN32
			printf("\033[0m");
#endif
		}

		if(__g_direction_ & __LOG_TO_FILE_)
			WRITE_FILE(msg->log_con, msg->pre, msg->thrd);
		if(__g_direction_ & __LOG_TO_SRV_)
			write_to_log_srv(msg->log_con, msg->pre);

		free(msg);
	}
#endif
	return NULL;
}

/* @brief 设日志输入设为异步的 */
void MYLOG_AYN()
{
#if 0
	bsyn_ = 0;
	if(NULL == hlog_que_)
	{
		hlog_que_ = MyMsgQueConstruct(NULL, 65535);
	}
	if(NULL == hlog_thread_)
	{
		hlog_thread_ = MyThreadConstruct(log_thrd_fun, NULL, 0, NULL);
	}
#endif
}

/* @brief 将日志输出设为同步的 */
void MYLOG_SYN()
{
	bsyn_ = 1;
}







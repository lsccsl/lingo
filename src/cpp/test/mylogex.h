#ifndef __MYLOG_H__
#define __MYLOG_H__

#include <stdarg.h>
#include <string.h>
#include <stdio.h>

extern void MYLOG_X(int level, const char * fmt, ...);
extern void mylog_err(char * fmt, ...);
extern void mylog_warn(char * fmt, ...);
extern void mylog_info(char * fmt, ...);
extern void mylog_trace(char * fmt, ...);
extern void mylog_debug(char * fmt, ...);
extern void mylog_debug_ex(const char * pcpre, char * log);
extern void mylog_trace_ex(const char * pcpre, char * log);
extern void mylog_info_ex(const char * pcpre, char * log);
extern void mylog_warn_ex(const char * pcpre, char * log);
extern void mylog_err_ex(const char * pcpre, char * log);
extern void mylog_dump_bin(const char * file, const int line, const void * buf, int buf_sz);
//extern void myprinttime(int level);
void get_log_time(char * pctemp, int& len);
extern int __g_level_;

enum{
	MYLOG_FLAG_DEBUG = 0x01,
	MYLOG_FLAG_TRACE = 0x02,
	MYLOG_FLAG_INFO = 0x04,
	MYLOG_FLAG_WARN = 0x08,
	MYLOG_FLAG_ERR = 0x010,
};

struct __log_obj_t_
{
	__log_obj_t_(const char * file, int line):line_num(line)
	{
#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
		strncpy(file_name, file ? file : "null", sizeof(file_name) - 1);
		file_name[sizeof(file_name) - 1] = 0;
	}

	inline void __format_log_ex_(const char * pcpre, const char * fmt, ...)
	{
		if(pcpre)
		{
			strncpy(acpre, pcpre, sizeof(acpre) - 1);
			acpre[sizeof(acpre) - 1] = 0;
		}
		else
		{
			strncpy(acpre, "log_", sizeof(acpre) - 1);
			acpre[sizeof(acpre) - 1] = 0;
		}

		{
			int pos = sizeof(this->log_context);
			get_log_time(log_context, pos);
			if(pos >= (sizeof(log_context) - 1))
				return;
#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
#ifdef WIN32
			pos += _snprintf(log_context + pos, sizeof(log_context) - 3 - pos, "[%s:%d]", file_name, line_num);
#else
			pos += snprintf(log_context + pos, sizeof(log_context) - 3 - pos, "[%s:%d]", file_name, line_num);
#endif
			if(pos >= (sizeof(log_context) - 3))
			{
				log_context[sizeof(log_context) - 3] = '\r';
				log_context[sizeof(log_context) - 2] = '\n';
				log_context[sizeof(log_context) - 1] = 0;
				return;
			}

			va_list var;
			va_start(var, fmt);
#ifdef WIN32
			pos += _vsnprintf(log_context + pos, sizeof(log_context) - 3 - pos, fmt, var);
#else
			pos += vsnprintf(log_context + pos, sizeof(log_context) - 3 - pos, fmt, var);
#endif
			va_end(var);

			if(pos >= (sizeof(log_context) - 3))
			{
				log_context[sizeof(log_context) - 3] = '\r';
				log_context[sizeof(log_context) - 2] = '\n';
				log_context[sizeof(log_context) - 1] = 0;
			}
			else
			{
				log_context[pos + 0] = '\r';
				log_context[pos + 1] = '\n';
				log_context[pos + 2] = 0;
			}
		}
	}

	inline void __format_log_(const char * fmt, ...)
	{
		int pos = sizeof(this->log_context);
		get_log_time(log_context, pos);
		if(pos >= (sizeof(log_context) - 1))
			return;
#pragma   warning(   disable   :   4996) /* fuck vc,why warning? */ 
#ifdef WIN32
		pos += _snprintf(log_context + pos, sizeof(log_context) - 3 - pos, "[%s:%d]", file_name, line_num);
#else
		pos += snprintf(log_context + pos, sizeof(log_context) - 3 - pos, "[%s:%d]", file_name, line_num);
#endif
		if(pos >= (sizeof(log_context) - 3))
		{
			log_context[sizeof(log_context) - 3] = '\r';
			log_context[sizeof(log_context) - 2] = '\n';
			log_context[sizeof(log_context) - 1] = 0;
			return;
		}

		va_list var;
		va_start(var, fmt);
#ifdef WIN32
		pos += _vsnprintf(log_context + pos, sizeof(log_context) - 3 - pos, fmt, var);
#else
		pos += vsnprintf(log_context + pos, sizeof(log_context) - 3 - pos, fmt, var);
#endif
		va_end(var);

		if(pos >= (sizeof(log_context) - 3))
		{
			log_context[sizeof(log_context) - 3] = '\r';
			log_context[sizeof(log_context) - 2] = '\n';
			log_context[sizeof(log_context) - 1] = 0;
		}
		else
		{
			log_context[pos + 0] = '\r';
			log_context[pos + 1] = '\n';
			log_context[pos + 2] = 0;
		}
	}

	char file_name[64];
	int line_num;

	char acpre[16];

	char log_context[1024];
};

#ifdef _NOT_LOG

	#define MYLOG_ERR(x) do{}while(0)
	#define MYLOG_WARN(x) do{}while(0)
	#define MYLOG_INFO(x) do{}while(0)
	#define MYLOG_DEBUG(x) do{}while(0)
	#define MYLOG_TRACE(x) do{}while(0)
	#define MYLOG_DUMP_BIN(buf, buf_sz) do{}while(0)

#else

	#define MYLOG_ERR(x) do{\
		if(__g_level_ <= 4)\
			{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ x;\
			mylog_err(__lo_.log_context);\
		}\
	}while(0)

	#define MYLOG_WARN(x) do{\
		if(__g_level_ <= 3)\
		{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ x;\
			mylog_warn(__lo_.log_context);\
		}\
	}while(0)

	#define MYLOG_INFO(x) do{\
		if(__g_level_ <= 2)\
		{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ x;\
			mylog_info(__lo_.log_context);\
		}\
	}while(0)

	#define MYLOG_TRACE(x) do{\
		if(__g_level_ <= 1)\
			{\
				__log_obj_t_ __lo_(__FILE__, __LINE__);\
				__lo_.__format_log_ x;\
				mylog_trace(__lo_.log_context);\
			}\
	}while(0)

	#define MYLOG_DEBUG(x) do{\
		if(__g_level_ <= 0)\
		{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ x;\
			mylog_debug(__lo_.log_context);\
		}\
	}while(0)


	#define MYLOG_ERREX(x) do{\
		if(__g_level_ <= 4)\
			{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ex_ x;\
			mylog_err_ex(__lo_.acpre, __lo_.log_context);\
			}\
	}while(0)

	#define MYLOG_WARNEX(x) do{\
		if(__g_level_ <= 3)\
			{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ex_ x;\
			mylog_warn_ex(__lo_.acpre, __lo_.log_context);\
		}\
	}while(0)

	#define MYLOG_INFOEX(x) do{\
		if(__g_level_ <= 2)\
			{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ex_ x;\
			mylog_info_ex(__lo_.acpre, __lo_.log_context);\
		}\
	}while(0)

	#define MYLOG_TRACEEX(x) do{\
		if(__g_level_ <= 1)\
			{\
				__log_obj_t_ __lo_(__FILE__, __LINE__);\
				__lo_.__format_log_ex_ x;\
				mylog_trace_ex(__lo_.acpre, __lo_.log_context);\
			}\
	}while(0)


	#define MYLOG_DEBUGEX(x) do{\
		if(__g_level_ <= 0)\
		{\
			__log_obj_t_ __lo_(__FILE__, __LINE__);\
			__lo_.__format_log_ex_ x;\
			mylog_debug_ex(__lo_.acpre, __lo_.log_context);\
		}\
	}while(0)
	

	#define MYLOG_DUMP_BIN(buf, buf_sz)  do{mylog_dump_bin(__FILE__, __LINE__, buf, buf_sz);}while(0)
#endif

/* @brief 设日志输入设为异步的 */
extern void MYLOG_AYN();
/* @brief 将日志输出设为同步的 */
extern void MYLOG_SYN();
extern void MYLOG_PRE(const char * pcPre);
extern void MYLOG_LEVEL(int level);
extern void MYLOG_SETLOGFILESZ(unsigned long nFileSZ_In_M);
extern void MYLOG_SET_LOG_DIRECTION(unsigned int d);
extern void MYLOG_SET_LOG_SRV(const char * ip, int port);
extern void MYLOG_GET_PRE(char * pre, int pre_sz);
extern int MYLOG_GETLEVEL();
extern int MYLOG_GET_LOG_DIRECTION();

#endif

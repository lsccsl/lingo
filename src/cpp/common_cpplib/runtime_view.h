/**
 * @file runtime_view
 */
#ifndef __RUNTIME_VIEW_H__
#define __RUNTIME_VIEW_H__

extern "C"
{
	#include "mylisterner.h"
}
#include "CMySocket.h"

class rt_view
{
public:

	/**
	 * @brief constructor
	 */
	rt_view(const char * ipc_name = NULL);

	/**
	 * @brief destructor
	 */
	virtual ~rt_view();

	/**
	 * @brief 命令回调
	 */
	virtual void command(int cmd) = 0;

	/**
	 * @brief 送命令
	 */
	static void send_cmd(int cmd, const char * ipc_name = NULL);

protected:

	/**
	 * @brief 输入回调
	 */
	static int _handle_input(unsigned long context_data, int fd);

protected:

	/**
	 * @brief listerner
	 */
	HMYLISTERNER hlsn_;

	/**
	 * @brief unix sock
	 */
	CMyUnixSocket us_;
};

#endif

#include "os.h"
#ifndef WIN32
#include <errno.h>
#include <unistd.h>


int WSAGetLastError()
{
	return errno;
}

int GetLastError()
{
	return errno;
}

void Sleep(int usec)
{
	usleep(usec);
}

#endif

#ifndef WIN32
#include "myepoll_linux.h"
#else
#pragma   warning(   disable   :   4786)
#include "myepoll_win32.h"
#endif

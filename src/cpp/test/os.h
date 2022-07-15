#pragma once

#ifndef WIN32

int WSAGetLastError();
int GetLastError();
void Sleep(int usec);

#endif
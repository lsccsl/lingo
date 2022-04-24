#ifndef ____TYPE_DEF_H__
#define ____TYPE_DEF_H__
//#ifdef _MBCS
//	//typedef unsigned __int64 uint64;
//#else
	typedef unsigned long long int uint64;
//#endif
	typedef unsigned int uint32;
	typedef unsigned short uint16;
	typedef unsigned char uint8;

#ifdef _MBCS
	typedef __int64 int64;
#else
	typedef long long int int64;
#endif
	typedef int int32;
	typedef short int16;
	typedef char int8;
#endif


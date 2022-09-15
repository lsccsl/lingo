#ifndef __DLLEXPORT_H_06457146_757F_49CA_97CD_BB5B563FE66B
#define __DLLEXPORT_H_06457146_757F_49CA_97CD_BB5B563FE66B

#ifdef _WINDOWS

	#ifdef _DLL_EXPORTS_
		#define DLL_EXPORT __declspec(dllexport)
	#else
		#define DLL_EXPORT __declspec(dllimport)
	#endif

#else

	#define DLL_EXPORT

#endif

#endif

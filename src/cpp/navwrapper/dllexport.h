#pragma once

#ifdef _WINDOWS

	#ifdef _DLL_EXPORTS_
		#define DLL_EXPORT __declspec(dllexport)
	#else
		#define DLL_EXPORT __declspec(dllimport)
	#endif

#else

	#define DLL_EXPORT

#endif

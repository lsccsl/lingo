#pragma once

#ifdef _DLL_EXPORTS_
#define DLL_EXPORT __declspec(dllexport)
#else
#define DLL_EXPORT __declspec(dllimport)
#endif

#pragma once

#include "dllexport.h"
#include <stdbool.h>
#include "RecastWrapperDef.h"

#ifdef __cplusplus
extern "C" {
#endif

DLL_EXPORT void* nav_create(const char * file_path);

/* @brief pos_path nav_freepath */
DLL_EXPORT void nav_findpath(void* ins_ptr, struct RecastPos * startPos, struct RecastPos * endPos, struct RecastPos** pos_path, int* pos_path_sz, bool bprint);
DLL_EXPORT void nav_freepath(struct RecastPos* pos_path);

#ifdef __cplusplus
}
#endif

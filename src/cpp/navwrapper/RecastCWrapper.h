#pragma once

#include "dllexport.h"
#include <stdbool.h>
#include "RecastWrapperDef.h"

#ifdef __cplusplus
extern "C" {
#endif

DLL_EXPORT void* nav_create(const char * file_path);


DLL_EXPORT void nav_findpath(void* ins_ptr, struct RecastPos * startPos, struct RecastPos * endPos, struct RecastPos** pos_path, int* pos_path_sz, bool bprint);
DLL_EXPORT void nav_findpath1(void * ins_ptr, const float startPos[3], const float endPos[3], bool bprint);

#ifdef __cplusplus
}
#endif

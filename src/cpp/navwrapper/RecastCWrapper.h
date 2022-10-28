#pragma once

#include "dllexport.h"
#include <stdbool.h>
#include "RecastWrapperDef.h"

#ifdef __cplusplus
extern "C" {
#endif

DLL_EXPORT void* nav_create(const char* file_path);
DLL_EXPORT void* nav_new();
DLL_EXPORT void nav_delete(void* ins_ptr);
DLL_EXPORT void nav_update(void* ins_ptr);

DLL_EXPORT void nav_reset_agent(void* ins_ptr, float agentHeight, float agentRadius, float agentMaxClimb, float agentMaxSlope);
DLL_EXPORT bool nav_load(void* ins_ptr, const char* file_path);
DLL_EXPORT bool nav_load_from_template(void* ins_ptr, void* template_ptr);


/* @brief pos_path nav_freepath */
DLL_EXPORT void nav_findpath(void* ins_ptr, struct RecastVec3f* startPos, struct RecastVec3f* endPos, struct RecastVec3f** pos_path, int* pos_path_sz, bool bprint);
DLL_EXPORT void nav_freepath(struct RecastVec3f* pos_path);

DLL_EXPORT unsigned int nav_add_obstacle(void* ins_ptr, struct RecastVec3f* center, struct RecastVec3f* halfExtents, const float yRadians);
DLL_EXPORT void nav_del_obstacle(void* ins_ptr, unsigned int obj_id);


DLL_EXPORT void* nav_temlate_create(const char* file_path,
	float agentHeight, float agentRadius, float agentMaxClimb, float agentMaxSlope);

#ifdef __cplusplus
}
#endif

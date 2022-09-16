#include "RecastCWrapper.h"
#include "RecastInstance.h"

void* nav_create(const char* file_path)
{
	if (NULL == file_path)
		return NULL;
	RecastInstance* ins = new RecastInstance;
	if (!ins->LoadFromObj(file_path))
	{
		delete ins;
		return NULL;
	}

	return ins;
}

void nav_findpath(void* ins_ptr, const float startPos[3], const float endPos[3], bool bprint)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	ins->FindPath(startPos, endPos, bprint);
}

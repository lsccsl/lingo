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

void nav_findpath1(void* ins_ptr, const float startPos[3], const float endPos[3], bool bprint)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	ins->FindPath(startPos, endPos, bprint);
}

void nav_findpath(void* ins_ptr, RecastPos* p_startPos, RecastPos* p_endPos, struct RecastPos** pos_path, int* pos_path_sz, bool bprint)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	float startPos[3] = {p_startPos->x, p_startPos->y, p_startPos->z};
	float endPos[3] = { p_endPos->x, p_endPos->y, p_endPos->z };

	std::vector<RecastPos> vPos;
	ins->FindPath(startPos, endPos, vPos, bprint);
	if (vPos.empty())
		return;
	RecastPos * pos = *pos_path = (RecastPos*)malloc(sizeof(RecastPos) * vPos.size());
	for (auto& it : vPos)
	{
		pos->x = it.x;
		pos->y = it.y;
		pos->z = it.z;
		pos++;
	}
	*pos_path_sz = vPos.size();
}

void nav_freepath(RecastPos* pos_path)
{
	printf("\r\nnav_freepath\r\n");
	free(pos_path);
}

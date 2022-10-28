#include "RecastCWrapper.h"
#include "RecastInstance.h"
#include "RecastTemplate.h"

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

void* nav_new()
{
	return new RecastInstance;
}
void nav_delete(void* ins_ptr)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	delete ins;
}

bool nav_load(void* ins_ptr, const char* file_path)
{
	if (NULL == ins_ptr)
		return false;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return false;

	return ins->LoadFromObj(file_path);
}

void nav_reset_agent(void* ins_ptr, float agentHeight, float agentRadius, float agentMaxClimb, float agentMaxSlope)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	ins->reset_agent(agentHeight, agentRadius, agentMaxClimb, agentMaxSlope);
}

void nav_findpath(void* ins_ptr, RecastVec3f* p_startPos, RecastVec3f* p_endPos, struct RecastVec3f** pos_path, int* pos_path_sz, bool bprint)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	float startPos[3] = {p_startPos->x, p_startPos->y, p_startPos->z};
	float endPos[3] = { p_endPos->x, p_endPos->y, p_endPos->z };

	std::vector<RecastVec3f> vPos;
	ins->FindPath(startPos, endPos, vPos, bprint);
	if (vPos.empty())
		return;
	RecastVec3f* pos = *pos_path = (RecastVec3f*)malloc(sizeof(RecastVec3f) * vPos.size());
	for (auto& it : vPos)
	{
		pos->x = it.x;
		pos->y = it.y;
		pos->z = it.z;
		pos++;
	}
	*pos_path_sz = vPos.size();
}

void nav_freepath(RecastVec3f* pos_path)
{
	printf("\r\nnav_freepath\r\n");
	free(pos_path);
}

unsigned int nav_add_obstacle(void* ins_ptr, RecastVec3f* center, RecastVec3f* halfExtents, const float yRadians)
{
	if (NULL == ins_ptr)
		return 0;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return 0;

	return ins->add_obstacle_with_y_rotation(center, halfExtents, yRadians);
}

void nav_del_obstacle(void* ins_ptr, unsigned int obj_id)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	ins->del_obstacle(obj_id);
}

void nav_update(void* ins_ptr)
{
	if (NULL == ins_ptr)
		return;

	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return;

	ins->UpdateNavInstance();
}

void* nav_temlate_create(const char* file_path,
	float agentHeight, float agentRadius, float agentMaxClimb, float agentMaxSlope)
{
	RecastTemplate * tmp = new RecastTemplate;
	tmp->reset_agent(agentHeight, agentRadius, agentMaxClimb, agentMaxSlope);
	tmp->LoadTemplate(file_path);

	return tmp;
}

bool nav_load_from_template(void* ins_ptr, void* template_ptr)
{
	RecastInstance* ins = static_cast<RecastInstance*>(ins_ptr);
	if (NULL == ins)
		return false;
	RecastTemplate* tmp = static_cast<RecastTemplate*>(template_ptr);
	if (NULL == tmp)
		return false;

	ins->LoadFromTemplate(tmp->GetGeom(), tmp->GetNavTemplateMem());
	return true;
}

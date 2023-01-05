#pragma once

#include "RecastCommon.h"
#include "dllexport.h"
#include "RecastInstance.h"

class DLL_EXPORT RecastTemplate
{
public:

	RecastTemplate();
	~RecastTemplate();

	void reset_agent(float agentHeight = 2.0f, float agentRadius = 0.6f, float agentMaxClimb = 0.9f, float agentMaxSlope = 45.0f);

	bool LoadTemplate(const std::string objFilePath);

	InputGeom* GetGeom();

	const NavTemplateMem& GetNavTemplateMem();

private:

	InputGeom * m_geom;
	BuildContext m_ctx;

	NavTemplateMem m_NavTemplateMem;
	RecastInstance navIns;

	float m_agentHeight = 2.0f;
	float m_agentRadius = 0.6f;
	float m_agentMaxClimb = 0.9f;
	float m_agentMaxSlope = 45.0f;
};

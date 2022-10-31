#pragma once

#include "RecastCommon.h"
#include "dllexport.h"
#include "RecastInstance.h"

class DLL_EXPORT RecastTemplate
{
public:

	RecastTemplate();
	~RecastTemplate();

	void reset_agent(float agentHeight = 2.0f, float agentRadius = 0.6f, float agentMaxClimb = 0.9f, float agentMaxSlope = 45.0f)
	{
		m_agentHeight = agentHeight;
		m_agentRadius = agentRadius;
		m_agentMaxClimb = agentMaxClimb;
		m_agentMaxSlope = agentMaxSlope;
	}

	bool LoadTemplate(const std::string objFilePath);

	InputGeom* GetGeom()
	{
		return m_geom;
	}

	const NavTemplateMem& GetNavTemplateMem()
	{
		return m_NavTemplateMem;
	}

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

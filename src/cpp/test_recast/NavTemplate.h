#pragma once

#include "NavCommon.h"

class NavTemplate
{
public:

	NavTemplate();
	~NavTemplate();

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
};

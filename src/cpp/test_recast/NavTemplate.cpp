#include "NavTemplate.h"
#include "NavCommon.h"
#include "NavInstance.h"

NavTemplate::NavTemplate()
{}

NavTemplate::~NavTemplate()
{
	// todo free all
}

bool NavTemplate::LoadTemplate(const std::string objFilePath)
{
	m_geom = new InputGeom;
	if (!m_geom->load(&m_ctx, objFilePath))
	{
		printf("fail load:%s", objFilePath.c_str());
		delete m_geom;
		m_geom = NULL;
		return false;
	}

	if (!m_geom || !m_geom->getMesh())
	{
		m_ctx.log(RC_LOG_ERROR, "buildTiledNavigation: No vertices and triangles.");
		return false;
	}

	NavInstance navIns;
	navIns.buildFromGeom(m_geom);
	navIns.SaveToTemplate(m_NavTemplateMem);

	//{
	//	float startpos[3] = { 40.5650635f, -1.71816540f, 22.0546188f };
	//	float endpos[3] = { 49.6740074f, -2.50520134f, -6.56286621f };
	//	navIns.FindPath(startpos, endpos);
	//}

	return true;
}

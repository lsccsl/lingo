#pragma once

#include <string>
#include <memory>
#include <vector>
#include "DetourNavMeshQuery.h"
#include "RecastCommon.h"
#include "dllexport.h"
#include "RecastWrapperDef.h"

class DLL_EXPORT RecastInstance
{
	friend class RecastTemplate;
public:

	RecastInstance();
	~RecastInstance();

	void UpdateNavInstance(const int max_count = 100);

	bool LoadFromObj(const std::string objFilePath);

	void LoadFromTemplate(InputGeom* geom, const NavTemplateMem& tmpMem);
	void SaveToTemplate(NavTemplateMem& tmpMem);

	void FindPath(const float startPos[3], const float endPos[3],
		bool bprint = true);
	void FindPath(const float startPos[3], const float endPos[3],
		std::vector<RecastVec3f>& vPos,
		bool bprint = true);

	unsigned int add_obstacle(RecastVec3f* center, RecastVec3f* halfExtents);
	unsigned int add_obstacle_with_y_rotation(RecastVec3f* center, RecastVec3f* halfExtents, float rotation_y = 0.0f);
	void del_obstacle(int idObj);
	bool IsWalkAble(float PosX, float PosY, float PosZ);

	void reset_agent(float agentHeight = 2.0f, float agentRadius = 0.6f, float agentMaxClimb = 0.9f, float agentMaxSlope = 45.0f)
	{
		m_agentHeight = agentHeight;
		m_agentRadius = agentRadius;
		m_agentMaxClimb = agentMaxClimb;
		m_agentMaxSlope = agentMaxSlope;
	}

	void GetNavBound(RecastVec3f* outBoundMin, RecastVec3f* outBoundMax)
	{
		if (!outBoundMin || !outBoundMax)
			return;

		outBoundMin->x = this->nav_bound_min[0];
		outBoundMin->y = this->nav_bound_min[1];
		outBoundMin->z = this->nav_bound_min[2];
		outBoundMax->x = this->nav_bound_max[0];
		outBoundMax->y = this->nav_bound_max[1];
		outBoundMax->z = this->nav_bound_max[2];
	}

protected:

	bool buildFromGeom(InputGeom* geom);

private:

	void resetCommonSettings();
	void initSettings();

	int rasterizeTileLayers(const int tx, const int ty, const rcConfig& cfg, struct TileCacheData* tiles, const int maxTiles);

private:
	bool m_keepInterResults;

	InputGeom* m_geom;
	FastLZCompressor* m_tcomp;
	BuildContext* m_ctx;
	MeshProcess* m_tmproc;

	LinearAllocator* m_talloc;
	dtTileCache* m_tileCache;
	dtNavMesh* m_navMesh;
	dtNavMeshQuery* m_navQuery;


	// settting
	float m_cellSize;
	float m_cellHeight;
	float m_agentHeight;
	float m_agentRadius;
	float m_agentMaxClimb;
	float m_agentMaxSlope;
	float m_regionMinSize;
	float m_regionMergeSize;
	float m_edgeMaxLen;
	float m_edgeMaxError;
	float m_vertsPerPoly;
	float m_detailSampleDist;
	float m_detailSampleMaxError;
	int m_partitionType;

	int m_maxTiles;
	int m_maxPolysPerTile;
	float m_tileSize;

	bool m_filterLowHangingObstacles;
	bool m_filterLedgeSpans;
	bool m_filterWalkableLowHeightSpans;


	// find path
	float m_polyPickExt[3];
	float m_spos[3];
	float m_epos[3];
	dtQueryFilter m_filter;
	dtPolyRef m_startRef;
	dtPolyRef m_endRef;
	dtPolyRef m_polys[MAX_POLYS];
	int m_npolys;
	int m_nstraightPath;
	float m_straightPath[MAX_POLYS * 3];
	unsigned char m_straightPathFlags[MAX_POLYS];
	dtPolyRef m_straightPathPolys[MAX_POLYS];
	int m_straightPathOptions;

	float nav_bound_max[3];
	float nav_bound_min[3];
};

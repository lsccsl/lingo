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

	void UpdateNavInstance();

	bool LoadFromObj(const std::string objFilePath);

	void LoadFromTemplate(InputGeom* geom, const NavTemplateMem& tmpMem);
	void SaveToTemplate(NavTemplateMem& tmpMem);

	void FindPath(const float startPos[3], const float endPos[3],
		bool bprint = true);
	void FindPath(const float startPos[3], const float endPos[3],
		std::vector<RecastPos>& vPos,
		bool bprint = true);

	int AddBlockObject(const float posCenter[3], float sizeX, float sizeY, float sizeZ);
	void DelBlockObject(int idObj);
	bool IsWalkAble(float PosX, float PosY, float PosZ);

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
};

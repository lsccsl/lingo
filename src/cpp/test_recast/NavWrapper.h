#pragma once

#include <string>
#include <memory>
#include "DetourNavMeshQuery.h"
#include "NavCommon.h"

class NavWrapper
{
public:

	NavWrapper();
	~NavWrapper();


	bool buildFromObj(const std::string& objFilePath);

	void loadFromBin(const std::string& binFilePath);
	void saveToBin(const std::string& binFilePath);

	void FindPath(const float startPos[3], const float endPos[3]);

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

	float m_cacheBuildTimeMs;
	int m_cacheCompressedSize;
	int m_cacheRawSize;
	int m_cacheLayerCount;
	unsigned int m_cacheBuildMemUsage;

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

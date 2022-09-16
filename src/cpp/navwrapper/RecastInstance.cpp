#include "RecastInstance.h"

RecastInstance::RecastInstance() :
	m_keepInterResults(false),
	m_tileCache(0),
	m_maxTiles(0),
	m_maxPolysPerTile(0),
	m_tileSize(48),

	m_geom(0),
	m_navMesh(0),
	m_navQuery(0),
	m_filterLowHangingObstacles(true),
	m_filterLedgeSpans(true),
	m_filterWalkableLowHeightSpans(true),
	m_ctx(0),

	m_straightPathOptions(0),
	m_startRef(0),
	m_endRef(0),
	m_npolys(0),
	m_nstraightPath(0)
{
	m_navQuery = dtAllocNavMeshQuery();
	resetCommonSettings();

	m_talloc = new LinearAllocator(32000);
	m_tcomp = new FastLZCompressor;
	m_tmproc = new MeshProcess;
	m_ctx = new BuildContext;

	m_filter.setIncludeFlags(SAMPLE_POLYFLAGS_ALL ^ SAMPLE_POLYFLAGS_DISABLED);
	m_filter.setExcludeFlags(0);

	//m_polyPickExt[0] = 2;
	//m_polyPickExt[1] = 4;
	//m_polyPickExt[2] = 2;
	m_polyPickExt[0] = 1;
	m_polyPickExt[1] = 1;
	m_polyPickExt[2] = 1;
}

RecastInstance::~RecastInstance()
{
	// todo free
	if(m_tcomp)
		delete m_tcomp;
	if(m_ctx)
		delete m_ctx;
	if(m_tmproc)
		delete m_tmproc;

	if (m_navMesh)
		dtFreeNavMesh(m_navMesh);
	if(m_tileCache)
		dtFreeTileCache(m_tileCache);
	if(m_navQuery)
		dtFreeNavMeshQuery(m_navQuery);

	if(m_talloc)
		delete m_talloc;
}

void RecastInstance::resetCommonSettings()
{
	m_cellSize = 0.3f;
	m_cellHeight = 0.2f;
	m_agentHeight = 2.0f;
	m_agentRadius = 0.6f;
	m_agentMaxClimb = 0.9f;
	m_agentMaxSlope = 45.0f;
	m_regionMinSize = 8;
	m_regionMergeSize = 20;
	m_edgeMaxLen = 12.0f;
	m_edgeMaxError = 1.3f;
	m_vertsPerPoly = 6.0f;
	m_detailSampleDist = 6.0f;
	m_detailSampleMaxError = 1.0f;
	m_partitionType = SAMPLE_PARTITION_WATERSHED;
}

void RecastInstance::initSettings()
{
	int gridSize = 1;
	if (m_geom)
	{
		const float* bmin = m_geom->getNavMeshBoundsMin();
		const float* bmax = m_geom->getNavMeshBoundsMax();

		int gw = 0, gh = 0;
		rcCalcGridSize(bmin, bmax, m_cellSize, &gw, &gh);
		const int ts = (int)m_tileSize;
		const int tw = (gw + ts - 1) / ts;
		const int th = (gh + ts - 1) / ts;

		int tileBits = rcMin((int)dtIlog2(dtNextPow2(tw * th * EXPECTED_LAYERS_PER_TILE)), 14);
		if (tileBits > 14) tileBits = 14;
		int polyBits = 22 - tileBits;
		m_maxTiles = 1 << tileBits;
		m_maxPolysPerTile = 1 << polyBits;

		gridSize = tw * th;
	}
	else
	{
		m_maxTiles = 0;
		m_maxPolysPerTile = 0;
	}
}

int RecastInstance::rasterizeTileLayers(
	const int tx, const int ty,
	const rcConfig& cfg,
	TileCacheData* tiles,
	const int maxTiles)
{
	if (!m_geom || !m_geom->getMesh() || !m_geom->getChunkyMesh())
	{
		m_ctx->log(RC_LOG_ERROR, "buildTile: Input mesh is not specified.");
		return 0;
	}

	FastLZCompressor comp;
	RasterizationContext rc;

	const float* verts = m_geom->getMesh()->getVerts();
	const int nverts = m_geom->getMesh()->getVertCount();
	const rcChunkyTriMesh* chunkyMesh = m_geom->getChunkyMesh();

	// Tile bounds.
	const float tcs = cfg.tileSize * cfg.cs;

	rcConfig tcfg;
	memcpy(&tcfg, &cfg, sizeof(tcfg));

	tcfg.bmin[0] = cfg.bmin[0] + tx * tcs;
	tcfg.bmin[1] = cfg.bmin[1];
	tcfg.bmin[2] = cfg.bmin[2] + ty * tcs;
	tcfg.bmax[0] = cfg.bmin[0] + (tx + 1) * tcs;
	tcfg.bmax[1] = cfg.bmax[1];
	tcfg.bmax[2] = cfg.bmin[2] + (ty + 1) * tcs;
	tcfg.bmin[0] -= tcfg.borderSize * tcfg.cs;
	tcfg.bmin[2] -= tcfg.borderSize * tcfg.cs;
	tcfg.bmax[0] += tcfg.borderSize * tcfg.cs;
	tcfg.bmax[2] += tcfg.borderSize * tcfg.cs;

	// Allocate voxel heightfield where we rasterize our input data to.
	rc.solid = rcAllocHeightfield();
	if (!rc.solid)
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Out of memory 'solid'.");
		return 0;
	}
	if (!rcCreateHeightfield(m_ctx, *rc.solid, tcfg.width, tcfg.height, tcfg.bmin, tcfg.bmax, tcfg.cs, tcfg.ch))
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Could not create solid heightfield.");
		return 0;
	}

	// Allocate array that can hold triangle flags.
	// If you have multiple meshes you need to process, allocate
	// and array which can hold the max number of triangles you need to process.
	rc.triareas = new unsigned char[chunkyMesh->maxTrisPerChunk];
	if (!rc.triareas)
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Out of memory 'm_triareas' (%d).", chunkyMesh->maxTrisPerChunk);
		return 0;
	}

	float tbmin[2], tbmax[2];
	tbmin[0] = tcfg.bmin[0];
	tbmin[1] = tcfg.bmin[2];
	tbmax[0] = tcfg.bmax[0];
	tbmax[1] = tcfg.bmax[2];
	int cid[512];// TODO: Make grow when returning too many items.
	const int ncid = rcGetChunksOverlappingRect(chunkyMesh, tbmin, tbmax, cid, 512);
	if (!ncid)
	{
		return 0; // empty
	}

	for (int i = 0; i < ncid; ++i)
	{
		const rcChunkyTriMeshNode& node = chunkyMesh->nodes[cid[i]];
		const int* tris = &chunkyMesh->tris[node.i * 3];
		const int ntris = node.n;

		memset(rc.triareas, 0, ntris * sizeof(unsigned char));
		rcMarkWalkableTriangles(m_ctx, tcfg.walkableSlopeAngle,
			verts, nverts, tris, ntris, rc.triareas);

		if (!rcRasterizeTriangles(m_ctx, verts, nverts, tris, rc.triareas, ntris, *rc.solid, tcfg.walkableClimb))
			return 0;
	}

	// Once all geometry is rasterized, we do initial pass of filtering to
	// remove unwanted overhangs caused by the conservative rasterization
	// as well as filter spans where the character cannot possibly stand.
	if (m_filterLowHangingObstacles)
		rcFilterLowHangingWalkableObstacles(m_ctx, tcfg.walkableClimb, *rc.solid);
	if (m_filterLedgeSpans)
		rcFilterLedgeSpans(m_ctx, tcfg.walkableHeight, tcfg.walkableClimb, *rc.solid);
	if (m_filterWalkableLowHeightSpans)
		rcFilterWalkableLowHeightSpans(m_ctx, tcfg.walkableHeight, *rc.solid);


	rc.chf = rcAllocCompactHeightfield();
	if (!rc.chf)
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Out of memory 'chf'.");
		return 0;
	}
	if (!rcBuildCompactHeightfield(m_ctx, tcfg.walkableHeight, tcfg.walkableClimb, *rc.solid, *rc.chf))
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Could not build compact data.");
		return 0;
	}

	// Erode the walkable area by agent radius.
	if (!rcErodeWalkableArea(m_ctx, tcfg.walkableRadius, *rc.chf))
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Could not erode.");
		return 0;
	}

	// (Optional) Mark areas.
	const ConvexVolume* vols = m_geom->getConvexVolumes();
	for (int i = 0; i < m_geom->getConvexVolumeCount(); ++i)
	{
		rcMarkConvexPolyArea(m_ctx, vols[i].verts, vols[i].nverts,
			vols[i].hmin, vols[i].hmax,
			(unsigned char)vols[i].area, *rc.chf);
	}

	rc.lset = rcAllocHeightfieldLayerSet();
	if (!rc.lset)
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Out of memory 'lset'.");
		return 0;
	}
	if (!rcBuildHeightfieldLayers(m_ctx, *rc.chf, tcfg.borderSize, tcfg.walkableHeight, *rc.lset))
	{
		m_ctx->log(RC_LOG_ERROR, "buildNavigation: Could not build heighfield layers.");
		return 0;
	}

	rc.ntiles = 0;
	for (int i = 0; i < rcMin(rc.lset->nlayers, MAX_LAYERS); ++i)
	{
		TileCacheData* tile = &rc.tiles[rc.ntiles++];
		const rcHeightfieldLayer* layer = &rc.lset->layers[i];

		// Store header
		dtTileCacheLayerHeader header;
		header.magic = DT_TILECACHE_MAGIC;
		header.version = DT_TILECACHE_VERSION;

		// Tile layer location in the navmesh.
		header.tx = tx;
		header.ty = ty;
		header.tlayer = i;
		dtVcopy(header.bmin, layer->bmin);
		dtVcopy(header.bmax, layer->bmax);

		// Tile info.
		header.width = (unsigned char)layer->width;
		header.height = (unsigned char)layer->height;
		header.minx = (unsigned char)layer->minx;
		header.maxx = (unsigned char)layer->maxx;
		header.miny = (unsigned char)layer->miny;
		header.maxy = (unsigned char)layer->maxy;
		header.hmin = (unsigned short)layer->hmin;
		header.hmax = (unsigned short)layer->hmax;

		dtStatus status = dtBuildTileCacheLayer(&comp, &header, layer->heights, layer->areas, layer->cons,
			&tile->data, &tile->dataSize);
		if (dtStatusFailed(status))
		{
			return 0;
		}
	}

	// Transfer ownsership of tile data from build context to the caller.
	int n = 0;
	for (int i = 0; i < rcMin(rc.ntiles, maxTiles); ++i)
	{
		tiles[n++] = rc.tiles[i];
		rc.tiles[i].data = 0;
		rc.tiles[i].dataSize = 0;
	}

	return n;
}

static int calcLayerBufferSize(const int gridWidth, const int gridHeight)
{
	const int headerSize = dtAlign4(sizeof(dtTileCacheLayerHeader));
	const int gridSize = gridWidth * gridHeight;
	return headerSize + gridSize * 4;
}

bool RecastInstance::buildFromGeom(InputGeom* geom)
{
	m_geom = geom;
	dtStatus status;

	if (!m_geom || !m_geom->getMesh())
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: No vertices and triangles.");
		return false;
	}

	const BuildSettings* buildSettings = m_geom->getBuildSettings();
	if (buildSettings)
	{
		m_cellSize = buildSettings->cellSize;
		m_cellHeight = buildSettings->cellHeight;
		m_agentHeight = buildSettings->agentHeight;
		m_agentRadius = buildSettings->agentRadius;
		m_agentMaxClimb = buildSettings->agentMaxClimb;
		m_agentMaxSlope = buildSettings->agentMaxSlope;
		m_regionMinSize = buildSettings->regionMinSize;
		m_regionMergeSize = buildSettings->regionMergeSize;
		m_edgeMaxLen = buildSettings->edgeMaxLen;
		m_edgeMaxError = buildSettings->edgeMaxError;
		m_vertsPerPoly = buildSettings->vertsPerPoly;
		m_detailSampleDist = buildSettings->detailSampleDist;
		m_detailSampleMaxError = buildSettings->detailSampleMaxError;
		m_partitionType = buildSettings->partitionType;
	}
	dtFreeTileCache(m_tileCache);
	m_tileCache = 0;
	dtFreeNavMesh(m_navMesh);
	m_navMesh = 0;

	initSettings();

	m_tmproc->init(m_geom);

	// Init cache
	const float* bmin = m_geom->getNavMeshBoundsMin();
	const float* bmax = m_geom->getNavMeshBoundsMax();
	int gw = 0, gh = 0;
	rcCalcGridSize(bmin, bmax, m_cellSize, &gw, &gh);
	const int ts = (int)m_tileSize;
	const int tw = (gw + ts - 1) / ts;
	const int th = (gh + ts - 1) / ts;

	// Generation params.
	rcConfig cfg;
	memset(&cfg, 0, sizeof(cfg));
	cfg.cs = m_cellSize;
	cfg.ch = m_cellHeight;
	cfg.walkableSlopeAngle = m_agentMaxSlope;
	cfg.walkableHeight = (int)ceilf(m_agentHeight / cfg.ch);
	cfg.walkableClimb = (int)floorf(m_agentMaxClimb / cfg.ch);
	cfg.walkableRadius = (int)ceilf(m_agentRadius / cfg.cs);
	cfg.maxEdgeLen = (int)(m_edgeMaxLen / m_cellSize);
	cfg.maxSimplificationError = m_edgeMaxError;
	cfg.minRegionArea = (int)rcSqr(m_regionMinSize);		// Note: area = size*size
	cfg.mergeRegionArea = (int)rcSqr(m_regionMergeSize);	// Note: area = size*size
	cfg.maxVertsPerPoly = (int)m_vertsPerPoly;
	cfg.tileSize = (int)m_tileSize;
	cfg.borderSize = cfg.walkableRadius + 3; // Reserve enough padding.
	cfg.width = cfg.tileSize + cfg.borderSize * 2;
	cfg.height = cfg.tileSize + cfg.borderSize * 2;
	cfg.detailSampleDist = m_detailSampleDist < 0.9f ? 0 : m_cellSize * m_detailSampleDist;
	cfg.detailSampleMaxError = m_cellHeight * m_detailSampleMaxError;
	rcVcopy(cfg.bmin, bmin);
	rcVcopy(cfg.bmax, bmax);

	// Tile cache params.
	dtTileCacheParams tcparams;
	memset(&tcparams, 0, sizeof(tcparams));
	rcVcopy(tcparams.orig, bmin);
	tcparams.cs = m_cellSize;
	tcparams.ch = m_cellHeight;
	tcparams.width = (int)m_tileSize;
	tcparams.height = (int)m_tileSize;
	tcparams.walkableHeight = m_agentHeight;
	tcparams.walkableRadius = m_agentRadius;
	tcparams.walkableClimb = m_agentMaxClimb;
	tcparams.maxSimplificationError = m_edgeMaxError;
	tcparams.maxTiles = tw * th * EXPECTED_LAYERS_PER_TILE;
	tcparams.maxObstacles = 128;

	dtFreeTileCache(m_tileCache);

	m_tileCache = dtAllocTileCache();
	if (!m_tileCache)
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not allocate tile cache.");
		return false;
	}
	status = m_tileCache->init(&tcparams, m_talloc, m_tcomp, m_tmproc);
	if (dtStatusFailed(status))
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not init tile cache.");
		return false;
	}

	dtFreeNavMesh(m_navMesh);

	m_navMesh = dtAllocNavMesh();
	if (!m_navMesh)
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not allocate navmesh.");
		return false;
	}

	dtNavMeshParams params;
	memset(&params, 0, sizeof(params));
	rcVcopy(params.orig, bmin);
	params.tileWidth = m_tileSize * m_cellSize;
	params.tileHeight = m_tileSize * m_cellSize;
	params.maxTiles = m_maxTiles;
	params.maxPolys = m_maxPolysPerTile;

	status = m_navMesh->init(&params);
	if (dtStatusFailed(status))
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not init navmesh.");
		return false;
	}

	status = m_navQuery->init(m_navMesh, 2048);
	if (dtStatusFailed(status))
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not init Detour navmesh query");
		return false;
	}

	// Preprocess tiles.

	m_ctx->resetTimers();

	for (int y = 0; y < th; ++y)
	{
		for (int x = 0; x < tw; ++x)
		{
			TileCacheData tiles[MAX_LAYERS];
			memset(tiles, 0, sizeof(tiles));
			int ntiles = rasterizeTileLayers(x, y, cfg, tiles, MAX_LAYERS);

			for (int i = 0; i < ntiles; ++i)
			{
				TileCacheData* tile = &tiles[i];
				status = m_tileCache->addTile(tile->data, tile->dataSize, DT_COMPRESSEDTILE_FREE_DATA, 0);
				if (dtStatusFailed(status))
				{
					dtFree(tile->data);
					tile->data = 0;
					continue;
				}
			}
		}
	}

	// Build initial meshes
	m_ctx->startTimer(RC_TIMER_TOTAL);
	for (int y = 0; y < th; ++y)
		for (int x = 0; x < tw; ++x)
			m_tileCache->buildNavMeshTilesAt(x, y, m_navMesh);
	m_ctx->stopTimer(RC_TIMER_TOTAL);

	const dtNavMesh* nav = m_navMesh;
	int navmeshMemUsage = 0;
	for (int i = 0; i < nav->getMaxTiles(); ++i)
	{
		const dtMeshTile* tile = nav->getTile(i);
		if (tile->header)
			navmeshMemUsage += tile->dataSize;
	}
	printf("navmeshMemUsage = %.1f kB", navmeshMemUsage / 1024.0f);

	return true;

}

void RecastInstance::FindPath(const float startPos[3], const float endPos[3], bool bprint)
{
	if (!m_navMesh)
		return;

	dtVcopy(m_spos, startPos);
	dtVcopy(m_epos, endPos);

	m_navQuery->findNearestPoly(m_spos, m_polyPickExt, &m_filter, &m_startRef, 0);
	m_navQuery->findNearestPoly(m_epos, m_polyPickExt, &m_filter, &m_endRef, 0);

	//m_pathFindStatus = DT_FAILURE;

#ifdef DUMP_REQS
	if (bprint)
	{
		printf("ps  %f %f %f  %f %f %f  0x%x 0x%x\n",
			m_spos[0], m_spos[1], m_spos[2], m_epos[0], m_epos[1], m_epos[2],
			m_filter.getIncludeFlags(), m_filter.getExcludeFlags());
}
#endif
	m_navQuery->findPath(m_startRef, m_endRef, m_spos, m_epos, &m_filter, m_polys, &m_npolys, MAX_POLYS);//lsc 按网格邻接查找网格路径,类A*算法
	m_nstraightPath = 0;
	if (m_npolys)
	{
		// In case of partial path, make sure the end point is clamped to the last polygon.
		float epos[3];
		dtVcopy(epos, m_epos);
		if (m_polys[m_npolys - 1] != m_endRef)
			m_navQuery->closestPointOnPoly(m_polys[m_npolys - 1], m_epos, epos, 0);//lsc 目标点不可达的情况

		m_navQuery->findStraightPath(m_spos, epos, m_polys, m_npolys,
			m_straightPath, m_straightPathFlags,
			m_straightPathPolys, &m_nstraightPath, MAX_POLYS, m_straightPathOptions);//lsc 寻找拐点,生成路径点

#if 1
		if (bprint)
		{
			printf("\r\n\r\nall grid:\r\n");
			for (int j = 0; j < m_npolys; j++)
			{
				printf("%d ", m_polys[j]);
			}
			printf("\r\n");
			printf("all point:\r\n");
			for (int i = 0; i < m_nstraightPath; i++)
			{
				printf("{%f,%f,%f}grid:%d flag:%02x\r\n",
					m_straightPath[i * 3], m_straightPath[i * 3 + 1], m_straightPath[i * 3 + 2],
					m_straightPathPolys[i],
					m_straightPathFlags[i]);
			}
		}
#endif
	}
}

void RecastInstance::LoadFromTemplate(InputGeom* geom, const NavTemplateMem& tmpMem)
{
	m_geom = geom;

	if (!m_geom || !m_geom->getMesh())
		return;

	const BuildSettings* buildSettings = m_geom->getBuildSettings();
	if (buildSettings)
	{
		m_cellSize = buildSettings->cellSize;
		m_cellHeight = buildSettings->cellHeight;
		m_agentHeight = buildSettings->agentHeight;
		m_agentRadius = buildSettings->agentRadius;
		m_agentMaxClimb = buildSettings->agentMaxClimb;
		m_agentMaxSlope = buildSettings->agentMaxSlope;
		m_regionMinSize = buildSettings->regionMinSize;
		m_regionMergeSize = buildSettings->regionMergeSize;
		m_edgeMaxLen = buildSettings->edgeMaxLen;
		m_edgeMaxError = buildSettings->edgeMaxError;
		m_vertsPerPoly = buildSettings->vertsPerPoly;
		m_detailSampleDist = buildSettings->detailSampleDist;
		m_detailSampleMaxError = buildSettings->detailSampleMaxError;
		m_partitionType = buildSettings->partitionType;
	}
	dtFreeTileCache(m_tileCache);
	m_tileCache = 0;
	dtFreeNavMesh(m_navMesh);
	m_navMesh = 0;

	initSettings();

	m_tmproc->init(m_geom);

	// Read header.
	const TileCacheSetHeader& header = tmpMem.m_TileCacheSetHeader;
	if (header.magic != TILECACHESET_MAGIC)
		return;

	if (header.version != TILECACHESET_VERSION)
		return;

	m_navMesh = dtAllocNavMesh();
	if (!m_navMesh)
		return;

	dtStatus status = m_navMesh->init(&header.meshParams);
	if (dtStatusFailed(status))
		return;

	status = m_navQuery->init(m_navMesh, 2048);
	if (dtStatusFailed(status))
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: Could not init Detour navmesh query");
		return;
	}

	m_tileCache = dtAllocTileCache();
	if (!m_tileCache)
		return;

	status = m_tileCache->init(&header.cacheParams, m_talloc, m_tcomp, m_tmproc);
	if (dtStatusFailed(status))
		return;

	// Read tiles.
	for (int i = 0; i < header.numTiles && i < tmpMem.m_vNavTemplateTile.size(); ++i)
	{
		const NavTemplateMem::NavTemplateTile& navTile = tmpMem.m_vNavTemplateTile[i];
		const TileCacheTileHeader& tileHeader = navTile.m_TileCacheTileHeader;

		if (!tileHeader.tileRef || !tileHeader.dataSize)
			break;

		unsigned char* data = (unsigned char*)dtAlloc(navTile.m_vData.size(), DT_ALLOC_PERM);
		if (!data)
			break;
		memcpy(data, &navTile.m_vData[0], navTile.m_vData.size());

		dtCompressedTileRef tile = 0;
		dtStatus addTileStatus = m_tileCache->addTile(data, tileHeader.dataSize, DT_COMPRESSEDTILE_FREE_DATA, &tile);
		if (dtStatusFailed(addTileStatus))
		{
			dtFree(data);
			break;
		}

		if (tile)
			m_tileCache->buildNavMeshTile(tile, m_navMesh);
	}
}

void RecastInstance::SaveToTemplate(NavTemplateMem& tmpMem)
{
	if (!m_tileCache) return;

	// Store header.
	tmpMem.m_TileCacheSetHeader.magic = TILECACHESET_MAGIC;
	tmpMem.m_TileCacheSetHeader.version = TILECACHESET_VERSION;
	tmpMem.m_TileCacheSetHeader.numTiles = 0;
	for (int i = 0; i < m_tileCache->getTileCount(); ++i)
	{
		const dtCompressedTile* tile = m_tileCache->getTile(i);
		if (!tile || !tile->header || !tile->dataSize) continue;
		tmpMem.m_TileCacheSetHeader.numTiles++;
	}
	memcpy(&tmpMem.m_TileCacheSetHeader.cacheParams, m_tileCache->getParams(), sizeof(dtTileCacheParams));
	memcpy(&tmpMem.m_TileCacheSetHeader.meshParams, m_navMesh->getParams(), sizeof(dtNavMeshParams));

	tmpMem.m_vNavTemplateTile.resize(m_tileCache->getTileCount());
	// Store tiles.
	for (int i = 0; i < m_tileCache->getTileCount() && i < tmpMem.m_vNavTemplateTile.size(); ++i)
	{
		const dtCompressedTile* tile = m_tileCache->getTile(i);
		if (!tile || !tile->header || !tile->dataSize) continue;

		NavTemplateMem::NavTemplateTile& tmpTile = tmpMem.m_vNavTemplateTile[i];
		tmpTile.m_TileCacheTileHeader.tileRef = m_tileCache->getTileRef(tile);
		tmpTile.m_TileCacheTileHeader.dataSize = tile->dataSize;

		tmpTile.m_vData.resize(tile->dataSize);
		memcpy(&tmpTile.m_vData[0], tile->data, tile->dataSize);
	}
}

int RecastInstance::AddBlockObject(const float posCenter[3], float sizeX, float sizeY, float sizeZ)
{
	if (!m_tileCache)
		return -1;

	float bmin[3];
	float bmax[3];
	const float ext[3] = { sizeX, sizeY, sizeZ };//{ 1.0f, 1.5f, 2.0f };
	dtVadd(bmax, posCenter, ext);
	dtVsub(bmin, posCenter, ext);
	dtObstacleRef result;
	dtStatus ret = m_tileCache->addBoxObstacle(bmin, bmax, &result);
	if (DT_SUCCESS != ret)
		return -1;

	return result;
}

void RecastInstance::DelBlockObject(int idObj)
{
	if (!m_tileCache)
		return;
	m_tileCache->removeObstacle(idObj);
}

bool RecastInstance::IsWalkAble(float PosX, float PosY, float PosZ)
{
	float pos[3] = { PosX, PosY, PosZ };

	float polyPickExt[3] = { 0.01f, 4.0f, 0.01f };

	dtPolyRef nearestRef = 0;
	dtStatus ret = m_navQuery->findNearestPoly(pos, polyPickExt, &m_filter, &nearestRef, 0);
	if (DT_SUCCESS != ret)
		return false;

	if (0 == nearestRef)
		return false;

	return true;
}

void RecastInstance::UpdateNavInstance()
{
	if (!m_navMesh)
		return;
	if (!m_tileCache)
		return;

	m_tileCache->update(1, m_navMesh);
}

bool RecastInstance::LoadFromObj(const std::string objFilePath)
{
	m_geom = new InputGeom;
	if (!m_geom->load(m_ctx, objFilePath))
	{
		printf("fail load:%s", objFilePath.c_str());
		delete m_geom;
		m_geom = NULL;
		return false;
	}

	if (!m_geom || !m_geom->getMesh())
	{
		m_ctx->log(RC_LOG_ERROR, "buildTiledNavigation: No vertices and triangles.");
		return false;
	}

	return this->buildFromGeom(m_geom);
}

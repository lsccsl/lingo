#include "RecastWrapper.h"

#include "RecastTemplate.h"
#include "RecastInstance.h"
#include "RecastCommon.h"
#include <Windows.h>

void test_1()
{
	printf("test_1\r\n");
	RecastWrapper nw;
	nw.buildFromObj("./test_mesh/nav_test.obj");

	float startpos[3] = { 40.5650635f, -1.71816540f, 22.0546188f };
	float endpos[3] = { 49.6740074f, -2.50520134f, -6.56286621f };
	nw.FindPath(startpos, endpos);

	nw.saveToBin("test.bin");
	nw.loadFromBin("test.bin");
	nw.FindPath(startpos, endpos);
}

void test_template()
{
	printf("test_template\r\n");
	RecastTemplate navTmp;
	navTmp.LoadTemplate("./test_mesh/nav_test.obj");

	float startpos[3] = { 40.5650635f, -1.71816540f, 22.0546188f };
	float endpos[3] = { 49.6740074f, -2.50520134f, -6.56286621f };

	RecastInstance navIns_1;
	{
		printf("\r\ninstance 1 find path\r\n");
		navIns_1.LoadFromTemplate(navTmp.GetGeom(), navTmp.GetNavTemplateMem());
		navIns_1.FindPath(startpos, endpos);
	}

	RecastInstance navIns_2;
	int objObstalce = 0;
	{
		auto tStart = ::getPerfTime();
		navIns_2.LoadFromTemplate(navTmp.GetGeom(), navTmp.GetNavTemplateMem());
		auto tEnd = ::getPerfTime();
		printf("load :%fs %dmicroseconds\r\n", getPerfTimeUsec(tEnd - tStart)/1000000.f, getPerfTimeUsec(tEnd - tStart));

		float posCenter[3] = { 48.2378387f, -1.40648651f, 8.61733246f, };
		tStart = ::getPerfTime();
		objObstalce = navIns_2.AddBlockObject(posCenter, 2.0f, 2.0f, 2.0f);
		navIns_2.UpdateNavInstance();
		navIns_2.UpdateNavInstance();
		tEnd = ::getPerfTime();

		printf("\r\ninstance 2 add obstacle:%fs %dmicroseconds\r\n", getPerfTimeUsec(tEnd - tStart) / 1000000.f, getPerfTimeUsec(tEnd - tStart));
		tStart = ::getPerfTime();
		navIns_2.FindPath(startpos, endpos, true);
		tEnd = ::getPerfTime();
		printf("\r\ninstance 2 find path:%fs %dmicroseconds\r\n", getPerfTimeUsec(tEnd - tStart) / 1000000.f, getPerfTimeUsec(tEnd - tStart));

		bool bret = navIns_2.IsWalkAble(48.2378387f, 0/*-1.40648651f*/, 8.61733246f);
		printf("\r\n\r\n walkable:%d", bret);
	}

	{
		printf("\r\ninstance 1 find path\r\n");
		navIns_1.LoadFromTemplate(navTmp.GetGeom(), navTmp.GetNavTemplateMem());
		navIns_1.FindPath(startpos, endpos);
	}

	// walkable
	{
		bool bret = navIns_2.IsWalkAble(100.0f, 0.0f, 100.0f);
		printf("\r\n\r\n walkable:%d", bret);

		bret = navIns_2.IsWalkAble(40.5650635f, 0 /*-1.71816540f*/, 22.0546188f);
		printf("\r\n\r\n walkable:%d", bret);
		float startpos[3] = {  };

		bret = navIns_2.IsWalkAble(41.4345322f, 0 /*16.7074146*/, -9.18507099);
		printf("\r\n\r\n walkable:%d", bret);
	}

	{
		navIns_2.DelBlockObject(objObstalce);
		navIns_2.UpdateNavInstance();
		navIns_2.UpdateNavInstance();

		printf("\r\ninstance 2 find path after del obstacle\r\n");
		navIns_2.FindPath(startpos, endpos);

		bool bret = navIns_2.IsWalkAble(48.2378387f, 0, 8.61733246f);
		printf("\r\n\r\n walkable:%d", bret);
	}
}

int main()
{
	test_template();

	test_1();
}


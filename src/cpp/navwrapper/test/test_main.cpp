#include "RecastWrapper.h"

#include "RecastTemplate.h"
#include "RecastInstance.h"
#include "RecastCommon.h"
#include "RecastCWrapper.h"

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

void test_instance()
{
	printf("test instance\r\n");
	RecastInstance ins;

	ins.LoadFromObj("./test_mesh/nav_test.obj");
	float startpos[3] = { 40.5650635f, -1.71816540f, 22.0546188f };
	float endpos[3] = { 49.6740074f, -2.50520134f, -6.56286621f };
	ins.FindPath(startpos, endpos);
}

void test_cwrapper()
{
	printf("\r\n\r\n ~~~~test_cwrapper\r\n");
	//void * ins_ptr = nav_create("./test_mesh/nav_test.obj");
	void* ins_ptr = nav_new();
	nav_reset_agent(ins_ptr, 2.0f, 0.6f, 0.9f, 45.0f);
	nav_load(ins_ptr, "./test_mesh/nav_test.obj");

	RecastVec3f posCenter = { 48.2378387f, -1.40648651f, 8.61733246f, };
	RecastVec3f halfExt = RecastVec3f{ 2.0f, 2.0f, 2.0f };

	nav_add_obstacle(ins_ptr, &posCenter, &halfExt, (45.0f / 360.0f) * 2.0f * 3.14f);
	nav_update(ins_ptr);
	nav_update(ins_ptr);
	//void* ins_ptr = nav_create("../../../../resource/test_scene.obj");
	RecastVec3f startpos;
	startpos.x = 40.5650635f;//702.190918f; 
	startpos.y = -1.71816540f;// 0.0f;// 1.53082275f; 
	startpos.z = 22.0546188f; //635.378662f; 
	RecastVec3f endpos;
	endpos.x = 49.6740074f; //710.805664f;
	endpos.y = -2.50520134f; //0.0f;// 1.00000000f;
	endpos.z = -6.56286621f; //851.753296f;
	RecastVec3f* pos_path = NULL;
	int pos_path_sz = 0;
	nav_findpath(ins_ptr, &startpos, &endpos, &pos_path, &pos_path_sz, true);
	printf("\r\n");
	for (int i = 0; i < pos_path_sz; i++)
	{
		printf("{%f,%f,%f}\r\n",
			pos_path[i].x, pos_path[i].y, pos_path[i].z);
	}
	printf("end test_cwrapper\r\n\r\n");
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
		navIns_1.reset_agent();
		navIns_1.LoadFromTemplate(navTmp.GetGeom(), navTmp.GetNavTemplateMem());
		navIns_1.FindPath(startpos, endpos);
	}

	RecastInstance navIns_2;
	int objObstalce = 0;
	{
		auto tStart = ::getPerfTime();
		navIns_2.reset_agent();
		navIns_2.LoadFromTemplate(navTmp.GetGeom(), navTmp.GetNavTemplateMem());
		auto tEnd = ::getPerfTime();
		printf("load :%fs %dmicroseconds\r\n", getPerfTimeUsec(tEnd - tStart)/1000000.f, getPerfTimeUsec(tEnd - tStart));

		RecastVec3f posCenter = { 48.2378387f, -1.40648651f, 8.61733246f, };
		RecastVec3f halfExt = RecastVec3f{ 2.0f, 2.0f, 2.0f };
		tStart = ::getPerfTime();
		objObstalce = navIns_2.add_obstacle_with_y_rotation(&posCenter, &halfExt, (45.0f / 360.0f) * 2.0f * 3.14f);
		navIns_2.UpdateNavInstance();
		navIns_2.UpdateNavInstance();
		tEnd = ::getPerfTime();

		printf("\r\ninstance 2 add obstacle:%fs %dmicroseconds\r\n", getPerfTimeUsec(tEnd - tStart) / 1000000.f, getPerfTimeUsec(tEnd - tStart));
		tStart = ::getPerfTime();
		navIns_2.FindPath(startpos, endpos);
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
		navIns_2.del_obstacle(objObstalce);
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
	test_instance();

	test_1();
	test_template();

	test_cwrapper();
}


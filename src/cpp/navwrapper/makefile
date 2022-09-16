INCL=-I./ -I./Detour/Include -I./DetourTileCache/Include -I./Recast/Include -I./RecastDemo
LIBDIR=-L./ -L./bin
LIB=-lstdc++ -lpthread -ldl
OUTPUT=./
FLAG=-fPIC -std=c++11 -g3
CFLAG=-fPIC -g3


.SUFFIXES: .cpp .c .cc

.cpp.o:
	g++ ${INCL} ${FLAG} -c $< -o $@
.cc.o:
	g++ ${INCL} ${FLAG} -c $< -o $@
.c.o:
	gcc ${INCL} ${CFLAG} -c $< -o $@

SRC_Detour = $(wildcard ./Detour/Source/*.cpp)
SRC_DetourTileCache = $(wildcard ./DetourTileCache/Source/*.cpp)
SRC_Recast = $(wildcard ./Recast/Source/*.cpp)
SRC_RecastDemo = $(wildcard ./RecastDemo/*.cpp)
SRC_RecastDemoC = $(wildcard ./RecastDemo/*.c)

OBJ_Detour = $(patsubst %.cpp, %.o, $(SRC_Detour))
OBJ_DetourTileCache = $(patsubst %.cpp, %.o, $(SRC_DetourTileCache))
OBJ_Recast = $(patsubst %.cpp, %.o, $(SRC_Recast))
OBJ_RecastDemo = $(patsubst %.cpp, %.o, $(SRC_RecastDemo))
OBJ_RecastDemoC = $(patsubst %.c, %.o, $(SRC_RecastDemoC))

TARGETS = $(OBJ_Detour) $(OBJ_DetourTileCache) $(OBJ_Recast) $(OBJ_RecastDemo) $(OBJ_RecastDemoC) RecastCommon.o RecastInstance.o RecastTemplate.o RecastWrapper.o RecastCWrapper.o

all: libnavwrapper.so test_nav

libnavwrapper.so: $(TARGETS)
	gcc -o $@ -shared $? $(LIBDIR) $(LIB)
	cp libnavwrapper.so ./bin/

test_nav: ./test/test_main.o
	gcc -o $@ $? $(LIBDIR) $(LIB) -lnavwrapper
	cp test_nav ./bin/test_nav

clean:
	rm -fr $(TARGETS)
	rm -fr libnavwrapper.so
	rm -fr test_nav

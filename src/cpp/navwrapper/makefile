INCL=-I./ -I./Detour/Include -I./DetourTileCache/Include -I./Recast/Include -I./RecastDemo
LIBDIR=-L./ -L./bin
LIB=-lstdc++ -lpthread -ldl
OUTPUT=./
FLAG=-fPIC -std=c++11 -g3


.SUFFIXES: .cpp .c .cc

.cpp.o:
	g++ ${INCL} ${FLAG} -c $< -o $@
.cc.o:
	g++ ${INCL} ${FLAG} -c $< -o $@
.c.o:
	gcc ${INCL} ${FLAG} -c $< -o $@

SRC_Detour = $(wildcard ./Detour/Source/*.cpp)
SRC_DetourTileCache = $(wildcard ./DetourTileCache/Source/*.cpp)
SRC_Recast = $(wildcard ./Recast/Source/*.cpp)
SRC_RecastDemo = $(wildcard ./RecastDemo/*.cpp)

OBJ_Detour = $(patsubst %.cpp, %.o, $(SRC_Detour))
OBJ_DetourTileCache = $(patsubst %.cpp, %.o, $(SRC_DetourTileCache))
OBJ_Recast = $(patsubst %.cpp, %.o, $(SRC_Recast))
OBJ_RecastDemo = $(patsubst %.cpp, %.o, $(SRC_RecastDemo))

TARGETS = $(OBJ_Detour) $(OBJ_DetourTileCache) $(OBJ_Recast) $(OBJ_RecastDemo) RecastCommon.o RecastInstance.o RecastTemplate.o RecastWrapper.o

all: libnavwrapper.so

libnavwrapper.so: $(TARGETS)
	gcc -o $@ -shared $? $(LIBDIR) $(LIB)
	mv libnavwrapper.so ./bin/

clean:
	$(warning ==== build files: $(TARGETS) ====)
	rm -fr $(TARGETS)
	rm -fr libnavwrapper.so

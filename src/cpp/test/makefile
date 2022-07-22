INCL=-I./. -I./protobuf-3.19.2/cmake/install/usr/local/include
LIBDIR=-L./ -L../lib -L./protobuf-3.19.2/cmake
LIB=-lstdc++ -lpthread -ldl -lprotobuf
OUTPUT=./
FLAG=-fPIC -std=c++11 -g3


.SUFFIXES: .cpp .c .cc

.cpp.o:
	g++ ${INCL} ${FLAG} -c $<
.cc.o:
	g++ ${INCL} ${FLAG} -c $<
.c.o:
	gcc ${INCL} ${FLAG} -c $<


all: test_client

test_client: test.o channel.o msg.pb.o mylogex.o testclient_mgr.o os.o testclient.o
	gcc -o $@ $? $(LIBDIR) $(LIB)
	cp ./test_client ../../../output/

clean:
	rm -fr ./*.o
	rm -fr ./*.so
	rm -fr ./test_client

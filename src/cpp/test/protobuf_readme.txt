tar -zxvf protobuf-3.17.0.tar.gz

cd protobuf-3.17.0

cd cmake

cmake . -Dprotobuf_BUILD_TESTS=OFF(第一步)

cmake --build . (或者使用make)

make install DESTDIR=安装目录(也可以在第一步中通过-DCMAKE_INSTALL_PREFIX来设置安装目录)

编译安装完成后，在安装目录下，将最底层中的文件夹bin(protoc可执行文件)、include(头文件)、lib64(库文件)放入到新文件夹protobuf中。（注意，这里生成的是静态库，如何生成动态库暂时不清楚。）。
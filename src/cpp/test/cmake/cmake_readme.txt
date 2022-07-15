linux升级cmake版本
2021-10-22 10:25:36
CMake版本升级
1、在网址 https://cmake.org/files/v3.1/ 下载 cmake-3.1.2.tar.gz
2、解压
3、执行 ./configure
4、执行 make
5、执行 sudo make install
6、执行 sudo update-alternatives --install /usr/bin/cmake  cmake /usr/local/bin/cmake  1 --force
7、运行 cmake --version 查看版本号

注意： 第6步 update-alternatives 命令用于处理linux系统中软件版本的切换。
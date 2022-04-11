# Microsoft Developer Studio Project File - Name="common_cpplib" - Package Owner=<4>
# Microsoft Developer Studio Generated Build File, Format Version 6.00
# ** DO NOT EDIT **

# TARGTYPE "Win32 (x86) Static Library" 0x0104

CFG=common_cpplib - Win32 Debug
!MESSAGE This is not a valid makefile. To build this project using NMAKE,
!MESSAGE use the Export Makefile command and run
!MESSAGE 
!MESSAGE NMAKE /f "common_cpplib.mak".
!MESSAGE 
!MESSAGE You can specify a configuration when running NMAKE
!MESSAGE by defining the macro CFG on the command line. For example:
!MESSAGE 
!MESSAGE NMAKE /f "common_cpplib.mak" CFG="common_cpplib - Win32 Debug"
!MESSAGE 
!MESSAGE Possible choices for configuration are:
!MESSAGE 
!MESSAGE "common_cpplib - Win32 Release" (based on "Win32 (x86) Static Library")
!MESSAGE "common_cpplib - Win32 Debug" (based on "Win32 (x86) Static Library")
!MESSAGE 

# Begin Project
# PROP AllowPerConfigDependencies 0
# PROP Scc_ProjName ""
# PROP Scc_LocalPath ""
CPP=cl.exe
RSC=rc.exe

!IF  "$(CFG)" == "common_cpplib - Win32 Release"

# PROP BASE Use_MFC 0
# PROP BASE Use_Debug_Libraries 0
# PROP BASE Output_Dir "Release"
# PROP BASE Intermediate_Dir "Release"
# PROP BASE Target_Dir ""
# PROP Use_MFC 0
# PROP Use_Debug_Libraries 0
# PROP Output_Dir "Release"
# PROP Intermediate_Dir "Release"
# PROP Target_Dir ""
# ADD BASE CPP /nologo /W3 /GX /O2 /D "WIN32" /D "NDEBUG" /D "_MBCS" /D "_LIB" /YX /FD /c
# ADD CPP /nologo /MT /W3 /GX /O2 /I "..\..\Rhapsody-0.2.0\Rhapsody-0.2.0\include" /I "..\..\pthreads.1" /D "WIN32" /D "NDEBUG" /D "_MBCS" /D "_LIB" /YX /FD /c
# ADD BASE RSC /l 0x804 /d "NDEBUG"
# ADD RSC /l 0x804 /d "NDEBUG"
BSC32=bscmake.exe
# ADD BASE BSC32 /nologo
# ADD BSC32 /nologo
LIB32=link.exe -lib
# ADD BASE LIB32 /nologo
# ADD LIB32 /nologo

!ELSEIF  "$(CFG)" == "common_cpplib - Win32 Debug"

# PROP BASE Use_MFC 0
# PROP BASE Use_Debug_Libraries 1
# PROP BASE Output_Dir "Debug"
# PROP BASE Intermediate_Dir "Debug"
# PROP BASE Target_Dir ""
# PROP Use_MFC 0
# PROP Use_Debug_Libraries 1
# PROP Output_Dir "Debug"
# PROP Intermediate_Dir "Debug"
# PROP Target_Dir ""
# ADD BASE CPP /nologo /W3 /Gm /GX /ZI /Od /D "WIN32" /D "_DEBUG" /D "_MBCS" /D "_LIB" /YX /FD /GZ /c
# ADD CPP /nologo /MTd /W3 /Gm /GX /ZI /Od /I "..\..\Rhapsody-0.2.0\Rhapsody-0.2.0\include" /I "..\..\pthreads.1" /D "WIN32" /D "_DEBUG" /D "_MBCS" /D "_LIB" /YX /FD /GZ /c
# ADD BASE RSC /l 0x804 /d "_DEBUG"
# ADD RSC /l 0x804 /d "_DEBUG"
BSC32=bscmake.exe
# ADD BASE BSC32 /nologo
# ADD BSC32 /nologo
LIB32=link.exe -lib
# ADD BASE LIB32 /nologo
# ADD LIB32 /nologo

!ENDIF 

# Begin Target

# Name "common_cpplib - Win32 Release"
# Name "common_cpplib - Win32 Debug"
# Begin Group "Source Files"

# PROP Default_Filter "cpp;c;cxx;rc;def;r;odl;idl;hpj;bat"
# Begin Source File

SOURCE=..\AutoLock.cpp

!IF  "$(CFG)" == "common_cpplib - Win32 Release"

!ELSEIF  "$(CFG)" == "common_cpplib - Win32 Debug"

!ENDIF 

# End Source File
# Begin Source File

SOURCE=..\AutoLock.h
# End Source File
# Begin Source File

SOURCE=..\CfgFile.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\CfgFile.h
# End Source File
# Begin Source File

SOURCE=..\channel.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\channel.h
# End Source File
# Begin Source File

SOURCE=..\CMySocket.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\CMySocket.h
# End Source File
# Begin Source File

SOURCE=..\common_def.h
# End Source File
# Begin Source File

SOURCE=..\crc32.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\crc32.h
# End Source File
# Begin Source File

SOURCE=..\Md5.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\Md5.h
# End Source File
# Begin Source File

SOURCE=..\myepoll.cpp

!IF  "$(CFG)" == "common_cpplib - Win32 Release"

!ELSEIF  "$(CFG)" == "common_cpplib - Win32 Debug"

!ENDIF 

# End Source File
# Begin Source File

SOURCE=..\myepoll.h
# End Source File
# Begin Source File

SOURCE=..\myepoll_linux.h
# End Source File
# Begin Source File

SOURCE=..\myepoll_win32.h
# End Source File
# Begin Source File

SOURCE=..\mylogex.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\mylogex.h
# End Source File
# Begin Source File

SOURCE=..\myos_cpp.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\myos_cpp.h
# End Source File
# Begin Source File

SOURCE=..\mythrdpoll.cpp

!IF  "$(CFG)" == "common_cpplib - Win32 Release"

!ELSEIF  "$(CFG)" == "common_cpplib - Win32 Debug"

!ENDIF 

# End Source File
# Begin Source File

SOURCE=..\mythrdpoll.h
# End Source File
# Begin Source File

SOURCE=..\osfile_cpp.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\osfile_cpp.h
# End Source File
# Begin Source File

SOURCE=..\OsInet.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\OsInet.h
# End Source File
# Begin Source File

SOURCE=..\PacketBase.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\PacketBase.h
# End Source File
# Begin Source File

SOURCE=..\runtime_view.cpp
# SUBTRACT CPP /YX /Yc /Yu
# End Source File
# Begin Source File

SOURCE=..\runtime_view.h
# End Source File
# Begin Source File

SOURCE=..\stringOp.cpp

!IF  "$(CFG)" == "common_cpplib - Win32 Release"

!ELSEIF  "$(CFG)" == "common_cpplib - Win32 Debug"

!ENDIF 

# End Source File
# Begin Source File

SOURCE=..\stringOp.h
# End Source File
# Begin Source File

SOURCE=..\type_def.h
# End Source File
# End Group
# Begin Group "Header Files"

# PROP Default_Filter "h;hpp;hxx;hm;inl"
# End Group
# End Target
# End Project

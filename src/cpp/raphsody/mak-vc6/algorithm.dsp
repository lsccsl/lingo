# Microsoft Developer Studio Project File - Name="algorithm" - Package Owner=<4>
# Microsoft Developer Studio Generated Build File, Format Version 6.00
# ** DO NOT EDIT **

# TARGTYPE "Win32 (x86) Static Library" 0x0104

CFG=algorithm - Win32 Debug
!MESSAGE This is not a valid makefile. To build this project using NMAKE,
!MESSAGE use the Export Makefile command and run
!MESSAGE 
!MESSAGE NMAKE /f "algorithm.mak".
!MESSAGE 
!MESSAGE You can specify a configuration when running NMAKE
!MESSAGE by defining the macro CFG on the command line. For example:
!MESSAGE 
!MESSAGE NMAKE /f "algorithm.mak" CFG="algorithm - Win32 Debug"
!MESSAGE 
!MESSAGE Possible choices for configuration are:
!MESSAGE 
!MESSAGE "algorithm - Win32 Release" (based on "Win32 (x86) Static Library")
!MESSAGE "algorithm - Win32 Debug" (based on "Win32 (x86) Static Library")
!MESSAGE 

# Begin Project
# PROP AllowPerConfigDependencies 0
# PROP Scc_ProjName ""
# PROP Scc_LocalPath ""
CPP=cl.exe
RSC=rc.exe

!IF  "$(CFG)" == "algorithm - Win32 Release"

# PROP BASE Use_MFC 0
# PROP BASE Use_Debug_Libraries 0
# PROP BASE Output_Dir "Release"
# PROP BASE Intermediate_Dir "Release"
# PROP BASE Target_Dir ""
# PROP Use_MFC 0
# PROP Use_Debug_Libraries 0
# PROP Output_Dir "../bin"
# PROP Intermediate_Dir "../output/algorithm/Release"
# PROP Target_Dir ""
# ADD BASE CPP /nologo /W3 /GX /O2 /D "WIN32" /D "NDEBUG" /D "_MBCS" /D "_LIB" /YX /FD /c
# ADD CPP /nologo /MT /W3 /GX /O2 /I "..\include" /D "WIN32" /D "NDEBUG" /D "_MBCS" /D "_LIB" /YX /FD /c
# ADD BASE RSC /l 0x804 /d "NDEBUG"
# ADD RSC /l 0x804 /d "NDEBUG"
BSC32=bscmake.exe
# ADD BASE BSC32 /nologo
# ADD BSC32 /nologo
LIB32=link.exe -lib
# ADD BASE LIB32 /nologo
# ADD LIB32 /nologo
# Begin Special Build Tool
SOURCE="$(InputPath)"
PostBuild_Cmds=md       ..\include\      	copy       ..\algorithm\*.h       ..\include\ 
# End Special Build Tool

!ELSEIF  "$(CFG)" == "algorithm - Win32 Debug"

# PROP BASE Use_MFC 0
# PROP BASE Use_Debug_Libraries 1
# PROP BASE Output_Dir "algorithm___Win32_Debug"
# PROP BASE Intermediate_Dir "algorithm___Win32_Debug"
# PROP BASE Target_Dir ""
# PROP Use_MFC 0
# PROP Use_Debug_Libraries 1
# PROP Output_Dir "../bin"
# PROP Intermediate_Dir "../output/algorithm/debug"
# PROP Target_Dir ""
# ADD BASE CPP /nologo /W3 /Gm /GX /ZI /Od /D "WIN32" /D "_DEBUG" /D "_MBCS" /D "_LIB" /YX /FD /GZ /c
# ADD CPP /nologo /MTd /W3 /Gm /GX /ZI /Od /I "..\include\\" /D "WIN32" /D "_DEBUG" /D "_MBCS" /D "_LIB" /YX /FD /GZ /c
# ADD BASE RSC /l 0x804 /d "_DEBUG"
# ADD RSC /l 0x804 /d "_DEBUG"
BSC32=bscmake.exe
# ADD BASE BSC32 /nologo
# ADD BSC32 /nologo
LIB32=link.exe -lib
# ADD BASE LIB32 /nologo
# ADD LIB32 /nologo /out:"../bin\algorithm-d.lib"
# Begin Special Build Tool
SOURCE="$(InputPath)"
PostBuild_Cmds=md       ..\include\      	copy       ..\algorithm\*.h       ..\include\ 
# End Special Build Tool

!ENDIF 

# Begin Target

# Name "algorithm - Win32 Release"
# Name "algorithm - Win32 Debug"
# Begin Source File

SOURCE=..\algorithm\__def_fun_def.c
# End Source File
# Begin Source File

SOURCE=..\algorithm\__def_fun_def.h
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyAlgorithm.h
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyBinarySearch.c
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyfunctionDef.h
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyHeapAlg.c
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyInsertSort.c
# End Source File
# Begin Source File

SOURCE=..\algorithm\MyQuickSort.c
# End Source File
# End Target
# End Project

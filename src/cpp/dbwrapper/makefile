#----------------------------------------------------------------------------
#       Macros
#----------------------------------------------------------------------------
CXX           = g++
LD            = $(CXX) $(CCFLAGS) $(CPPFLAGS)
AR            = ar
PICFLAGS      = -fPIC
CPPFLAGS      = $(PICFLAGS) $(GENFLAGS) -D_REENTRANT -I"./" -I"/usr/local/include/mysql" -I"../../lib/common_cpplib" -I"/usr/local/freetds/include/"
OBJEXT        = .o
OUTPUT_OPTION = -o "$@"
COMPILE.cc    = $(CXX) $(CCFLAGS) $(CPPFLAGS) -c
LDFLAGS       = -L"."
CCC           = $(CXX)
MAKEFILE      = Makefile.dbwrapper
DEPENDENCIES  = .depend.$(MAKEFILE)
LTARGETDIR    = ./
CAT           = cat
MV            = mv -f
RM            = rm -rf
CP            = cp -p
NUL           = /dev/null
MKDIR         = mkdir -p
EXEEXT        = 
LIBPREFIX     = lib
LIBSUFFIX     = 
GENFLAGS      = -O
LDLIBS        = -ldl -lpthread $(subst lib,-l,$(sort $(basename $(notdir $(wildcard /usr/lib/librt.so /lib/librt.so)))))
OBJS          = ./mysqlwrapper$(OBJEXT)
AREXT         = .a
LIB           = $(LTARGETDIR)$(LIBPREFIX)dbwrapper$(LIBSUFFIX)$(AREXT)
SRC           = ./mysqlwrapper.cpp
LINK.cc       = $(LD) $(LDFLAGS)
EXPORTFLAGS   = $(STATICFLAGS)

#----------------------------------------------------------------------------
#       Local targets
#----------------------------------------------------------------------------

all: $(LIB)

$(LTARGETDIR):
	@$(MKDIR) "$@"

$(LIB): $(LTARGETDIR) $(OBJS)
	$(AR) $(ARFLAGS) $(LIB) $(OBJS)

generated: $(GENERATED_DIRTY)
	@-:

./mysqlwrapper$(OBJEXT): ./mysqlwrapper.cpp
	$(COMPILE.cc) $(EXPORTFLAGS) $(OUTPUT_OPTION) ./mysqlwrapper.cpp

clean:
	-$(RM) $(OBJS)

realclean: clean
	-$(RM) $(SHLIB) $(LIB)

CPPFLAGS += -Wno-deprecated -g -D_FILE_OFFSET_BITS=64

#----------------------------------------------------------------------------
#       Dependencies
#----------------------------------------------------------------------------

$(DEPENDENCIES):
	@touch $(DEPENDENCIES)

depend:
	-makedepend -Y $(CFLAGS) $(CCFLAGS) $(CPPFLAGS) -f $(DEPENDENCIES) $(SRC) 2> $(NUL)

include $(DEPENDENCIES)

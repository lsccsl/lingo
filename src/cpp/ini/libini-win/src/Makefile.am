# The libraries

lib_LTLIBRARIES   = libini.la libtclini.la
libini_la_LDFLAGS = -version-info $(LIBVERSION)
libini_la_SOURCES = ini.cpp ini.h \
headings.h headings.i \
keys.h keys.i \
list.h list.i \
types.i

# Other language support.  Just rename the libini_wrap.c generated
# by SWIG to libini_<language>.c
libtclini_la_LDFLAGS = -version-info $(LIBVERSION)
libtclini_la_LIBADD  = ./libini.la
libtclini_la_SOURCES = libini_tcl.c

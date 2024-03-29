/***************************************************************************
                          libini.h  -  Header file of functions to
                                       manipulate an ini file.
                             -------------------
    begin                : Fri Apr 21 2000
    copyright            : (C) 2000 by Simon White
    email                : s_a_white@email.com
 ***************************************************************************/

/***************************************************************************
 *                                                                         *
 *   This program is free software; you can redistribute it and/or modify  *
 *   it under the terms of the GNU General Public License as published by  *
 *   the Free Software Foundation; either version 2 of the License, or     *
 *   (at your option) any later version.                                   *
 *                                                                         *
 ***************************************************************************/

#ifndef _libini_h_
#define _libini_h_

#ifdef __cplusplus
extern "C" {
#endif

/* Rev 1.3 Added scripting support using Swig 1.3a5 */
#ifdef SWIG
%module libini
#endif

#include <stdio.h>

#define INI_ADD_LIST_SUPPORT
#define INI_ADD_EXTRA_TYPES

typedef void* ini_fd_t;

/* DLL building support on win32 hosts */
#ifndef INI_EXTERN
#   ifdef DLL_EXPORT          /* defined by libtool (if required) */
#       define INI_EXTERN __declspec(dllexport)
#   endif
#   ifdef LIBINI_DLL_IMPORT   /* define if linking with this dll */
#       define INI_EXTERN extern __declspec(dllimport)
#   endif
#   ifndef INI_EXTERN         /* static linking or !_WIN32 */
#       define INI_EXTERN extern
#   endif
#endif


#ifdef SWIG
%include typemaps.i
%apply int    *BOTH { int    *value };
%apply long   *BOTH { long   *value };
%apply double *BOTH { double *value };
%name (ini_readString)
    int ini_readFileToBuffer    (ini_fd_t fd, ini_buffer_t *buffer);
%name (ini_writeString)
    int ini_writeFileFromBuffer (ini_fd_t fd, ini_buffer_t *buffer);

ini_buffer_t *ini_createBuffer        (unsigned long size);
void          ini_deleteBuffer        (ini_buffer_t *buffer);
char         *ini_getBuffer           (ini_buffer_t *buffer);
int           ini_setBuffer           (ini_buffer_t *buffer, char *str);

%{
#include "libini.h"

typedef struct 
{
    char  *buffer;
    size_t size;
} ini_buffer_t;

ini_buffer_t *ini_createBuffer (unsigned int size)
{
    ini_buffer_t *b;
    /* Allocate memory to structure */
    if (size == ( ((unsigned) -1 << 1) >> 1 ))
        return 0; /* Size is too big */
    b = malloc (sizeof (ini_buffer_t));
    if (!b)
        return 0;

    /* Allocate memory to buffer */
    b->buffer = malloc (sizeof (char) * (size + 1));
    if (!b->buffer)
    {
        free (b);
        return 0;
    }
    b->size = size;

    /* Returns address to tcl */
    return b;
}

void ini_deleteBuffer (ini_buffer_t *buffer)
{
    if (!buffer)
        return;
    free (buffer->buffer);
    free (buffer);
}

/*************************************************************
 * SWIG helper functions to create C compatible string buffers
 *************************************************************/
int ini_readFileToBuffer (ini_fd_t fd, ini_buffer_t *buffer)
{
    if (!buffer)
        return -1;
    return ini_readString (fd, buffer->buffer, buffer->size + 1);
}

int ini_writeFileFromBuffer (ini_fd_t fd, ini_buffer_t *buffer)
{
    if (!buffer)
        return -1;
    return ini_writeString (fd, buffer->buffer);
}

char *ini_getBuffer (ini_buffer_t *buffer)
{
    if (!buffer)
        return "";
    return buffer->buffer;
}

int ini_setBuffer (ini_buffer_t *buffer, char *str)
{
    size_t len;
    if (!buffer)
        return -1;
    len = strlen (str);
    if (len > buffer->size)
        len = buffer->size;

    memcpy (buffer->buffer, str, len);
    buffer->buffer[len] = '\0';
    return len;
}

%}

#endif /* SWIG */


/* Rev 1.2 Added new fuction */
INI_EXTERN ini_fd_t ini_new      (char *name);
INI_EXTERN ini_fd_t ini_open     (char *name);
INI_EXTERN int      ini_close    (ini_fd_t fd);
INI_EXTERN int      ini_flush    (ini_fd_t fd);

/* Rev 1.2 Added these functions to make life a bit easier, can still be implemented
 * through ini_writeString though. */
INI_EXTERN int ini_locateKey     (ini_fd_t fd, char *key);
INI_EXTERN int ini_locateHeading (ini_fd_t fd, char *heading);
INI_EXTERN int ini_deleteKey     (ini_fd_t fd);
INI_EXTERN int ini_deleteHeading (ini_fd_t fd);

/* Returns the number of bytes required to be able to read the key as a
 * string from the file. (1 should be added to this length to account
 * for a NULL character) */
INI_EXTERN int ini_dataLength (ini_fd_t fd);

/* Default Data Type Operations
 * Arrays implemented to help with reading, for writing you should format the
 * complete array as a string and perform an ini_writeString. */
#ifndef SWIG
INI_EXTERN int ini_readString  (ini_fd_t fd, char *str, size_t size);
INI_EXTERN int ini_writeString (ini_fd_t fd, char *str);
#endif /* SWIG */
INI_EXTERN int ini_readInt     (ini_fd_t fd, int  *value);


#ifdef INI_ADD_EXTRA_TYPES
    /* Read Operations */
    INI_EXTERN int ini_readLong    (ini_fd_t fd, long   *value);
    INI_EXTERN int ini_readDouble  (ini_fd_t fd, double *value);

    /* Write Operations */
    INI_EXTERN int ini_writeInt    (ini_fd_t fd, int    value);
    INI_EXTERN int ini_writeLong   (ini_fd_t fd, long   value);
    INI_EXTERN int ini_writeDouble (ini_fd_t fd, double value);
#endif /* INI_ADD_EXTRA_TYPES */


#ifdef INI_ADD_LIST_SUPPORT
    /* Rev 1.1 Added - List Operations (Used for read operations only)
     * Be warned, once delimiters are set, every key that is read will be checked for the
     * presence of sub strings.  This will incure a speed hit and therefore once a line
     * has been read and list/array functionality is not required, set delimiters
     * back to NULL.
     */

    /* Returns the number of elements in an list being seperated by the provided delimiters */
    INI_EXTERN int ini_listLength      (ini_fd_t fd);
    /* Change delimiters, default "" */
    INI_EXTERN int ini_listDelims      (ini_fd_t fd, char *delims);
    /* Set index to access in a list.  When read the index will automatically increment */
    INI_EXTERN int ini_listIndex       (ini_fd_t fd, unsigned long index);
    /* Returns the length of an indexed sub string in the list */
    INI_EXTERN int ini_listIndexLength (ini_fd_t fd);
#endif // INI_ADD_LIST_SUPPORT

#ifdef __cplusplus
}
#endif

#endif /* _libini_h_ */

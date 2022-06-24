#ifndef STRING_H
#define STRING_H

#include <stdlib.h>

int pstrncmp(const char *p, const char *q, size_t n);

char *pstrncpy(char *s, const char *t, int n);

char *pstrcpy(char *s, const char *t, int n);

int pstrlen(const char *s);

typedef struct _String
{
    size_t length;
    char *data;
} String;

String *string_create_by_length(size_t length);

String *string_create_by_str(const char *str);

void string_destroy(String *str);

int string_length(String *str);

String *string_concat(String *str1, String *str2);

#endif
#include "string.h"

int pstrncmp(const char *p, const char *q, size_t n)
{
  while (n > 0 && *p && *p == *q)
    n--, p++, q++;
  if (n == 0)
    return 0;
  return (char)*p - (char)*q;
}

char *pstrncpy(char *s, const char *t, int n)
{
  char *os;

  os = s;
  while (n-- > 0 && (*s++ = *t++) != 0)
    ;
  while (n-- > 0)
    *s++ = 0;
  return os;
}

// Like strncpy but guaranteed to NUL-terminate.
char *pstrcpy(char *s, const char *t, int n)
{
  char *os;

  os = s;
  if (n <= 0)
    return os;
  while (--n > 0 && (*s++ = *t++) != 0)
    ;
  *s = 0;
  return os;
}

int pstrlen(const char *s)
{
  int n;

  for (n = 0; s[n]; n++)
    ;
  return n;
}

String *string_create_by_length(size_t length)
{
    // 申请内存
    size_t size = sizeof(size_t) + sizeof(unsigned char) * (length + 1);
    String *pstr = (String *)malloc(size);
    // 设置字符串长度
    pstr->length = length;
    // 设置数据指针
    pstr->data = (char *)(pstr + sizeof(size_t));

    return pstr;
}

String *string_create_by_str(const char *str)
{
    size_t str_length = pstrlen(str);

    String *pstr = string_create_by_length(str_length);

    // 拷贝数据
    pstrcpy(pstr->data, str, pstr->length);
    return pstr;
}

void string_destroy(String *str)
{
    free(str);
}

int string_length(String *str)
{
    return str->length;
}

String *string_concat(String *str1, String *str2)
{
    size_t size = sizeof(size_t) + sizeof(unsigned char) * (str1->length + str2->length + 1);
    // 申请内存
    String *pstr = (String *)malloc(size);
    // 设置字符串长度
    pstr->length = str1->length + str2->length;
    // 设置数据指针
    pstr->data = (char *)(pstr + sizeof(size_t));
    // 拷贝数据
    pstrcpy(pstr->data, str1->data, str1->length);
    pstrcpy(pstr->data + str1->length, str2->data, str2->length);
    return pstr;
}
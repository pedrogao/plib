#ifndef OP_H
#define OP_H

#include <stdint.h>

char isspace(char ch);

// < 10000
int u32toa_small(char *out, uint32_t val);

#endif
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char const *argv[])
{
    int *a = (int *)malloc(sizeof(int));
    *a = 1;

    printf("%d\n", *a);

    free(a);
    return 0;
}

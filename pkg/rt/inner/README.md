1. generate asm by clang

```shell
clang -mno-red-zone -fno-asynchronous-unwind-tables -fno-builtin -fno-exceptions \
-fno-rtti -fno-stack-protector -nostdlib -O3 -msse4 -mavx -mno-avx2 -DUSE_AVX=1 \
 -DUSE_AVX2=0 -S ./*.c
```

2. transform x86 asm to plan9

```shell
../../tools/asm2asm.py ./op.s ./inner/op.s
```
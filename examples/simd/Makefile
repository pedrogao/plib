all:
	@clang -S -mavx2 -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 \
	-fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -o ./lib/sum_float64.s  -c ./lib/sum_float64.c
	# @go install  github.com/minio/c2goasm
	# @go install  github.com/minio/asm2plan9s
	@c2goasm -a -f ./lib/sum_float64.s sum_float64.s
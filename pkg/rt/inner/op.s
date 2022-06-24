	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 11, 0
	.globl	_isspace                        ## -- Begin function isspace
	.p2align	4, 0x90
_isspace:                               ## @isspace
## %bb.0:
	pushq	%rbp
	movq	%rsp, %rbp
	movb	$1, %al
	cmpb	$13, %dil
	je	LBB0_3
## %bb.1:
	cmpb	$32, %dil
	jne	LBB0_2
LBB0_3:
	movzbl	%al, %eax
	popq	%rbp
	retq
LBB0_2:
	addb	$-9, %dil
	cmpb	$2, %dil
	setb	%al
	movzbl	%al, %eax
	popq	%rbp
	retq
                                        ## -- End function
	.globl	_u32toa_small                   ## -- Begin function u32toa_small
	.p2align	4, 0x90
_u32toa_small:                          ## @u32toa_small
## %bb.0:
	pushq	%rbp
	movq	%rsp, %rbp
	movl	%esi, %eax
	imulq	$1374389535, %rax, %rax         ## imm = 0x51EB851F
	shrq	$37, %rax
	leaq	(%rax,%rax), %rdx
	imull	$100, %eax, %eax
	movl	%esi, %ecx
	subl	%eax, %ecx
	addq	%rcx, %rcx
	cmpl	$1000, %esi                     ## imm = 0x3E8
	jb	LBB1_2
## %bb.1:
	leaq	_Digits(%rip), %rax
	movb	(%rdx,%rax), %al
	movb	%al, (%rdi)
	movl	$1, %eax
	jmp	LBB1_3
LBB1_2:
	xorl	%eax, %eax
	cmpl	$100, %esi
	jb	LBB1_4
LBB1_3:
	movl	%edx, %edx
	orq	$1, %rdx
	leaq	_Digits(%rip), %rsi
	movb	(%rdx,%rsi), %dl
	movl	%eax, %esi
	addl	$1, %eax
	movb	%dl, (%rdi,%rsi)
LBB1_5:
	leaq	_Digits(%rip), %rdx
	movb	(%rcx,%rdx), %dl
	movl	%eax, %esi
	addl	$1, %eax
	movb	%dl, (%rdi,%rsi)
LBB1_6:
	movl	%ecx, %ecx
	orq	$1, %rcx
	leaq	_Digits(%rip), %rdx
	movb	(%rcx,%rdx), %cl
	movl	%eax, %edx
	addl	$1, %eax
	movb	%cl, (%rdi,%rdx)
	popq	%rbp
	retq
LBB1_4:
	xorl	%eax, %eax
	cmpl	$10, %esi
	jae	LBB1_5
	jmp	LBB1_6
                                        ## -- End function
	.section	__TEXT,__const
	.p2align	4                               ## @Digits
_Digits:
	.ascii	"00010203040506070809101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899"

.subsections_via_symbols

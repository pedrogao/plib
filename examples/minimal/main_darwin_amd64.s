
TEXT ·entry(SB), $0-0
    MOVL    $33, DI             // arg 1 exit status
    MOVL    $(0x2000000+1), AX  // syscall entry
    SYSCALL
    RET



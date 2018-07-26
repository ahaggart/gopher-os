#include "textflag.h"

TEXT ·Interrupt(SB),NOSPLIT,$0
	INT $0x80
    RET

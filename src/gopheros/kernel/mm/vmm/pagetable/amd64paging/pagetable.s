#include "textflag.h"

TEXT LoadCR3(SB),NOSPLIT,$0-8
	MOVQ AX, CR3
	MOVQ BX, $((1 << 48 - 1) << 12)
	ANDQ AX, BX
	ORQ AX, p(SB)
	MOVQ CR3, AX
	RET

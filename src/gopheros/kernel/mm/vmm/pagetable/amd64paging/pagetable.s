#include "textflag.h"

TEXT LoadCR3(SB),NOSPLIT,$0-8
	MOVQ AX, CR3
	MOVQ BX, cr3Mask(SB)
	ANDQ AX, BX
	ORQ AX, p(FP)
	MOVQ CR3, AX
	RET

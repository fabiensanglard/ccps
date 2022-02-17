	.include "Defines.inc"

	.globl	LoadPalette
	.globl	LoadVRAM
	.globl	FillVRAM
	.globl	ReadJoystickA
	.globl	ReadJoystickB
	.globl	SetPaletteRed
	.globl	SetPaletteGreen
	.globl	SetPaletteBlue
	.globl	DrawNumber
	.globl	DrawString
	.globl	ClearHardwareSprites
	.globl	ClearScrolls
	.globl	TurnOffScreen
	.globl	TurnOnScreen
	.globl	WaitVBlank
	.globl	VBlank
	.globl	InitSound
	.globl	PlaySound
	.globl	Seed
	.globl	Random

	.extern	SpritesChanged

* void LoadPalette(DWORD Source, DWORD Dest, DWORD Length);
* A0 - Source address A1 - Dest address  D0 - Size of data	
LoadPalette:
	.set	_ARGS, 4
	move.l	_ARGS(sp), a0
	move.l	_ARGS+4(sp), a1
	move.l	_ARGS+8(sp), d0

	movem.l d0-d7/a0-a4, -(sp)

	* Divide size by 4 and adjust for loop
	asr.l	#2, d0
	sub.l	#1, d0

LoadPaletteLoop:
	move.l	(a0)+, (a1)+
	dbra	d0, LoadPaletteLoop

	movem.l (sp)+, d0-d7/a0-a4

	rts

* void LoadVRAM(DWORD Source, DWORD Dest, DWORD Length);
* A0 - Source address A1 - Dest address  D0 - Size of data	
LoadVRAM:
	.set	_ARGS, 4
	move.l	_ARGS(sp), a0
	move.l	_ARGS+4(sp), a1
	move.l	_ARGS+8(sp), d0

	movem.l d0-d7/a0-a4, -(sp)

	* Divide size by 4 and adjust for loop
	asr.l	#2, d0
	sub.l	#1, d0

LoadVRAMLoop:
	move.l	(a0)+, (a1)+
	dbra	d0, LoadVRAMLoop

	movem.l (sp)+, d0-d7/a0-a4

	rts

* void FillVRAM(DWORD Address, DWORD Value, DWORD Count);
* A0 - VRAM address D0 - DWORDS to fill D1 - Fill Value, 
FillVRAM:
	.set	_ARGS, 4
	move.l	_ARGS(sp), a0
	move.l	_ARGS+4(sp), d1
	move.l	_ARGS+8(sp), d0

	movem.l d0-d7/a0-a4, -(sp)

	* Adjust for loop
	sub.l	#1, d0

FillVRAMLoop:
	move.l	d1, (a0)+
	dbra	d0, FillVRAMLoop

	movem.l (sp)+, d0-d7/a0-a4

	rts

* Turn off screen
TurnOffScreen:
	rts

* Turn on screen
TurnOnScreen:
	rts

* DWORD ReadJoystickA();
ReadJoystickA:
	move.l	JoystickValA, d0

	rts

* DWORD ReadJoystickB();
ReadJoystickB:
	move.l	JoystickValB, d0

	rts

* Clear hardware sprite table
ClearHardwareSprites:
	lea		SPRITETABLE, a0
	move.l	#255, d0

ClearHardwareSpritesLoop:
	clr.l	(a0)+
	clr.l	(a0)+

	dbra	d0, ClearHardwareSpritesLoop

	rts

* Clear all 3 scrolls
ClearScrolls:
	lea		SCROLL1, a0
	move.l	#10000/4-1, d0

ClearScrollsLoop:
	clr.l	(a0)+
	dbra	d0, ClearScrollsLoop
	
	rts

* void SetPaletteRed(DWORD Dest, DWORD Red);
* D0 - Palette pointer
* D1 - Color
SetPaletteRed:
	.set	_ARGS, 4
	move.l	_ARGS(sp), d0
	move.l	_ARGS+4(sp), d1

	movem.l	d0-d2/a0, -(sp)

	lea		BACKPAL, a0
	add.l	d0, a0

	* Get current color, and out the red component. Need to keep upper 4 bits.
	move.w	(a0), d2
	and.w	#0xF0FF, d2
	
	* Shift color bits into positions
	asl.w	#8, d1
	or.w	d1, (a0)

	movem.l (sp)+, d0-d2/a0

	rts

* void SetPaletteGreen(DWORD Dest, DWORD Green);
* D0 - Palette pointer
* D1 - Color
SetPaletteGreen:
	.set	_ARGS, 4
	move.l	_ARGS(sp), d0
	move.l	_ARGS+4(sp), d1

	movem.l	d0-d2/a0, -(sp)

	lea		BACKPAL, a0
	add.l	d0, a0

	* Get current color, and out the green component. Need to keep upper 4 bits.
	move.w	(a0), d2
	and.w	#0xFF0F, d2
	
	* Shift color bits into positions
	asl.w	#4, d1
	or.w	d1, (a0)

	movem.l (sp)+, d0-d2/a0

	rts

* void SetPaletteBlue(DWORD Dest, DWORD Blue);
* D0 - Palette pointer
* D1 - Color
SetPaletteBlue:
	.set	_ARGS, 4
	move.l	_ARGS(sp), d0
	move.l	_ARGS+4(sp), d1

	movem.l	d0-d2/a0, -(sp)

	lea		BACKPAL, a0
	add.l	d0, a0

	* Get current color, and out the blue component. Need to keep upper 4 bits.
	move.w	(a0), d2
	and.w	#0xFFF0, d2
	
	or.w	d1, (a0)

	movem.l (sp)+, d0-d2/a0

	rts

* void DrawNumber(DWORD Value, DWORD Palette, DWORD x, DWORD y);
* D1 - Value, D2 - Palette, D3 - x, D4 - y, D5 - Location
DrawNumber:
	.set	_ARGS, 4

	* Get value
	move.l	_ARGS(sp), d1		
	
	* Get palette
	move.l	_ARGS+4(sp), d2		

	* Get x
	move.l	_ARGS+8(sp), d3

	* Get y
	move.l	_ARGS+12(sp), d4		

	* Divide score by 10
	move.l	#10, -(sp)			
	move.l	d1, -(sp)
	bsr		ldiv
	addq.w	#8, sp

	add.l	#48, d1
	jsr		DrawChar

	cmp.b	#0, d0
	beq		NoDrawNumber

	sub.l	#1, d3
	move.l	d0, d1
	add.l	#48, d1
	jsr		DrawChar

NoDrawNumber:
	rts

* void DrawString(char *String, DWORD Palette, DWORD x, DWORD y);
* A0 - String, D2 - Palette, D3 - x, D4 - y, D5 - Location
DrawString:
	.set	_ARGS, 4

	* Get value
	move.l	_ARGS(sp), a0
	
	* Get palette
	move.l	_ARGS+4(sp), d2		

	* Get x
	move.l	_ARGS+8(sp), d3

	* Get y
	move.l	_ARGS+12(sp), d4		

	clr.w	d1

ContinueDrawString:
	move.b	(a0)+, d1

	cmp.b	#0, d1
	beq.s	EndDrawString

	jsr		DrawChar

	add.l	#1, d3

	bra.s	ContinueDrawString

EndDrawString:
	rts

* D1 - Character, D2 - Palette, D3 - x, D4 - y
DrawChar:
	movem.l d3-d4/a0, -(sp)

	sub.w	#32, d1

	* Multiply x by 128, screen is layed out in columns
	asl.l	#7, d3

	* Mulitply y by 4
	asl.l	#2, d4

	lea		SCROLL1, a0
	add.l	d3, a0
	add.l	d4, a0

	move.w	d1, (a0)+
	move.w	d2, (a0)+

	movem.l (sp)+, d3-d4/a0

	rts

* void WaitVBlank();
WaitVBlank:
	cmp.l	#1, VertBlank
	bne		WaitVBlank

	cmp.l	#1, SpritesChanged
	bne.s	EndWaitVBlank
	
	jsr		UpdateSprites

	move.l	#0, SpritesChanged

EndWaitVBlank:
	clr.l	VertBlank

	rts

* Perform normal Vertical Blank Tasks (transfer sprites, read joystick)
VBlank:
	* jsr		ReadJoysticks
	* jsr		UpdateSound
    jsr onVSync
	move.l	#1, VertBlank
	
	rte

* Read joysticks
ReadJoysticks:
	movem.l	d0-d1, -(sp)

	move.w	JOYSTICK, d0
	not.w	d0
	move.w	d0, d1

	and.l	#0x00FF, d0
	move.l	d0, JoystickValA

	move.w	d1, d0
	asr.w	#8, d0

	and.l	#0x00FF, d0
	move.l	d0, JoystickValB

	movem.l	(sp)+, d0-d1

	rts

UpdateSprites:
	movem.l d0/a0-a1, -(sp)

	move.l	#MAXSPRITES, d0
	lea		SpriteBlock, a0
	lea		SPRITETABLE, a1

UpdateSpritesLoop:
	move.l	(a0)+, (a1)+
	move.l	(a0)+, (a1)+
	
	add.l	#SPRITEBLOCKSIZE-8, a0

	dbra	d0, UpdateSpritesLoop
		
	movem.l (sp)+, d0/a0-a1

	rts

	.globl	UpdateSound

UpdateSound:
	cmp.l	#1, SoundsChanged
	bne.s	UpdateSoundSkip

	move.l	SoundID, d0
	move.w	d0, 0x800188
	move.w	d0, 0x800180
	
	move.l	#0, SoundsChanged

	rts

UpdateSoundSkip:
	move.w	#0xFF, 0x800180

	rts

* void InitSound();
InitSound:
	move.l	#0xF0, SoundID
	move.l	#1, SoundsChanged

	rts

* void PlaySound(DWORD SoundID);
PlaySound:
	.set	_ARGS, 4

	* Get sound id
	move.l	_ARGS(sp), d0

	move.l	d0, SoundID
	move.l	#1, SoundsChanged

	rts

* D0 - Seed
Seed:
	* user seed in d0 (d1 too)
	add.l	d0, d1
	movem.l	d0-d1, RandomSeed

* D0 - Limit
LongRnd:
	movem.l d2-d3, -(sp)
	
	* d0=LSB's, d1=MSB's of random number
	movem.l RandomSeed, d0-d1
	
	* ensure upper 59 bits are an...
	andi.b	#0x0E, d0
	
	* ...odd binary number
	ori.b	#0x20, d0
	move.l	d0, d2
	move.l	d1, d3
	
	* accounts for 1 of 17 left shifts
	add.l	d2, d2
	
	* [d2/d3] = RND*2
	addx.l	d3, d3
	add.l	d2, d0
	
	* [d0/d1] = RND*3
	addx.l	d3, d1
	
	* shift [d2/d3] additional 16 times
	swap	d3
	swap	d2
	move.w	d2, d3
	clr.w	d2
	
	* add to [d0/d1]
	add.l	d2, d0
	addx.l	d3, d1
	
	* save for next time through
	movem.l d0-d1, RandomSeed
	
	* most random part to d0
	move.l	d1, d0
	movem.l (sp)+, d2-d3

	rts

* void Random(DWORD Limit);
* D0 - Limit
Random:
	.set	_ARGS, 4

	* Get value
	move.l	_ARGS(sp), d0

	move.l	d2, -(sp)

	* save upper limit
	move.w	d0, d2
	cmp 	#0, d0

	* range of 0 returns 0 always
	beq.s	Zero

	* get a longword random number
	bsr 	LongRnd

	* use upper word (it's most random)
	clr.w	d0
	swap	d0

	* divide by range...
	divu.w	d2, d0

	* ...and use remainder for the value
	clr.w	d0

	* result in d0.w
	swap	d0

Zero:
	move.l	(sp)+, d2

	rts

ldiv:
	move.l  4(sp),d0
	bpl     ld1
	neg.l   d0
ld1:
	move.l  8(sp),d1
	bpl     ld2
	neg.l   d1
	eor.b   #0x80,4(sp)
ld2:
	bsr     i_ldiv          /* d0 = d0/d1 */
	tst.b   4(sp)
	bpl     ld3
	neg.l   d0
ld3:
	rts

lmul:
	move.l  4(sp),d0
	bpl     lm1
	neg.l   d0
lm1:
	move.l  8(sp),d1
	bpl     lm2
	neg.l   d1
	eor.b   #0x80,4(sp)
lm2:
	bsr     i_lmul          /* d0 = d0*d1 */
	tst.b   4(sp)
	bpl     lm3
	neg.l   d0

lm3:
	rts

*
* A in d0, B in d1, return A*B in d0
*
i_lmul:
	move.l  d3,a2           /* save d3 */
	move.w  d1,d2
	mulu    d0,d2           /* d2 = Al * Bl */

	move.l  d1,d3
	swap    d3
	mulu    d0,d3           /* d3 = Al * Bh */

	swap    d0
	mulu    d1,d0           /* d0 = Ah * Bl */

	add.l   d3,d0           /* d0 = (Ah*Bl + Al*Bh) */
	swap    d0
	clr.w   d0              /* d0 = (Ah*Bl + Al*Bh) << 16 */

	add.l   d2,d0           /* d0 = A*B */
	move.l  a2,d3           /* restore d3 */
	rts
*
*A in d0, B in d1, return A/B in d0, A%B in d1
*
i_ldiv:
	tst.l   d1
	bne     nz1

	move.l  #0x80000000,d0
	move.l  d0,d1
	rts
nz1:
	move.l  d3,a2           /* save d3 */
	cmp.l   d1,d0
	bhi     norm
	beq     is1
*       A<B, so ret 0, rem A
	move.l  d0,d1
	clr.l   d0
	move.l  a2,d3           /* restore d3 */
	rts
*       A==B, so ret 1, rem 0
is1:
	moveq.l #1,d0
	clr.l   d1
	move.l  a2,d3           /* restore d3 */
	rts
*       A>B and B is not 0
norm:
	cmp.l   #1,d1
	bne     not1
*       B==1, so ret A, rem 0
	clr.l   d1
	move.l  a2,d3           /* restore d3 */
	rts
*  check for A short (implies B short also)
not1:
	cmp.l   #0xffff,d0
	bhi     slow
*  A short and B short -- use 'divu'
	divu    d1,d0           /* d0 = REM:ANS */
	swap    d0              /* d0 = ANS:REM */
	clr.l   d1
	move.w  d0,d1           /* d1 = REM */
	clr.w   d0
	swap    d0
	move.l  a2,d3           /* restore d3 */
	rts
* check for B short
slow:
	cmp.l   #0xffff,d1
	bhi     slower
* A long and B short -- use special stuff from gnu
	move.l  d0,d2
	clr.w   d2
	swap    d2
	divu    d1,d2           /* d2 = REM:ANS of Ahi/B */
	clr.l   d3
	move.w  d2,d3           /* d3 = Ahi/B */
	swap    d3

	move.w  d0,d2           /* d2 = REM << 16 + Alo */
	divu    d1,d2           /* d2 = REM:ANS of stuff/B */

	move.l  d2,d1
	clr.w   d1
	swap    d1              /* d1 = REM */

	clr.l   d0
	move.w  d2,d0
	add.l   d3,d0           /* d0 = ANS */
	move.l  a2,d3           /* restore d3 */
	rts
*       A>B, B > 1
slower:
	move.l  #1,d2
	clr.l   d3
moreadj:
	cmp.l   d0,d1
	bhs     adj
	add.l   d2,d2
	add.l   d1,d1
	bpl     moreadj
* we shifted B until its >A or sign bit set
* we shifted #1 (d2) along with it
adj:
	cmp.l   d0,d1
	bhi     ltuns
	or.l    d2,d3
	sub.l   d1,d0
ltuns:
	lsr.l   #1,d1
	lsr.l   #1,d2
	bne     adj
* d3=answer, d0=rem
	move.l  d0,d1
	move.l  d3,d0
	move.l  a2,d3           /* restore d3 */

	rts

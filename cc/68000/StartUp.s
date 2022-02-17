* CPS-1 Frog Feast (c) RasterSoft

	.include "Defines.inc"

********************** Exported Symbols **********************
	.globl	_start
	.globl	atexit
	.globl	RandomSeed
	.globl	VertBlank
	.globl	JoyStickValA
	.globl	JoyStickValB

********************** Imported Symbols **********************
	.extern	run
	.extern	_end
	.extern	VBlank
	.extern	Seed

	dc.l	0xFFF000, _start, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, VBlank,  Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default
	dc.l	Default, Default, Default, Default, Default, Default, Default, Default

	.align	4

Default:
	rte

	.align	4

HBlank:
	rte

	.align	4

* Dummy atexit (does nothing for now)
atexit:
	moveq	#0, d0
	rts

	.align	4

_start:
	* Set stack pointer
	move.l	#0xFFF000, sp


	* Initialize BSS section
	move.l	#_end, d0
	sub.l	#__bss_start, d0
	move.l	d0, -(sp)
	clr.l	-(sp)
	pea		__bss_start

	* Seed random number generator
	move.w	SCROLL1, d0
	jsr		Seed

	* Clear scrolls
	jsr		ClearScrolls

	* Clear hardware sprites
	jsr		ClearHardwareSprites

	* Reset coin control
	move.b	#0x80, 0x800030
	nop
	nop
	nop
	nop
	move.b	#0x00, 0x800030

	* Screen setup
	move.w	#0x000E, 0x800122
	
	* Set control register
	move.w	#0x003F, 0x80016A
	
	* Set video control register (Enables graphic scrolls)
	move.w	#0x12CE, 0x80016E

	* Sets object base / 256 (in gfx ram).
	move.w	#0x9000, 0x800100

	* Sets scroll 1 base / 256 (in gfx ram).
	move.w	#0x9080, 0x800102

	* Sets scroll 2 base / 256 (in gfx ram).
	move.w	#0x90c0, 0x800104

	* Sets scroll 3 base / 256 (in gfx ram).
	move.w	#0x9100, 0x800106

	* Sets scroll distortion base / 256 (in gfx ram). (Need to find values)
	move.w	#0x9200, 0x800108

	* Sets palette base / 256 (in gfx ram).
	move.w	#0x9140, 0x80010A

	* Sets scroll 1 x
	move.w	#0xFFC0, 0x80010C

	* Sets scroll 1 y
	move.w	#0xFFF0, 0x80010E

	* Sets scroll 2 x
	move.w	#0xFFC0, 0x800110

	* Sets scroll 2 y
	move.w	#0xFFF0, 0x800112
	
	* Sets scroll 3 x
	move.w	#0xFFC0, 0x800114

	* Sets scroll 3 y
	move.w	#0xFFF0, 0x800116

	* Sets rowscroll matrix offset. (Need to find values)
	move.w	#0x9200, 0x800120

	* Enable interrupts
	move.w	#0x2000, sr

	* Jump to mainloop
	jbsr	run

EndLoop:
	bra.s	EndLoop

	.align	4


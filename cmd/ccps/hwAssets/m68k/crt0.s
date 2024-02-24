	.globl	_start
	.globl	atexit

	.extern hardwareInit
	.extern onVSync
	.extern	run
	.extern	_end

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
* Dummy atexit (does nothing for now)
atexit:
	moveq	#0, d0
	rts

.align	4
VBlank:
	jsr onVSync
	rte

.align	4
_start:
	* Reset coin control
	move.b	#0x80, 0x800030
	nop
	nop
	nop
	move.b	#0x00, 0x800030

	* Init Main RAM
	lea 	0xff0000, A0
	lea 	0xffffff, A1
loopram:
	move.w	#0x0000, (A0)+
	cmpa.l	A1, A0
	bls 	loopram

	* Init GFX Memory
	lea 	0x900000, A0
	lea 	0x92ffff, A1
loopgfx:
	move.w	#0x0000, (A0)+
	cmpa.l	A1, A0
	bls 	loopgfx

	* Initialize BSS section
	move.l	#_end, d0
	sub.l	#__bss_start, d0
	move.l	d0, -(sp)
	clr.l	-(sp)
	pea		__bss_start

	* Call Board initialization before interrups are enabled
	* for CPSA & CPSB registers and any other data before calling run
	jsr 	hardwareInit

	* Enable interrupts
	move.w	#0x2000, sr

	* Jump to mainloop
	jbsr	run

EndLoop:
	bra.s	EndLoop



	.globl	_start
	.globl	atexit

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
	* Initialize BSS section
	move.l	#_end, d0
	sub.l	#__bss_start, d0
	move.l	d0, -(sp)
	clr.l	-(sp)
	pea		__bss_start

	* Enable interrupts
	move.w	#0x2000, sr

	* Jump to mainloop
	jbsr	run

EndLoop:
	bra.s	EndLoop



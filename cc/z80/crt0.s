   .module crt0
   .globl  _main
   .area  _HEADER (ABS)

;------------------
; Z-80 starts here!
;------------------
.org 0
   jp init  

;--------------
; INTERRUPT
;--------------
.org 0x38
   DI                     ; Disable Interrupt
   call _interrupt        ; Process Interrupt
   call _requestInterrupt ; Request the interrupt re-trigger
   EI                     ; Enable  Interrupt
   RET 

;--------------
; INIT and MAIN
;--------------
.org 0x100
init:
   
   ld  sp,#0xd7ff    ; Setup stack
   IM 1              ; Set Interupt mode 1
   call  gsinit      ; Init global variables
main:  
   call  _main       ; Infinite loop
   jp    _exit       ; Never happens

   ; Ordering of segments for the linker.
   .area _HOME
   .area _CODE
   .area _INITIALIZER
   .area _GSINIT
   .area _GSFINAL

   .area _DATA
   .area _INITIALIZED
   .area _BSEG
   .area _BSS
   .area _HEAP

   .area _CODE

_exit:
   jp main

;----------------------------
; Initialize global variables
; Copy values from ROM > RAM.
;----------------------------
   .area _GSINIT
gsinit::
   ld  bc, #l__INITIALIZER
   ld  a, b
   or  a, c
   jr  Z, gsinit_next
   ld  de, #s__INITIALIZED
   ld  hl, #s__INITIALIZER
   ldir
gsinit_next:
   .area _GSFINAL
   ret


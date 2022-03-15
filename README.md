# CPS-1 SDK

This a development kit for CPS-1 board.

Compilers required
==================

Tested with GCC 10.3.0
Tested with SDCC 4.0

Setup :
=======

CCPS is a compiler/tools driver. It needs GCC to compile code to 68000 and SDCC to compile code to Z-80. Everything else is self-contained. 

On Linux:

```
sudo apt install gcc-m68k-linux-gnu g++-m68k-linux-gnu binutils-m68k-linux-gnu
sudo apt install sdcc
```

On Windows:
Use WSL and follow Linux steps.

On MacOS X:
Never tested. 

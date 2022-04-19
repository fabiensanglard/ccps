__at (0xF000) char YM2151_ADD;
__at (0xF001) char YM2151_DAT;
__at (0xF002) char OKI;
__at (0xF008) char LATCH1;

#define YM2151_REG_CLKA1 0x10
#define YM2151_REG_CLKA2 0x11
#define YM2151_REG_CTRL  0x14

#define NO_OP 0xFF

unsigned char latch = 0;
unsigned char lastLatch = NO_OP;

void interrupt() {
  latch = LATCH1;

  if (lastLatch == latch) {
    return;
  }

  lastLatch = latch;

  if (latch == NO_OP) {
    return;
  }

  OKI = 0x80 | latch; // First bit must be 1 then sound ID.
  OKI = 0x80; // 0x80 = Channel 4  !   0x00 = No sound reduction.
}

void requestInterrupt() {
  // The YM2151 runs at the same speed at the Z-80 (3.579MHz)
  // We want to be interrupted at 4ms interval. In YM2151 ticks,
  // that means setting timer to 0xC000 (800) so it ticks at
  // 64 * ( 1024 - 800) / 3579
  YM2151_ADD = YM2151_REG_CLKA1;
  YM2151_DAT = 0xC0;

  YM2151_ADD = YM2151_REG_CLKA2;
  YM2151_DAT = 0x00;
}

void main() {
  // Enable timer A
  YM2151_ADD = YM2151_REG_CTRL;
  YM2151_DAT = 0x15;

  // Request the first interrupt (after that the interrupt handler
  // will call requestInterrupt() after calling interrupt().
  requestInterrupt();

  while(1) {
    interrupt();
   }
}
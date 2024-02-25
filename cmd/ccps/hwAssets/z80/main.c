__at (0xF000) char YM2151_ADR;
__at (0xF001) char YM2151_DAT;
__at (0xF002) char OKI;
__at (0xF008) char LATCH1;

#define YM2151_REG_CLKA1 0x10
#define YM2151_REG_CLKA2 0x11
#define YM2151_REG_CTRL  0x14

#define NO_OP 0xFF

unsigned char latch = 0;
unsigned char lastLatch = NO_OP;

void YM2151_writeReg(char adr, char dat) {
  while(YM2151_DAT == 0x80); // Wait until YM2151 is ready for write
  YM2151_ADR = adr;
  YM2151_DAT = dat;
}

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
  // We want to be interrupted at 4ms interval. That means
  // setting timer A to 800 so it ticks at 64 * ( 1024 - 800) / 3579.
  // The Timer A 10 bit is split across 2 registers (A1 and A2). We
  // write 0xC8 in the A1 and 0 in A2 resulting in value 0xC8 * 4 = 800.

  // 8 msb
  YM2151_writeReg(YM2151_REG_CLKA1, 0xC8);
  // 2 lsb
  YM2151_writeReg(YM2151_REG_CLKA2, 0x00);
  // Re-enable timer A
  YM2151_writeReg(YM2151_REG_CTRL, 0x15);
}

void main() {
  // Reset timer flags
  YM2151_writeReg(YM2151_REG_CTRL, 0x30);

  // Request the first interrupt. The interrupt will trigger a jump to
  // intructions at 0x38 (see crt0.s) which will call interrupt and
  // then requestInterrupt.
  requestInterrupt();

  // infinite loop, all will be done via timer and interrupts
  while(1) {
  }
}
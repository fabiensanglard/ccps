__at (0xF002) char OKI;
__at (0xF008) char LATCH1;
#define NO_OP 0xFF

void interrupt() {

}

void requestInterrupt() {
}

void main() {


  unsigned char latch = 0;
  unsigned char lastLatch = NO_OP;
  while(1) {
      latch = LATCH1;

	  // Tick one
	  if (lastLatch == latch) {
	 	continue;
	  }
      lastLatch = latch;

	  if (latch == NO_OP) {
		continue;
	  }

	  OKI = 0x80 | latch; // First bit must be 1 then sound ID.
      OKI = 0x80; // 0x80 = Channel 4  !   0x00 = No sound reduction.
   }
}
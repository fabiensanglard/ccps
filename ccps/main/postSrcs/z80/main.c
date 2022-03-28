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
	  //while(interruptCounter == mainCounter){
	  //}
	  //mainCounter++;
      latch = LATCH1;

	  // Tick one
	  if (lastLatch == latch) {
	 	continue;
	  }
      lastLatch = latch;

	  if (latch == NO_OP) {
		continue;
	  }

	  OKI = 0x80 | latch;
    OKI = 0x81;
   }
}
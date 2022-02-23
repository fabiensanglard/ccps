// https://gcc.gnu.org/onlinedocs/gcc-3.2/gcc/Variable-Attributes.html
// https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/html_chapter/ld_3.html
// https://mcuoneclipse.com/2016/11/01/getting-the-memory-range-of-sections-with-gnu-linker-files/
#define GFXRAM   __attribute__ ((section (".gfx_data")))  
#define CPSA_REG __attribute__ ((section (".cpsa_reg")))  
#define CPSB_REG __attribute__ ((section (".cpsb_reg")))  
#define BYTE		unsigned char
#define WORD		unsigned short
#define DWORD		unsigned int
#define BOOL		unsigned int

#define	MAXSPRITES	256

typedef struct
{
	WORD	x;			// Sprite x position
	WORD	y;			// Sprite y position
	WORD	Tile;		// Sprite tile
	WORD	Attribute;	// Sprite attribute
	WORD	Used;		// Is this SpriteBlock used?
} Sprite;

GFXRAM Sprite sprites [MAXSPRITES]  =  {};
CPSA_REG WORD  cpsa_reg[0x19] = {};
CPSB_REG WORD cpsb_reg[0x19] = {};

int vsyncCounter = 0;
int soundCounter = 0;
int someCounter = 0x6666;
int uinVARRRRRR;
const int mYcOnStT = 1;

void onVSync() {
   vsyncCounter++;
   if (vsyncCounter >= 60) {
   	vsyncCounter = 0;
   	soundCounter += 1;
   	*((char*)0x800180) = (0x22 + soundCounter);
   	
   } else {
   	*((char*)0x800180) = 0xFF;
   }
}



int run() {
	while(1) {

	}
	return 0;
}
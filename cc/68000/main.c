
#define	BYTE		unsigned char
#define	WORD		unsigned short
#define	DWORD		unsigned int
#define	BOOL		unsigned int

#define	MAXSPRITES	8

typedef struct
{
	WORD	x;			// Sprite x position
	WORD	y;			// Sprite y position
	WORD	Tile;		// Sprite tile
	WORD	Attribute;	// Sprite attribute
	WORD	Used;		// Is this SpriteBlock used?
} tagSpriteBlock;

tagSpriteBlock	SpriteBlock[MAXSPRITES];

int SpritesChanged = 1;

int vsyncCounter = 0;
int soundCounter = 0;
int someCounter = 0x6666;

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
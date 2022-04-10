#ifndef CCPS_GFX
#define CCPS_GFX

#define BYTE		unsigned char
#define WORD		unsigned short
#define DWORD		unsigned int
#define BOOL		unsigned int

typedef struct {
   short  height;
   short  width;
   short id;
} GFXSprite;

typedef struct {
   short x;
   short y;
   short id;
} GFXShapeTile;

typedef struct {
   short numTiles;
   GFXShapeTile tiles[];
} GFXShape;

#define Color WORD
typedef struct {
  Color colors[16];
} Palette;

#endif // CCPS_GFX
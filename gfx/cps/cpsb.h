
// Registers location and meaning change.
// Add here supported CPSB version and their registers as needed

/*
 CPSB_REG_PALETTE_CONTROL

  bit 0: copy page 0 (sprites)
  bit 1: copy page 1 (scroll1)
  bit 2: copy page 2 (scroll2)
  bit 3: copy page 3 (scroll3)
  bit 4: copy page 4 (stars1)
  bit 5: copy page 5 (stars2)
*/

#if CPSB_VERSION == 01
#define CPSB_REG_LAYER_CTRL      (0x26 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x30 / 2)
#endif // CPSB_VERSION == 01

#if CPSB_VERSION == 02
#define CPSB_REG_LAYER_CTRL      (0x2c / 2)
#define CPSB_REG_PALETTE_CONTROL (0x22 / 2)
#endif // CPSB_VERSION == 02

#if CPSB_VERSION == 03
#define CPSB_REG_LAYER_CTRL      (0x30 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x26 / 2)
#endif // CPSB_VERSION == 03

#if CPSB_VERSION == 04
#define CPSB_REG_LAYER_CTRL      (0x2e / 2)
#define CPSB_REG_LAYER_PRI0      (0x26 / 2)
#define CPSB_REG_LAYER_PRI1      (0x30 / 2)
#define CPSB_REG_LAYER_PRI2      (0x28 / 2)
#define CPSB_REG_LAYER_PRI3      (0x32 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x2a / 2)
#endif // CPSB_VERSION == 04

#if CPSB_VERSION == 05
#define CPSB_REG_LAYER_CTRL      (0x28 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x32 / 2)
#endif // CPSB_VERSION == 05

#if CPSB_VERSION == 11
#define CPSB_REG_LAYER_CTRL      (0x26 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x30 / 2)
#endif // CPSB_VERSION == 11

#if CPSB_VERSION == 12
#define CPSB_REG_LAYER_CTRL      (0x2c / 2)
#define CPSB_REG_PALETTE_CONTROL (0x22 / 2)
#endif // CPSB_VERSION == 12

#if CPSB_VERSION == 13
#define CPSB_REG_LAYER_CTRL      (0x22 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x2c / 2)
#endif // CPSB_VERSION == 13

#if CPSB_VERSION == 14
#define CPSB_REG_LAYER_CTRL      (0x12 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x1c / 2)
#endif // CPSB_VERSION == 14

#if CPSB_VERSION == 15
#define CPSB_REG_LAYER_CTRL      (0x02 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x0c / 2)
#endif // CPSB_VERSION == 15

#if CPSB_VERSION == 16
#define CPSB_REG_LAYER_CTRL      (0x0c / 2)
#define CPSB_REG_PALETTE_CONTROL (0x02 / 2)
#endif // CPSB_VERSION == 16

#if CPSB_VERSION == 17
#define CPSB_REG_LAYER_CTRL      (0x14 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x0a / 2)
#endif // CPSB_VERSION == 17

#if CPSB_VERSION == 18
#define CPSB_REG_LAYER_CTRL      (0x1c / 2)
#define CPSB_REG_PALETTE_CONTROL (0x12 / 2)
#endif // CPSB_VERSION == 18

#if CPSB_VERSION == 21
#define CPSB_REG_LAYER_CTRL      (0x26 / 2)
#define CPSB_REG_PALETTE_CONTROL (0x30 / 2)
#endif // CPSB_VERSION == 21

// Registers location and meaning change.
// Add here supported CPSB version and their registers as needed


#if CPSB_VERSION == 05
// Palette control
// bit 0: copy page 0 (sprites)
// bit 1: copy page 1 (scroll1)
// bit 2: copy page 2 (scroll2)
// bit 3: copy page 3 (scroll3)
// bit 4: copy page 4 (stars1)
// bit 5: copy page 5 (stars2)
#define CPSB_REG_PALETTE_CONTROL (0x32 / 2)

#endif // CPSB_VERSION == 05

#if CPSB_VERSION == 11
// Palette control
// bit 0: copy page 0 (sprites)
// bit 1: copy page 1 (scroll1)
// bit 2: copy page 2 (scroll2)
// bit 3: copy page 3 (scroll3)
// bit 4: copy page 4 (stars1)
// bit 5: copy page 5 (stars2)
#define CPSB_REG_PALETTE_CONTROL (0x30 / 2)

#endif // CPSB_VERSION == 11
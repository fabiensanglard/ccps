module ccps/ccps

go 1.17

replace ccps/z80 => ../z80
replace ccps/m68k => ../m68k
replace ccps/boards => ../boards
replace ccps/ym2151 => ../ym2151
replace ccps/oki => ../oki
replace ccps/gfx => ../gfx
replace ccps/mus => ../mus

require ccps/z80 v0.0.0-00010101000000-000000000000
require ccps/m68k v0.0.0-00010101000000-000000000000
require ccps/boards v0.0.0-00010101000000-000000000000
require ccps/ym2151 v0.0.0-00010101000000-000000000000
require ccps/oki v0.0.0-00010101000000-000000000000
require ccps/gfx v0.0.0-00010101000000-000000000000
require ccps/mus v0.0.0-00010101000000-000000000000

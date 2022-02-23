mkdir -p out
m68k-linux-gnu-as -m68000 --register-prefix-optional -o out/StartUp.o StartUp.s
m68k-linux-gnu-gcc -m68000 -nostdlib -c -O0 -o out/main.o main.c
m68k-linux-gnu-gcc -Llib -m68000 -Wall -nostartfiles -nodefaultlibs -fno-builtin -fomit-frame-pointer -ffast-math -Wl,-Map,out/game.map -Wl,--build-id=none -T CPS1.x -o out/game.obj out/StartUp.o out/main.o
m68k-linux-gnu-objcopy --gap-fill=0xFF --pad-to=1048576  -R .data  --output-target=binary out/game.obj out/game.rom
go run split.go
mv out/sf2e_30g.11e ~/mame/roms/sf2
mv out/sf2e_37g.11f ~/mame/roms/sf2
mv out/sf2e_31g.12e ~/mame/roms/sf2
mv out/sf2e_38g.12f ~/mame/roms/sf2
mv out/sf2e_28g.9e ~/mame/roms/sf2
mv out/sf2e_35g.9f ~/mame/roms/sf2
mv out/sf2_29b.10e ~/mame/roms/sf2
mv out/sf2_36b.10f ~/mame/roms/sf2
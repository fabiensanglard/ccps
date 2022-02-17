m68k-linux-gnu-as -m68000 --register-prefix-optional -o Routines.o Routines.s
m68k-linux-gnu-as -m68000 --register-prefix-optional -o StartUp.o StartUp.s
m68k-linux-gnu-gcc -m68000 -nostdlib -c -O0 -o main.o main.c
m68k-linux-gnu-gcc -Llib -m68000 -Wall -nostartfiles -nodefaultlibs -fno-builtin -fomit-frame-pointer -ffast-math -Wl,-Map,game.map -Wl,--build-id=none -T CPS1.x -o game.obj StartUp.o Routines.o main.o
m68k-linux-gnu-objcopy --gap-fill=0xFF --pad-to=1048576  -R .data  --output-target=binary game.obj game.rom
go run split.go
mv sf2e_30g.11e ~/mame/roms/sf2
mv sf2e_37g.11f ~/mame/roms/sf2
mv sf2e_31g.12e ~/mame/roms/sf2
mv sf2e_38g.12f ~/mame/roms/sf2
mv sf2e_28g.9e ~/mame/roms/sf2
mv sf2e_35g.9f ~/mame/roms/sf2
mv sf2_29b.10e ~/mame/roms/sf2
mv sf2_36b.10f ~/mame/roms/sf2
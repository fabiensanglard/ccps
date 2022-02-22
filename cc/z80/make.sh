sdasz80 -plogff -o crt0.rel crt0.s
sdcc --compile-only -mz80 --data-loc 0xd000 --no-std-crt0 main.c 
sdldz80 -nf main.lk
objcopy --input-target=ihex --output-target=binary main.ihx main.rom
dd if=/dev/zero of=main.rom bs=1 count=1 seek=65536
cp main.rom ~/mame/roms/sf2/sf2_9.12a
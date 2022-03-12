rm -fr out
mkdir -p out
cd cc/z80
./make.sh
cd ../..
cd cc/68000
./make.sh
cd ../..
cp out/* ~/mame/roms/sf2

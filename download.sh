#!/bin/sh -e

archive=df_40_09_linux.tar.bz2

# make core dumps if something crashes.
ulimit -c unlimited

# delete the old work directory and make a new one with DF2014 in it.
if [ ! -e "work/$archive" ]
then
	rm -rf work
	mkdir work
	cd work
	wget "http://www.bay12games.com/dwarves/$archive"
	tar xvf "$archive"
else
	cd work
fi

cd df_linux

# disable graphics and sound so it works on headless systems.
sed -e 's/\[PRINT_MODE:2D\]/[PRINT_MODE:TEXT]/' -i 'data/init/init.txt'
sed -e 's/\[SOUND:YES\]/[SOUND:NO]/'            -i 'data/init/init.txt'

# make the max year lower so the generation process doesn't take as long.
sed -e 's/\[END_YEAR:1050\]/[END_YEAR:250]/'    -i 'data/init/world_gen.txt'

# generate some worlds!
./libs/Dwarf_Fortress -gen 1 12345 "LARGE REGION"  || true
./libs/Dwarf_Fortress -gen 2 54321 "POCKET ISLAND" || true

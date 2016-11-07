#!/bin/bash -e

exec > generate.go

echo "//go:generate bash generate.bash > generate.go"
echo
echo "package namegen"

language_file() {
	echo
	echo "const language$1 = \`"
	curl -sSL "https://raw.githubusercontent.com/BenLubar/raws/v0.42.06/objects/language_$1.txt"
	echo "\`"
}

language_file "words"
language_file "DWARF"
language_file "HUMAN"
language_file "GOBLIN"
language_file "ELF"

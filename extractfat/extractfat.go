package main

// extractfat <fat_file>

import (
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <fat_file>\n", os.Args[0])
		os.Exit(2)
	}

	file := os.Args[1]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if len(data) < 24 {
		panic(fmt.Sprintf("file %s too small", file))
	}
	// fat binary header is BE
	magicBE := binary.BigEndian.Uint32(data[0:4])
	if magicBE != macho.MagicFat {
		panic(fmt.Sprintf("input %s is not a macho file, magic=%x, expected=%x", file, magicBE, macho.MagicFat))
	}

	mf, err := macho.OpenFat(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer mf.Close()

	fmt.Printf("file contains %d arches\n", len(mf.Arches))

	for _, arch := range mf.Arches {
		fah := arch.FatArchHeader
		fmt.Printf("- %#v: ", fah)
		filename := os.Args[1] + "." + fah.Cpu.String()
		content := data[fah.Offset : fah.Offset+fah.Size]
		fmt.Printf("writing %s\n", filename)
		// #nosec
		if err = ioutil.WriteFile(filename, content, 0o755); err != nil {
			panic(err)
		}
	}
}

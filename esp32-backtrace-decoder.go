package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("unable to determine home directory: ", err)
		os.Exit(-1)
	}

	// Default path to addr2line
	addr2lineDefault := path.Join(home, ".platformio/packages/toolchain-xtensa32/bin/xtensa-esp32-elf-addr2line")

	var addr2line = flag.String("addr2line", addr2lineDefault, "path to xtensa-esp32-elf-addr2line")
	var elfFile = flag.String("elf", "", "path to program elf file")

	flag.Parse()

	// Validate the flags
	if *elfFile == "" {
		fmt.Println("must specify elf file")
		os.Exit(-1)
	}
	if _, err := os.Lstat(*elfFile); err != nil {
		fmt.Printf("unable to find %v: %v\n", *elfFile, err)
		os.Exit(-1)
	}
	if _, err := os.Lstat(*addr2line); err != nil {
		fmt.Printf("unable to find addr2line: %v\n", err)
		os.Exit(-1)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("must provide backtrace")
		os.Exit(-1)
	}

	backtrace := strings.Join(flag.Args(), " ")

	addresses := regexp.MustCompile(`0x[0-9a-f]{8}:0x[0-9a-f]{8}`).FindAllString(backtrace, -1)
	for _, address := range addresses {
		cmd := exec.Command(*addr2line, "-pfiaC", "-e", *elfFile, address)
		if out, err := cmd.Output(); err != nil {
			fmt.Println("failed to run addr2line: ", err)
			os.Exit(-1)
		} else {
			fmt.Println(strings.TrimSpace(string(out)))
		}
	}
}

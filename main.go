package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vbatts/power_info/helper"
	"github.com/vbatts/power_info/linux"
	"os"
	"path/filepath"
	"time"
)

var (
	quiet bool
)

func main() {
	flag.Parse()
	linux.SetQuiet(quiet)

	powers, err := filepath.Glob(linux.PowerSupplyPath + "*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(2)
	}
	for _, power := range powers {
		info := linux.Info{
			Key:     filepath.Base(power),
			Time:    time.Now().UnixNano(),
			Values:  map[string]string{},
			Load:    linux.GetLoadAvg(),
			Version: linux.GetVersion(),
		}

		files, err := filepath.Glob(power + "/*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		for _, file := range files {
			basename := filepath.Base(file)
			if helper.IsFile(file) {
				str, err := helper.StringFromFile(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: reading file [%s]: %s\n", file, err)
				} else {
					info.Values[basename] = str
				}
			} else {
				if !quiet {
					fmt.Fprintf(os.Stderr, "WARN: [%s] does not appear to be a file\n", file)
				}
			}
		}
		buf, err := json.Marshal(info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		fmt.Printf("%s\n", buf)
	}
}

func init() {
	flag.BoolVar(&quiet, "quiet", false, "less output")
}

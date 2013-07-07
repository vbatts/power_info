package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vbatts/power_info/linux"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
			if linux.IsFile(file) {
				fh, err := os.Open(file)
				if err != nil {
					if !quiet {
						fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
					}
					continue
				}
				b, err := ioutil.ReadAll(fh)
				if err != nil {
					if !quiet {
						fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
					}
					continue
				}
				info.Values[basename] = strings.TrimRight(bytes.NewBuffer(b).String(), " \n")
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

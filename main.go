package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ParseToCsv(filename string, output io.Writer) (err error) {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	b_reader := bufio.NewReader(fh)
	c_writer := csv.NewWriter(output)
  c_writer.Write([]string{
    "Epoch Nano","Time","Key","Version",
  })
  _ = b_reader

	return nil
}

func main() {
	flag.Parse()

	if len(fileToParseToCSV) > 0 {
		os.Exit(0)
	}

	powers, err := filepath.Glob(POWER + "*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(2)
	}
	for _, power := range powers {
		info := Info{
			Key:     filepath.Base(power),
			Time:    time.Now().UnixNano(),
			Values:  map[string]string{},
			Load:    GetLoadAvg(),
			Version: GetVersion(),
		}

		files, err := filepath.Glob(power + "/*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		for _, file := range files {
			basename := filepath.Base(file)
			if IsFile(file) {
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
	flag.StringVar(&fileToParseToCSV, "parse", "", "parse this json.log file and render it to a CSV (for spreadsheeting)")
}

var (
	POWER            = "/sys/class/power_supply/"
	LOAD_AVG         = "/proc/loadavg"
	VERSION          = "/proc/version"
	quiet            bool
	fileToParseToCSV string
)

// a representation of /proc/loadavg, leaving of the PID of the last process
type LoadAvg struct {
	Avg1, Avg5, Avg15    string
	Schedulers, Entities string
}

// get the current LoadAvg
func GetLoadAvg() LoadAvg {
	fh, err := os.Open(LOAD_AVG)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return LoadAvg{}
	}
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return LoadAvg{}
	}
	values := strings.Split(strings.TrimRight(bytes.NewBuffer(b).String(), " \n"), " ")
	sAndE := strings.Split(values[3], "/")
	return LoadAvg{
		values[0],
		values[1],
		values[2],
		sAndE[0],
		sAndE[1],
	}
}

// Version of the current running kernel, /proc/version
func GetVersion() string {
	fh, err := os.Open(VERSION)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return ""
	}
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return ""
	}
	return strings.TrimRight(bytes.NewBuffer(b).String(), " \n")
}

/*
Info set on a /sys/class/power_supply item

Time is UnixNano
Values are the file [name]contents
Load is the LoadAvg when that Info was collected
*/
type Info struct {
	Time    int64
	Key     string
	Values  map[string]string
	Load    LoadAvg `json:",omitempty"`
	Version string
}

// Convenience Method for checking files
func IsFile(filename string) bool {
	if fi, _ := os.Stat(filename); fi.Mode().IsRegular() {
		return true
	}
	return false
}

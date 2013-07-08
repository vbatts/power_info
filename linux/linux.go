package linux

import (
	"fmt"
	"github.com/vbatts/power_info/helper"
	"os"
	"strings"
)

var (
	PowerSupplyPath = "/sys/class/power_supply/"
	LoadAvgPath     = "/proc/loadavg"
	VersionPath     = "/proc/version"
	quiet           bool
)

// set quiet: true or false
func SetQuiet(state bool) {
	quiet = state
}

// a representation of /proc/loadavg, leaving of the PID of the last process
type LoadAvg struct {
	Avg1, Avg5, Avg15    string
	Schedulers, Entities string
}

// get the current LoadAvg
func GetLoadAvg() LoadAvg {
	str, err := helper.StringFromFile(LoadAvgPath)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return LoadAvg{}
	}
	values := strings.Split(str, " ")
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
	str, err := helper.StringFromFile(VersionPath)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "WARN: %s\n", err)
		}
		return ""
	}
	return str
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

package linux

import (
	"fmt"
	"github.com/vbatts/power_info/helper"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
Abstract interface for what constitutes a power supply

Not sure what to do with this yet.
*/
type PowerSupply interface {
	GetInfo() (*Info, error)
}

func NewGenericPowerSupply(key string) GenericPowerSupply {
	return GenericPowerSupply{Key: key}

}

type GenericPowerSupply struct {
	Key string
}

func (g *GenericPowerSupply) GetInfo() (*Info, error) {
	info := Info{
		Key:     g.Key,
		Time:    time.Now().UnixNano(),
		Values:  map[string]string{},
		Load:    GetLoadAvg(),
		Version: GetVersion(),
	}
	files, err := filepath.Glob(filepath.Join(PowerSupplyPath, g.Key, "*"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if helper.IsFile(file) {
			if str, err := helper.StringFromFile(file); err == nil {
				info.Values[filepath.Base(file)] = str
			}
		}
	}

	return &info, nil
}

func GetPowerSupplies() (pss []GenericPowerSupply) {
	powers, err := filepath.Glob(PowerSupplyPath + "*")
	if err != nil {
		return
	}
	for _, power := range powers {
		pss = append(pss, NewGenericPowerSupply(filepath.Base(power)))
	}

	return
}

/*
Given a key like "AC" or "BAT0", determine whether it is a battery type
*/
func IsBattery(key string) bool {
	type_file := filepath.Join(PowerSupplyPath, key, "type")
	if str, err := helper.StringFromFile(type_file); err == nil && strings.Contains(str, "Battery") {
		return true
	}
	return false
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

func NewBattery(key string) Battery {
	return Battery{Key: key}
}

type Battery struct {
	Key string
}

func (b *Battery) Status() string {
	str, _ := helper.StringFromFile(psp(b.Key, "status"))
	return str
}

func (b *Battery) ChargeNow() int64 {
	var str string
	if helper.IsFile(psp(b.Key, "energy_now")) {
		str, _ = helper.StringFromFile(psp(b.Key, "energy_now"))
	} else {
		str, _ = helper.StringFromFile(psp(b.Key, "charge_now"))
	}
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (b *Battery) ChargeFull() int64 {
	var str string
	if helper.IsFile(psp(b.Key, "energy_full")) {
		str, _ = helper.StringFromFile(psp(b.Key, "energy_full"))
	} else {
		str, _ = helper.StringFromFile(psp(b.Key, "charge_full"))
	}
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (b *Battery) Percent() float64 {
	return float64(100) * float64(b.ChargeNow()) / float64(b.ChargeFull())
}

func (b *Battery) GetInfo() (*Info, error) {
	info := Info{
		Key:     b.Key,
		Time:    time.Now().UnixNano(),
		Values:  map[string]string{},
		Load:    GetLoadAvg(),
		Version: GetVersion(),
	}
	files, err := filepath.Glob(psp(b.Key, "*"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if helper.IsFile(file) {
			if str, err := helper.StringFromFile(file); err == nil {
				info.Values[filepath.Base(file)] = str
			}
		}
	}

	return &info, nil
}

// given the Battery.Key, get the path for this power_supply
func psp(key, value string) string {
	if len(value) == 0 {
		return filepath.Join(PowerSupplyPath, key)
	}
	return filepath.Join(PowerSupplyPath, key, value)
}

// Get an Array of the Batteries available on this system
func GetBatteries() (batts []Battery, err error) {
	possibilities, err := filepath.Glob(PowerSupplyPath + "*/type")
	if err != nil {
		return batts, err
	}
	for _, poss := range possibilities {
		key := filepath.Base(filepath.Dir(poss))
		if IsBattery(key) {
			batts = append(batts, Battery{Key: key})
		}
	}

	return batts, nil
}

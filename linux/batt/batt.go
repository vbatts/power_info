package batt

import (
	"github.com/vbatts/power_info/linux"
	"path/filepath"
	"strings"
  "strconv"
)

type Battery struct {
	Key string
}

func (b *Battery) Status() string {
  str, _ := linux.StringFromFile(psp(b.Key, "status"))
  return str
}

func (b *Battery) ChargeNow() int64 {
  var str string
  if linux.IsFile(psp(b.Key, "energy_now")) {
    str, _ = linux.StringFromFile(psp(b.Key, "energy_now"))
  } else {
    str, _ = linux.StringFromFile(psp(b.Key, "charge_now"))
  }
  i, _ := strconv.ParseInt(str, 10, 64)
  return i
}

func (b *Battery) ChargeFull() int64 {
  var str string
  if linux.IsFile(psp(b.Key, "energy_full")) {
    str, _ = linux.StringFromFile(psp(b.Key, "energy_full"))
  } else {
    str, _ = linux.StringFromFile(psp(b.Key, "charge_full"))
  }
  i, _ := strconv.ParseInt(str, 10, 64)
  return i
}

func (b *Battery) Percent() float64 {
  return float64(b.ChargeNow()) / float64(b.ChargeFull())
}

// given the Battery.Key, get the path for this power_supply
func psp(key, value string) string {
	if len(value) == 0 {
		return filepath.Join(linux.PowerSupplyPath, key)
	}
	return filepath.Join(linux.PowerSupplyPath, key, value)
}

// Get an Array of the Batteries available on this system
func GetBatteries() (batts []Battery, err error) {
	possibilities, err := filepath.Glob(linux.PowerSupplyPath + "*/type")
	if err != nil {
		return batts, err
	}
	for _, poss := range possibilities {
		str, _ := linux.StringFromFile(poss)
		if strings.Contains(str, "Battery") {
			key := filepath.Base(filepath.Dir(poss))
			batts = append(batts, Battery{Key: key})
		}
	}

	return batts, nil
}

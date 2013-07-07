package batt

import (
	"github.com/vbatts/power_info/linux"
	"path/filepath"
	"strings"
)

type Battery struct {
	Key string
}

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

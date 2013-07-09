package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vbatts/power_info/linux"
	"os"
	"strings"
)

var (
	quiet   bool
	percent bool
)

func main() {
	flag.Parse()
	linux.SetQuiet(quiet)

	if percent {
		batts, err := linux.GetBatteries()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		switch len(batts) {
		case 0:
			fmt.Fprintf(os.Stderr, "No batteries present\n", err)
			os.Exit(2)
		case 1:
			fmt.Printf("%s: %3.2f%% (%s)\n",
				batts[0].Key,
				batts[0].Percent(),
				batts[0].Status())
		default:
			var (
				charge_sum int64 = 0
				full_sum   int64 = 0
				batt_strs  []string
			)
			for _, batt := range batts {
				charge_sum = charge_sum + batt.ChargeNow()
				full_sum = full_sum + batt.ChargeFull()
				batt_strs = append(batt_strs, batt.Key)
			}
			fmt.Printf("%s: %3.2f%% \n",
				strings.Join(batt_strs, ", "),
				(float64(100) * float64(charge_sum) / float64(full_sum)))
		}
		os.Exit(0)
	}

	for _, power := range linux.GetPowerSupplies() {
		info, err := power.GetInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		buf, err := json.Marshal(*info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}
		fmt.Printf("%s\n", buf)
	}
}

func init() {
	flag.BoolVar(&quiet, "quiet", false, "less output")
	flag.BoolVar(&percent, "percent", false, "show battery percent and exit")
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vbatts/power_info/linux"
	"os"
	"strings"
	"time"
)

var (
	quiet           bool
	battery_rate    bool
	battery_percent bool
)

func main() {
	flag.Parse()
	linux.SetQuiet(quiet)

	// print battery percentage, and exit
	if battery_percent {
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
			if quiet {
				fmt.Printf("%3.2f%%[%s]\n",
					batts[0].Percent(),
					string(batts[0].Status()[0]))
			} else {
				fmt.Printf("%s: %3.2f%% (%s)\n",
					batts[0].Key,
					batts[0].Percent(),
					batts[0].Status())
			}
		default:
			var batt_strs []string
			for _, batt := range batts {
				batt_strs = append(batt_strs, batt.Key)
			}
			fmt.Printf("%s: %3.2f%% \n",
				strings.Join(batt_strs, ", "),
				linux.Percent(batts))
		}
		os.Exit(0)
	}

	/*
	  roll with the printing of battery rate used

	  Big and nasty. Perhaps this could be bottled into the Battery struct.
	  such that you could see the rate and duration of each battery independently?
	*/
	if battery_rate {
		var (
			c_curr_charge_ema = make(chan float64)
			c_rate            = make(chan float64)
		)
		batts, err := linux.GetBatteries()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(2)
		}

		// subroutine to get the current avg and last diff
		go func() {
			var (
				count                            int64
				diff                             int64
				prev_charge, curr_charge         int64
				prev_time, curr_time             int64
				prev_charge_ema, curr_charge_ema float64 = 1, 1
				prev_rate_ema, curr_rate_ema     float64 = 1, 1
			)
			prev_time = time.Now().UnixNano()

			for {
				curr_charge = linux.ChargeNow(batts)
				if curr_charge == prev_charge {
					time.Sleep(333 * time.Millisecond)
					continue
				}
				count++

				diff = (curr_charge - prev_charge)
				prev_charge = curr_charge

				// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
				curr_charge_ema = A(count)*float64(diff) + (1-A(count))*prev_charge_ema
				prev_charge_ema = curr_charge_ema

				curr_time = time.Now().UnixNano()
				rate := float64(diff) / float64(curr_time-prev_time)
				curr_rate_ema = A(count)*rate + (1-A(count))*prev_rate_ema
				prev_rate_ema = curr_rate_ema
				prev_time = curr_time
				fmt.Println(rate)
				fmt.Println(curr_rate_ema)

				c_rate <- curr_rate_ema
				c_curr_charge_ema <- curr_charge_ema

				time.Sleep(500 * time.Millisecond)
			}
		}()

		for {
			duration := (<-c_rate * <-c_curr_charge_ema)
			fmt.Printf("%3.2f%% est: %3.2fs\n",
				linux.Percent(batts), duration)
		}

		os.Exit(0) // should never make it here
	}

	// default is to print power_supply info and quit
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

// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
func A(n int64) float64 {
	return float64(2) / float64(n+1)
}

func init() {
	flag.BoolVar(&quiet, "quiet", false, "less output")
	flag.BoolVar(&battery_percent, "percent", false, "show battery percent and exit")
	flag.BoolVar(&battery_rate, "rate", false, "show battery ")
}

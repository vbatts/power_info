package batt

import (
	"fmt"
	"testing"
)

func TestFuncs(t *testing.T) {
	str := psp("BAT0", "")
	fmt.Println(str)

	str = psp("BAT0", "status")
	fmt.Println(str)

	if IsBattery("AC") {
		t.Errorf("AC should not be seen as a battery")
	}

}

// figure out how to test this
func TestBatts(t *testing.T) {
	batts, err := GetBatteries()
	if err != nil {
		t.Error(err)
	}
	for _, batt := range batts {
		fmt.Println(batt.Status())
		fmt.Println(batt.ChargeNow())
		fmt.Println(batt.ChargeFull())
		fmt.Println(batt.Percent())
		fmt.Println(batt.GetInfo())
	}

}

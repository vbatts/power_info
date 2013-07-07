package batt

import (
  "testing"
  "fmt"
)

func TestFuncs(t *testing.T) {
  str := psp("BAT0","")
  fmt.Println(str)

  str = psp("BAT0","status")
  fmt.Println(str)

}

func TestBatts(t *testing.T) {
  batts ,err := GetBatteries()
  if err != nil {
    t.Error(err)
  }
  for _, batt := range batts {
    fmt.Println(batt.Status())
    fmt.Println(batt.ChargeNow())
    fmt.Println(batt.ChargeFull())
    fmt.Println(batt.Percent())
  }

}

package batt

import (
  "testing"
  "fmt"
)

func TestBatts(t *testing.T) {
  batts ,err := GetBatteries()
  if err != nil {
    t.Error(err)
  }
  fmt.Println(batts)
}

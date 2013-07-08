package helper

import (
  "testing"
  "strings"
)

func TestHelper(t *testing.T) {
  filename := "helper_test.go"
  if !IsFile(filename) {
    t.Errorf("this is a file!")
  }

  str, err := StringFromFile(filename)
  if err != nil {
    t.Error(err)
  }
  if !strings.Contains(str, "package helper") {
    t.Errorf("this file ought to contain 'package helper'")
  }

}


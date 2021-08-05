package main

import "testing"

func TestValidSymbol(t *testing.T) {
	var symbolsInvalid = []string {"391231", "598754", "424594", "311564", "611234", "011654", "198756", "335489", "695151"}
	for _, symbol := range symbolsInvalid {
		isValid := isValidSymbol(symbol)
		if isValid != false {
			t.Fail()
			t.Logf("symbol: %v, excepted: %v, actual: %v", symbol, false, isValid)
		}
	}

	var symbolsValid = []string { "300000", "301234", "600000", "601234", "000000", "001235", "685464"}
	for _, symbol := range symbolsValid {
		isValid := isValidSymbol(symbol)
		if isValid != true {
			t.Fail()
			t.Logf("symbol: %v, excepted: %v, actual: %v", symbol, true, isValid)
		}
	}
}
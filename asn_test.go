package asnbuilder

import (
	numberRange "./numberRange"
	"net/http"
	"regexp"
	"testing"
)

func TestRangeGenerator(t *testing.T) {
	if numberRange.GetRegex(400, 569) != "_4[0-9][0-9]|_5[0-6][0-9]" {
		t.Error("Failed generating correct regex")
	}
}

func TestPostiveMatch(t *testing.T) {
	/* test  for my home AS196922 , Hofmeir Media GmbH Datacenter  */
	/* positive match */
	regex := numberRange.GetRegex(196608, 197631)
	positiveMatch, err := regexp.Compile(regex)
	if err != nil {
		t.Error("Cant compile regex")
	}

	if !positiveMatch.MatchString("_196922") {
		t.Error("Regex did not match _196922!")
	}
}

func TestNegativeMatch(t *testing.T) {
	/* negative match */
	regex := numberRange.GetRegex(264605, 265628)
	negativeMatch, err := regexp.Compile(regex)
	if err != nil {
		t.Error("Cant compile regex")
	}

	if negativeMatch.MatchString("_196922") {
		t.Error("Regex did match 196922!")
	}
}

func TestRechabilityIANA(t *testing.T) {
	_, err := http.Get("http://www.iana.org")
	if err != nil {
		t.Error("Cant connect to iana.org")
	}

}

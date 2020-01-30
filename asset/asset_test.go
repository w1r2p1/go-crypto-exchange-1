package asset

import (
	"fmt"
	"math/big"
	"regexp"
	"testing"
)

type testData struct {
	Value     interface{}
	Iscorrect bool
}

func TestNew(t *testing.T) {
	_, err := New(*big.NewInt(0), 15, "TEST")
	if err == nil {
		t.Errorf("Expected zero value precision error")
	}

	_, err = New(*big.NewInt(-3456467865865), 13, "TEST")
	if err != nil {
		t.Errorf("Expected success. Got error: %v", err)
	}

	_, err = New(*big.NewInt(3456467865865), 20, "TEST")
	if err == nil {
		t.Errorf("Expected precision overflow error")
	}

	_, err = New(*big.NewInt(3456467865865), 13, "TEST")
	if err != nil {
		t.Errorf("Expected success. Got error: %v", err)
	}
}

func TestRegexp(t *testing.T) {
	var testStrings = [...]testData{
		{"", false}, {"31.41e10", false}, {"3.", false}, {".14", false},
		{".", false}, {"3.1.4", false}, {".1.4", false}, {"3.-14", false},

		{"3.14", true}, {"3", true}, {"33.14", true}, {"0.999999999999999999", true},
		{"-3.14", true}, {"999999999999999999.0", true}, {"0003.0", true},
		{"+3.14", true}, {"03.14", true},
	}

	reg, err := regexp.Compile(FloatValidationRegexp)

	if err != nil {
		t.Errorf("Regexp compilation error: %v", err)
	}

	for _, element := range testStrings {
		match := reg.MatchString(fmt.Sprintf("%v", element.Value))
		if match != element.Iscorrect {
			t.Errorf("Regexp validation error in value: \"%s\"\n Expected: %v. Got: %v",
				element.Value, element.Iscorrect, match)
		}
	}
}

func TestNewNatural(t *testing.T) {
	var testCurrencies = [...]testData{
		{"3.14", true}, {"3.143456", true}, {"3", true},
		{"-3.14", true}, {"35345345.133", true}, {"0", true},
		{"+3.14", true}, {"03.14", true},

		{"0.00000", false},
	}

	for _, element := range testCurrencies {
		_, err := NewNatural(element.Value.(string), "TEST")
		if err != nil {
			if element.Iscorrect {
				t.Errorf("Somethig went wrong: %v", err)
			}
		}
	}

}

func TestNatural(t *testing.T) {
	res, _ := NewNatural("14.0003", "TEST")
	if res.Natural() != "14.0003" {
		t.Errorf("Wanted: \"14.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = New(*big.NewInt(140003), 4, "TEST")
	if res.Natural() != "14.0003" {
		t.Errorf("Wanted: \"14.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("0.0003", "TEST")
	if res.Natural() != "0.0003" {
		t.Errorf("Wanted: \"0.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = New(*big.NewInt(3), 4, "TEST")
	if res.Natural() != "0.0003" {
		t.Errorf("Wanted: \"0.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("0", "TEST")
	if res.Natural() != "0" {
		t.Errorf("Wanted: \"0\". Got: \"%s\"", res.Natural())
	}

	res, _ = New(*big.NewInt(0), 0, "TEST")
	if res.Natural() != "0" {
		t.Errorf("Wanted: \"0\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("-0.0003", "TEST")
	if res.Natural() != "-0.0003" {
		t.Errorf("Wanted: \"-0.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = New(*big.NewInt(-3), 4, "TEST")
	if res.Natural() != "-0.0003" {
		t.Errorf("Wanted: \"-0.0003\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("314", "ASD")
	if res.Natural() != "314" {
		t.Errorf("Wanted: \"314\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("3463734634", "ASD")
	if res.Natural() != "3463734634" {
		t.Errorf("Wanted: \"3463734634\". Got: \"%s\"", res.Natural())
	}

	res, _ = NewNatural("-1", "ASD")
	if res.Natural() != "-1" {
		t.Errorf("Wanted: \"-1\". Got: \"%s\"", res.Natural())
	}
}

func TestCmp(t *testing.T) {
	asset1, _ := New(*big.NewInt(314), 2, "TEST")
	asset2, _ := NewNatural("3.14", "test")

	// want true
	if !asset1.Cmp(asset2) {
		fmt.Println(asset1)
		fmt.Println(asset2)
		t.Errorf("Something wrong:")
	}
	asset2, _ = New(*big.NewInt(314), 2, "Error")

	// want false
	if asset1.Cmp(asset2) {
		fmt.Println(asset1)
		fmt.Println(asset2)
		t.Errorf("Something wrong:")
	}

	asset1, _ = New(*big.NewInt(0), 0, "TEST")
	asset2, _ = NewNatural("0", "test")

	// want true
	if !asset1.Cmp(asset2) {
		fmt.Println(asset1)
		fmt.Println(asset2)
		t.Errorf("Something wrong:")
	}

	asset1, _ = New(*big.NewInt(-314), 2, "TEST")
	asset2, _ = NewNatural("3.14", "test")

	// want false
	if asset1.Cmp(asset2) {
		fmt.Println(asset1)
		fmt.Println(asset2)
		t.Errorf("Something wrong:")
	}

	asset1, _ = New(*big.NewInt(314), 2, "TEST")
	asset2, _ = NewNatural("03.14", "test")

	// want true
	if !asset1.Cmp(asset2) {
		fmt.Println(asset1)
		fmt.Println(asset2)
		t.Errorf("Something wrong:")
	}
}

func TestAdd(t *testing.T) {
	x, _ := NewNatural("3.14", "Test")
	y, _ := NewNatural("3.14", "Test")
	res, _ := Add(x, y)
	if res.Natural() != "6.28" {
		t.Errorf("Want: 6.28. Got: %s", res.Natural())
	}

	x, _ = NewNatural("3.14", "Test")
	y, _ = NewNatural("-3.14", "Test")
	res, _ = Add(x, y)
	if res.Natural() != "0" {
		t.Errorf("Want: 0. Got: %s", res.Natural())
	}
}

func TestSub(t *testing.T) {
	x, _ := NewNatural("3.14", "Test")
	y, _ := NewNatural("3.14", "Test")
	res, _ := Sub(x, y)
	if res.Natural() != "0" {
		t.Errorf("Want: 0. Got: %s", res.Natural())
	}

	x, _ = NewNatural("3.14", "Test")
	y, _ = NewNatural("-3.14", "Test")
	res, _ = Sub(x, y)
	if res.Natural() != "6.28" {
		t.Errorf("Want: 6.28. Got: %s", res.Natural())
	}
}

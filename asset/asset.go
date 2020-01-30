package asset

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// MaxAssetPrecision is set up to 18 to permit coin have at least 1 digit in
// integer position and 0-18 in fractional position.
//
// For example: 0.999999999999999999.
const MaxAssetPrecision int = 18

// FloatValidationRegexp is used to validate floa numbers'
const FloatValidationRegexp string = "^[+-]?[0-9]+([.]+[0-9]+)?$"

// MaxFloatStringSize is the maximum size permitted for a float string.
// That wil guard from int64 overflow
const MaxFloatStringSize int = 18

// Make an instance of Asset woth zero values initialization.
func Make() *Asset {
	return &Asset{*big.NewInt(0), 0, ""}
}

// New creates and returns an instance of asset type.
func New(amount big.Int, precision int, symbol string) (*Asset, error) {
	if (amount.Int64() == int64(0)) && precision != 0 {
		return Make(),
			fmt.Errorf("Zero value amount cannot have precision: %d", precision)
	}

	if precision > MaxAssetPrecision {
		return Make(),
			fmt.Errorf("Precision overflow. Max value: %v, got: %v", MaxAssetPrecision, precision)
	}

	symbol = strings.ToUpper(symbol)

	return &Asset{amount, precision, symbol}, nil
}

// NewNatural creates and returns new asset from float representation.
//
// For example: "1.1" .
func NewNatural(amount, symbol string) (*Asset, error) {
	if len(amount) > MaxFloatStringSize {
		return Make(),
			fmt.Errorf("Amount length is too big. Max length: %v. Got: %v",
				MaxFloatStringSize, len(amount))
	}

	reg, err := regexp.Compile(FloatValidationRegexp)
	if err != nil {
		return Make(),
			fmt.Errorf("Regexp compilation error: %v", err)

	}

	match := reg.MatchString(amount)
	if !match {
		return Make(),
			fmt.Errorf("Float validation error: %s", amount)

	}

	precision := 0
	pointIdx := strings.IndexRune(amount, '.')
	if pointIdx > 0 {
		precision = len(amount) - pointIdx - 1
	}

	amount = strings.Replace(amount, ".", "", -1)
	bigAmount := new(big.Int)
	bigAmount.SetString(amount, 10)

	if bigAmount.Int64() == int64(0) && precision > 0 {
		return Make(),
			fmt.Errorf("Zero value amount cannot have precision: %d", precision)
	}

	symbol = strings.ToUpper(symbol)

	return &Asset{*bigAmount, precision, symbol}, nil
}

// Asset is a type to unify basic data and operations to work with
// fiat currencies like USD, RUB, etc...
// and cryptocurrencies like BTC, ETH, EOS, MNT and others.
type Asset struct {
	amount    big.Int //must be stored as significand.
	precision int     //must be at least 1.
	symbol    string
}

// Amount return asset amount as significand.
func (x *Asset) Amount() int64 {
	return x.amount.Int64()
}

// Natural return asset amount in natural representation.
func (x *Asset) Natural() string {
	bigAmount := big.NewInt(0).Abs(&x.amount)
	amount := bigAmount.Text(10)
	result := ""
	pointIdx := len(amount) - int(x.precision)

	if x.amount.Sign() < 0 {
		result = "-" + result
	}

	if x.precision == 0 {
		return result + amount
	}

	if x.amount.Int64() != int64(0) {
		if pointIdx > 0 {
			result = result + amount[:pointIdx] + "." + amount[pointIdx:]
		} else {
			result = result + "0." + strings.Repeat("0", x.precision-1) + amount
		}
	}

	return result
}

// Precision return asset precision.
func (x *Asset) Precision() int {
	return x.precision
}

// Symbol return asset symbol.
func (x *Asset) Symbol() string {
	return x.symbol
}

// Cmp compares x and y by amount, precision and symbol.
// 	Returns	true if they are equal and false if not
func (x *Asset) Cmp(y *Asset) bool {
	amount := x.amount.Cmp(&y.amount) == 0
	precision := x.precision == y.precision
	symbol := x.symbol == y.symbol

	return amount && precision && symbol
}

// Add y to x and return x
// x and y must have same symbol and precision
func Add(x, y *Asset) (*Asset, error) {
	if x.precision != y.precision {
		return nil,
			fmt.Errorf("Precision mismatch error: %v and %v",
				x.precision, y.precision)
	}

	if x.symbol != y.symbol {
		return nil,
			fmt.Errorf("Symbol mismatch error: %v and %v",
				x.symbol, y.symbol)
	}

	result := x
	z := x
	result.amount = *big.NewInt(0).Add(&z.amount, &y.amount)
	if result.Amount() == 0 {
		result.precision = 0
	}

	return result, nil
}

// Sub y from x and return x
// x and y must have same symbol and precision
func Sub(x, y *Asset) (*Asset, error) {
	if x.precision != y.precision {
		return nil,
			fmt.Errorf("Precision mismatch error: %v and %v",
				x.precision, y.precision)
	}

	if x.symbol != y.symbol {
		return nil,
			fmt.Errorf("Symbol mismatch error: %v and %v",
				x.symbol, y.symbol)
	}

	result := x
	z := x
	result.amount = *big.NewInt(0).Sub(&z.amount, &y.amount)
	if result.Amount() == 0 {
		result.precision = 0
	}

	return result, nil
}

// Mul multiplies asset by m
func Mul(x *Asset, m int64) (*Asset, error) {

	result := x
	z := x
	result.amount = *big.NewInt(0).Mul(&z.amount, big.NewInt(m))
	if result.Amount() == 0 {
		result.precision = 0
	}

	return result, nil
}

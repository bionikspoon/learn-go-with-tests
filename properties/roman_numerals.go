package properties

import (
	"fmt"
	"strings"
)

type RomanNumeral struct {
	Value  uint16
	Symbol string
}

type RomanNumerals []RomanNumeral

var romanNumerals = RomanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func (romanNumerals RomanNumerals) ValueOf(symbols ...byte) uint16 {
	symbol := string(symbols)

	for _, romanNumeral := range romanNumerals {
		if romanNumeral.Symbol == symbol {
			return romanNumeral.Value
		}
	}

	panic(fmt.Sprintf("unexpected symbol %#v", symbol))
}

func (romanNumerals RomanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)
	for _, romanNumeral := range romanNumerals {
		if romanNumeral.Symbol == symbol {
			return true
		}

	}
	return false

}

func ConvertToRoman(arabic uint16) string {
	var result strings.Builder

	for _, numeral := range romanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}

func ConvertToArabic(roman string) (total uint16) {

	for _, symbols := range windowedRoman(roman).Symbols() {
		total += romanNumerals.ValueOf(symbols...)
	}

	return

}

type windowedRoman string

func (roman windowedRoman) Symbols() (symbols [][]byte) {
	for index := 0; index < len(roman); index++ {
		symbol := roman[index]
		atEnd := index+1 >= len(roman)

		if !atEnd && isSubtractive(symbol) && romanNumerals.Exists(symbol, roman[index+1]) {
			symbols = append(symbols, []byte{byte(symbol), byte(roman[index+1])})
			index++
		} else {
			symbols = append(symbols, []byte{byte(symbol)})
		}
	}
	return
}

func isSubtractive(symbol uint8) bool {
	return symbol == 'I' || symbol == 'X' || symbol == 'C'
}

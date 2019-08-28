package properties

import (
	"fmt"
	"testing"
	"testing/quick"
)

var cases = []struct {
	arabic uint16
	roman  string
}{
	{1, "I"},
	{2, "II"},
	{3, "III"},
	{4, "IV"},
	{5, "V"},
	{6, "VI"},
	{7, "VII"},
	{8, "VIII"},
	{9, "IX"},
	{10, "X"},
	{11, "XI"},
	{14, "XIV"},
	{18, "XVIII"},
	{20, "XX"},
	{39, "XXXIX"},
	{40, "XL"},
	{47, "XLVII"},
	{49, "XLIX"},
	{50, "L"},
	{100, "C"},
	{90, "XC"},
	{400, "CD"},
	{500, "D"},
	{900, "CM"},
	{1000, "M"},
	{1984, "MCMLXXXIV"},
	{3999, "MMMCMXCIX"},
	{2014, "MMXIV"},
	{1006, "MVI"},
	{798, "DCCXCVIII"},
}

func TestConvertingToRomanNumerals(t *testing.T) {

	for _, tt := range cases {

		t.Run(fmt.Sprintf("%d get converted to %q", tt.arabic, tt.roman), func(t *testing.T) {
			if got := ConvertToRoman(tt.arabic); got != tt.roman {
				t.Errorf("got %q want %q", got, tt.roman)
			}
		})
	}
}

func TestConvertingToArabic(t *testing.T) {

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%q get converted to %d", tt.roman, tt.arabic), func(t *testing.T) {
			if got := ConvertToArabic(tt.roman); got != tt.arabic {
				t.Errorf("got %d want %d", got, tt.arabic)
			}
		})
	}
}

func TestPropertiesOfConversion(t *testing.T) {
	assertion := func(arabic uint16) bool {
		if arabic > 3999 {
			return true
		}
		roman := ConvertToRoman(arabic)
		from_roman := ConvertToArabic(roman)

		return from_roman == arabic
	}
	config := &quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(assertion, config); err != nil {
		t.Error("failed checks", err)
	}
}

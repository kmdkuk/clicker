package formatter

import (
	"fmt"
	"math"
)

// 3桁ごとの単位定義
var units = []string{"", "K", "M", "B", "T"}

// FormatLargeNumber は大きな数値を3桁ごとの指数表記に変換します
// 例: 1000 -> 1.00K, 1500 -> 1.50K, 1000000 -> 1.00M
func FormatLargeNumber(value float64) string {
	// 0や負の値は特別扱い
	if value == 0 {
		return "0.00"
	}
	if value < 0 {
		return "-" + FormatLargeNumber(-value)
	}

	// 値が1000未満の場合、範囲に応じて異なる小数点以下の桁数を使用
	// For values under 10, use two decimal places to provide finer precision for small numbers.
	if value < 10 {
		return fmt.Sprintf("%.2f", floor(value, 2))
	}
	// For values between 10 and 100, use one decimal place to balance precision and readability.
	if value < 100 {
		return fmt.Sprintf("%.1f", floor(value, 1))
	}
	// For values between 100 and 1000, use no decimal places to simplify the representation.
	if value < 1000 {
		return fmt.Sprintf("%.0f", floor(value, 0))
	}

	// 3桁ごとの指数を計算
	exp := int(math.Floor(math.Log10(value) / 3))
	if exp >= len(units) {
		// 単位定義を超えた大きさの場合は標準的な科学的記数法を使用
		return fmt.Sprintf("%.2e", value)
	}

	// 対応する単位で値をスケーリング
	scaledValue := value / math.Pow(1000, float64(exp))
	formattedValue := fmt.Sprintf("%.0f", floor(scaledValue, 0))
	if scaledValue < 10 {
		formattedValue = fmt.Sprintf("%.2f", floor(scaledValue, 2))
	} else if scaledValue < 100 {
		formattedValue = fmt.Sprintf("%.1f", floor(scaledValue, 1))
	}

	// 表示用の文字列を整形

	return formattedValue + units[exp]
}

// FormatCurrency は通貨値を整形します（通貨記号付き）
func FormatCurrency(value float64, symbol string) string {
	return symbol + " " + FormatLargeNumber(value)
}

func floor(value float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Floor(value*p) / p
}

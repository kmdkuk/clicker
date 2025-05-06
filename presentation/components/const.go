package components

import "image/color"

const (
	// アイテム表示関連
	ItemHeight        = 40 // リスト項目の高さ
	ItemTextPadding   = 5  // テキストと背景の余白
	ItemVerticalShift = 1  // 背景矩形の垂直位置調整
	TextSize          = 24 // テキストサイズ

	// スクロールバー関連
	ScrollbarWidth      = 8  // スクロールバーの幅
	ScrollbarMargin     = 5  // スクロールバーの余白
	MinimumHandleHeight = 10 // スクロールバーハンドルの最小高さ
	ViewportSize        = 8  // ビューポートのサイズ（表示可能なアイテム数）
)

// カラー定義
var (
	// 選択状態のカラー
	SelectedBgColor   = color.RGBA{R: 50, G: 100, B: 200, A: 160}
	SelectedTextColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// 非選択状態のカラー
	NormalBgColor   = color.RGBA{R: 40, G: 40, B: 40, A: 120}
	NormalTextColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}

	// スクロールバーのカラー
	ScrollbarTrackColor  = color.RGBA{R: 80, G: 80, B: 80, A: 180}
	ScrollbarHandleColor = color.RGBA{R: 180, G: 180, B: 180, A: 255}
)

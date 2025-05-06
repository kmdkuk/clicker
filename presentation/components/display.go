package components

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/assets/fonts"
	"github.com/kmdkuk/clicker/presentation/formatter"
)

// DisplayComponent shows basic game information
type Display struct {
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) DrawMoney(screen *ebiten.Image, playerDTO *dto.Player) {
	moneyText := fmt.Sprintf("Money: %s (Total Generate Rate: %s/s)",
		formatter.FormatCurrency(playerDTO.GetMoney(), "$"),
		formatter.FormatCurrency(playerDTO.GetTotalGenerateRate(), "$"),
	)
	x := 10
	y := 10
	// 画面の幅を取得してアイテムの背景幅を決定
	screenWidth := screen.Bounds().Dx()

	// アイテムの幅を計算（スクロールバー分を引く）
	itemWidth := screenWidth - x - ScrollbarWidth - ScrollbarMargin*2

	bgColor := NormalBgColor

	// 背景矩形を描画
	rectY := float32(y)
	rectWidth := float32(itemWidth)
	rectHeight := float32(ItemHeight - ItemVerticalShift)
	vector.DrawFilledRect(screen, float32(x), rectY, rectWidth, rectHeight, bgColor, false)

	// テキストの色を設定（選択中かどうかで分ける）
	textColor := NormalTextColor

	// フォントフェイスを作成
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.BebasNeueRegular_ttf))
	if err != nil {
		fmt.Printf("Error loading font: %v", err)
		return
	}

	face := &text.GoTextFace{
		Source: s,
		Size:   float64(TextSize),
	}

	rectCenterY := float32(y) + rectHeight/2
	textX := float64(x + ItemTextPadding)
	textY := float64(rectCenterY)

	// テキスト描画
	txtOp := &text.DrawOptions{}
	txtOp.PrimaryAlign = text.AlignStart
	txtOp.SecondaryAlign = text.AlignCenter
	txtOp.GeoM.Translate(textX, textY)
	txtOp.ColorScale.ScaleWithColor(textColor)

	text.Draw(screen, moneyText, face, txtOp)
}

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
	x int
	y int
}

func NewDisplay(x, y int) *Display {
	return &Display{
		x: x,
		y: y,
	}
}

func (d *Display) calcItemWidthHeight(screenWidth int) (float32, float32) {
	itemWidth := screenWidth - d.x - ScrollbarWidth - ScrollbarMargin*2
	itemHeight := ItemHeight - ItemVerticalShift
	return float32(itemWidth), float32(itemHeight)
}

func (d *Display) DrawMoney(screen *ebiten.Image, playerDTO *dto.Player) {
	moneyText := fmt.Sprintf("Money: %s (Total Generate Rate: %s/s)",
		formatter.FormatCurrency(playerDTO.GetMoney(), "$"),
		formatter.FormatCurrency(playerDTO.GetTotalGenerateRate(), "$"),
	)

	bgColor := NormalBgColor

	// 背景矩形を描画
	rectWidth, rectHeight := d.calcItemWidthHeight(screen.Bounds().Dx())
	vector.DrawFilledRect(screen, float32(d.x), float32(d.y), rectWidth, rectHeight, bgColor, false)

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

	rectCenterY := float32(d.y) + rectHeight/2
	textX := float64(d.x) + ItemTextPadding
	textY := float64(rectCenterY)

	// テキスト描画
	txtOp := &text.DrawOptions{}
	txtOp.PrimaryAlign = text.AlignStart
	txtOp.SecondaryAlign = text.AlignCenter
	txtOp.GeoM.Translate(textX, textY)
	txtOp.ColorScale.ScaleWithColor(textColor)

	text.Draw(screen, moneyText, face, txtOp)
}

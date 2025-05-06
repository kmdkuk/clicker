package components

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/kmdkuk/clicker/assets/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tab struct {
	titles     []string
	x          int
	y          int
	activePage int
}

func NewTab(items []string, defaultPage, x, y int) *Tab {
	return &Tab{
		titles:     items,
		x:          x,
		y:          y,
		activePage: defaultPage,
	}
}

func (t *Tab) SetActivePage(index int) {
	if index >= 0 && index < len(t.titles) {
		t.activePage = index
	}
}

func (t *Tab) GetActivePage() int {
	return t.activePage
}

// Draw はタブを描画します
func (t *Tab) Draw(screen *ebiten.Image) {
	if len(t.titles) == 0 {
		return
	}

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

	// 現在のX位置（水平方向に並べるため）
	currentX := t.x

	// 各タブを描画
	for i, item := range t.titles {
		if i != 0 {
			// 最初のタブ以外は、前のタブの右側に配置
			currentX += ItemVerticalShift
		}
		label := item
		isSelected := i == t.activePage
		tabWidth := ((screen.Bounds().Dx() - (10 + ScrollbarWidth + ScrollbarMargin*2)) / len(t.titles)) - ItemVerticalShift*(len(t.titles)-1)

		// 背景色を設定（選択中かホバー中かで分ける）
		var bgColor color.RGBA
		if isSelected {
			bgColor = SelectedBgColor
		} else {
			bgColor = NormalBgColor
		}

		// タブの背景を描画（上部に丸みをつける）
		t.drawTabBackground(screen, currentX, t.y+ItemVerticalShift, tabWidth, ItemHeight-(ItemVerticalShift*2), bgColor)

		// テキストの色を設定
		var textColor color.RGBA
		if isSelected {
			textColor = SelectedTextColor
		} else {
			textColor = NormalTextColor
		}

		// 背景矩形の中央を計算
		rectCenterX := float64(currentX + tabWidth/2)
		rectCenterY := float64(t.y + ItemHeight/2)

		// テキスト描画
		txtOp := &text.DrawOptions{}
		txtOp.PrimaryAlign = text.AlignCenter
		txtOp.SecondaryAlign = text.AlignCenter
		txtOp.GeoM.Translate(rectCenterX, rectCenterY)
		txtOp.ColorScale.ScaleWithColor(textColor)

		text.Draw(screen, label, face, txtOp)

		// 次のタブの開始位置
		currentX += tabWidth + ItemVerticalShift
	}
}

// drawTabBackground はタブの背景を描画します（上側の角のみ丸くする）
func (t *Tab) drawTabBackground(screen *ebiten.Image, x, y, width, height int, bgColor color.RGBA) {
	// 標準の矩形描画
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(width),
		float32(height),
		bgColor,
		false,
	)
}

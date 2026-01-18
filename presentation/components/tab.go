package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tab struct {
	source     *text.GoTextFaceSource
	titles     []string
	x          int
	y          int
	activePage int
}

func NewTab(source *text.GoTextFaceSource, items []string, defaultPage, x, y int) *Tab {
	return &Tab{
		source:     source,
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

	face := &text.GoTextFace{
		Source: t.source,
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
		tabWidth, tabHeight := t.getTabSize(screen.Bounds().Dx())

		// 背景色を設定（選択中かホバー中かで分ける）
		var bgColor color.RGBA
		if isSelected {
			bgColor = SelectedBgColor
		} else {
			bgColor = NormalBgColor
		}

		// タブの背景を描画（上部に丸みをつける）
		t.drawTabBackground(screen, currentX, t.y+ItemVerticalShift, tabWidth, tabHeight, bgColor)

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
	vector.FillRect(
		screen,
		float32(x),
		float32(y),
		float32(width),
		float32(height),
		bgColor,
		false,
	)
}

func (t *Tab) getTabSize(screenWidth int) (int, int) {
	width := ((screenWidth - (10 + ScrollbarWidth + ScrollbarMargin*2)) / len(t.titles)) - ItemVerticalShift*(len(t.titles)-1)
	height := ItemHeight - (ItemVerticalShift * 2)
	return width, height
}

func (t *Tab) GetHoverPage(screenWidth, mouseX, mouseY int) int {
	// タブの幅を計算
	tabWidth, tabHeight := t.getTabSize(screenWidth)

	// 各タブをチェック
	xOffset := 10
	for i := 0; i < len(t.titles); i++ {
		if mouseX >= xOffset && // left side
			mouseX < xOffset+tabWidth && // right side
			mouseY >= t.y+ItemVerticalShift && // top side
			mouseY < t.y+tabHeight-ItemVerticalShift { // bottom side
			return i
		}
		xOffset += (tabWidth + 2*ItemVerticalShift)
	}
	return -1
}

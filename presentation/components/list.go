package components

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/assets/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func ConvertBuildingToListItems(buildings []dto.Building) []ListItem {
	items := make([]ListItem, len(buildings))
	for i := range buildings {
		items[i] = &buildings[i]
	}
	return items
}

func ConvertUpgradeToListItems(upgrades []dto.Upgrade) []ListItem {
	items := make([]ListItem, len(upgrades))
	for i := range upgrades {
		items[i] = &upgrades[i]
	}
	return items
}

type ListItem interface {
	String() string
}

type List struct {
	source       *text.GoTextFaceSource
	Items        []ListItem
	Visible      bool
	x            int
	y            int
	scrollPos    int // Current scroll position
	viewportSize int // Number of items visible at once
}

func NewList(source *text.GoTextFaceSource, defaultVisible bool, x, y int) *List {
	return NewListWithViewport(source, defaultVisible, x, y, ViewportSize)
}

// NewListWithViewport creates a new list with a specified viewport size and font
func NewListWithViewport(source *text.GoTextFaceSource, defaultVisible bool, x, y, viewportSize int) *List {
	return &List{
		source:       source,
		Items:        []ListItem{},
		Visible:      defaultVisible,
		x:            x,
		y:            y,
		scrollPos:    0,
		viewportSize: viewportSize,
	}
}

func (l *List) Draw(screen *ebiten.Image, cursor int) {
	if (!l.Visible) || (len(l.Items) == 0) {
		return
	}

	// Adjust scroll position if cursor moves outside the viewport
	if cursor < l.scrollPos {
		l.scrollPos = cursor
	} else if cursor >= l.scrollPos+l.viewportSize {
		l.scrollPos = cursor - l.viewportSize + 1
	}

	// Ensure scroll position stays within valid range
	if l.scrollPos < 0 {
		l.scrollPos = 0
	}
	maxScroll := len(l.Items) - l.viewportSize
	if maxScroll < 0 {
		maxScroll = 0
	}
	if l.scrollPos > maxScroll {
		l.scrollPos = maxScroll
	}

	// Calculate ending index for display range
	endIdx := l.scrollPos + l.viewportSize
	if endIdx > len(l.Items) {
		endIdx = len(l.Items)
	}

	// Draw only items within the current viewport
	for i := l.scrollPos; i < endIdx; i++ {
		item := l.Items[i]
		displayY := l.y + (i-l.scrollPos)*ItemHeight + ItemVerticalShift
		l.drawItem(screen, item, l.x, displayY, i == cursor)
	}

	// Draw scrollbar instead of indicators
	l.drawScrollBar(screen, endIdx)
}

func (l *List) calcItemWidthHeight(screenWidth int, x, y int) (float32, float32) {
	itemWidth := screenWidth - x - ScrollbarWidth - ScrollbarMargin*2
	itemHeight := ItemHeight - ItemVerticalShift*2

	return float32(itemWidth), float32(itemHeight)

}

func (l *List) drawItem(screen *ebiten.Image, item ListItem, x, y int, isSelected bool) {
	var bgColor color.RGBA
	if isSelected {
		bgColor = SelectedBgColor
	} else {
		bgColor = NormalBgColor
	}

	// 背景矩形を描画
	rectWidth, rectHeight := l.calcItemWidthHeight(screen.Bounds().Dx(), x, y)
	vector.DrawFilledRect(screen, float32(x), float32(y), rectWidth, rectHeight, bgColor, false)

	// テキストの色を設定（選択中かどうかで分ける）
	var textColor color.RGBA
	if isSelected {
		textColor = SelectedTextColor
	} else {
		textColor = NormalTextColor
	}

	// テキストを描画
	textStr := item.String()
	if isSelected {
		textStr = "> " + textStr // 選択中の項目には矢印をつける
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

	rectCenterY := float32(y) + rectHeight/2
	textX := float64(x + ItemTextPadding)
	textY := float64(rectCenterY)

	// テキスト描画
	txtOp := &text.DrawOptions{}
	txtOp.PrimaryAlign = text.AlignStart
	txtOp.SecondaryAlign = text.AlignCenter
	txtOp.GeoM.Translate(textX, textY)
	txtOp.ColorScale.ScaleWithColor(textColor)

	text.Draw(screen, textStr, face, txtOp)
}

// drawScrollBar draws a scrollbar on the right edge of the screen
func (l *List) drawScrollBar(screen *ebiten.Image, endIdx int) {
	// Get screen dimensions
	screenWidth := screen.Bounds().Dx()

	// Skip drawing scrollbar if all items fit in viewport
	if len(l.Items) <= l.viewportSize {
		return
	}

	// Calculate scrollbar position
	scrollbarX := float64(screenWidth - ScrollbarWidth - ScrollbarMargin)
	scrollbarY := float64(l.y)

	// Calculate scrollbar height based on visible range
	visibleCount := endIdx - l.scrollPos
	listHeight := float64(visibleCount * ItemHeight)

	// Draw scrollbar background (track)
	vector.DrawFilledRect(screen, float32(scrollbarX), float32(scrollbarY),
		float32(ScrollbarWidth), float32(listHeight), ScrollbarTrackColor, false)

	// Calculate handle size and position
	totalItems := len(l.Items)
	handleRatio := float64(visibleCount) / float64(totalItems)
	handleHeight := listHeight * handleRatio
	if handleHeight < MinimumHandleHeight {
		handleHeight = MinimumHandleHeight // 最小ハンドル高さを確保
	}

	// Position the handle based on scroll position
	scrollRatio := 0.0
	if totalItems > visibleCount { // ゼロ除算を避ける
		scrollRatio = float64(l.scrollPos) / float64(totalItems-visibleCount)
	}
	handleY := scrollbarY + scrollRatio*(listHeight-handleHeight)

	// Draw scrollbar handle
	vector.DrawFilledRect(screen, float32(scrollbarX), float32(handleY),
		float32(ScrollbarWidth), float32(handleHeight), ScrollbarHandleColor, false)
}

// Scroll manually adjusts the scroll position by the specified amount
func (l *List) Scroll(amount int) {
	l.scrollPos += amount

	// Prevent scrolling out of bounds
	if l.scrollPos < 0 {
		l.scrollPos = 0
	}

	maxScroll := len(l.Items) - l.viewportSize
	if maxScroll < 0 {
		maxScroll = 0
	}
	if l.scrollPos > maxScroll {
		l.scrollPos = maxScroll
	}
}

// GetVisibleRange returns the start and end indices of currently visible items (for testing)
func (l *List) GetVisibleRange() (start, end int) {
	end = l.scrollPos + l.viewportSize
	if end > len(l.Items) {
		end = len(l.Items)
	}
	return l.scrollPos, end
}

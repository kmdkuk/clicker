package components

import (
	"image/color"

	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 吹き出し型ポップアップの設定定数
const (
	PopupMargin       = 20  // ポップアップの余白
	PopupPadding      = 15  // 内部パディング
	PopupCornerRadius = 10  // 角の丸み
	PopupMaxWidth     = 400 // 最大幅
	PopupTextSize     = 18  // テキストサイズ
	PopupLineHeight   = 28  // 行の高さ
	PopupTailSize     = 15  // 吹き出しの尻尾のサイズ
)

// ポップアップの色設定
var (
	PopupBgColor     = color.RGBA{R: 50, G: 50, B: 70, A: 240}    // 背景色
	PopupBorderColor = color.RGBA{R: 80, G: 80, B: 120, A: 255}   // 枠線色
	PopupTextColor   = color.RGBA{R: 240, G: 240, B: 240, A: 255} // テキスト色
)

type Popup struct {
	source  *text.GoTextFaceSource
	Message string // 表示メッセージ
	Active  bool   // アクティブ状態
}

func NewPopup(source *text.GoTextFaceSource) *Popup {
	return &Popup{
		source:  source,
		Message: "",
		Active:  false,
	}
}

func (p *Popup) Draw(screen *ebiten.Image) {
	if !p.IsActive() {
		return
	}

	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	face := &text.GoTextFace{
		Source: p.source,
		Size:   float64(PopupTextSize),
	}

	// ポップアップのサイズ設定
	popupWidth := float32(screenWidth)/2 + float32(PopupPadding*2)
	popupHeight := float32(screenHeight/2) + float32(PopupPadding*2)

	// 画面中央の座標を計算
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)

	// ポップアップの左上座標
	popupX := centerX - popupWidth/2
	popupY := centerY - popupHeight/2

	p.drawBackground(screen, popupX, popupY, popupWidth, popupHeight)

	txtOp := &text.DrawOptions{}
	txtOp.PrimaryAlign = text.AlignCenter
	txtOp.SecondaryAlign = text.AlignCenter
	txtOp.GeoM.Translate(float64(centerX), float64(centerY))
	txtOp.ColorScale.ScaleWithColor(PopupTextColor)

	text.Draw(screen, p.Message, face, txtOp)

	// クローズヒントの表示
	hintText := "Press Enter or Space to close"
	hintX := popupX + popupWidth - PopupPadding
	hintY := popupY + popupHeight - PopupPadding

	txtOp = &text.DrawOptions{}
	txtOp.PrimaryAlign = text.AlignEnd
	txtOp.SecondaryAlign = text.AlignEnd
	txtOp.GeoM.Translate(float64(hintX), float64(hintY))
	txtOp.ColorScale.ScaleWithColor(color.RGBA{R: 180, G: 180, B: 180, A: 180})
	text.Draw(screen, hintText, face, txtOp)
}

// 吹き出し背景の描画
func (p *Popup) drawBackground(screen *ebiten.Image, x, y, width, height float32) {
	// メイン背景の描画
	vector.DrawFilledRect(screen,
		float32(x), float32(y),
		float32(width), float32(height),
		PopupBorderColor, false)
	vector.DrawFilledRect(screen,
		float32(x+1), float32(y+1),
		float32(width-2), float32(height-2),
		PopupBgColor, false)

}

func (p *Popup) HandleInput(keyType input.KeyType) {
	if p.IsActive() && keyType == input.KeyTypeDecision {
		p.Close()
	}
}

func (p *Popup) IsActive() bool {
	return p.Active
}

func (p *Popup) GetMessage() string {
	return p.Message
}

func (p *Popup) Show(message string) {
	p.Message = message
	p.Active = true
}

func (p *Popup) Close() {
	p.Active = false
}

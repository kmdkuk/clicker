package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyType int

const (
	KeyTypeUp       KeyType = iota // 上
	KeyTypeDown                    // 下
	KeyTypeLeft                    // 左
	KeyTypeRight                   // 右
	KeyTypeDecision                // 決定
	KeyTypeNone                    // その他のキーまたは未入力
)

type InputHandler struct {
	pressedKey ebiten.Key // 押されたキーを記録
}

// Update メソッドで押されたキーを記録
func (ih *InputHandler) Update() {
	ih.pressedKey = ebiten.Key(0) // 初期化

	// 押されたキーを記録
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		ih.pressedKey = key
		break // 最初に押されたキーのみ記録
	}
}

// GetPressedKey メソッドで押されたキーを分類して取得
func (ih *InputHandler) GetPressedKey() KeyType {
	switch ih.pressedKey {
	case ebiten.KeyArrowUp, ebiten.KeyW, ebiten.KeyK:
		return KeyTypeUp // 方向キー: 上
	case ebiten.KeyArrowDown, ebiten.KeyS, ebiten.KeyJ:
		return KeyTypeDown // 方向キー: 下
	case ebiten.KeyArrowLeft, ebiten.KeyA, ebiten.KeyH:
		return KeyTypeLeft // 方向キー: 左
	case ebiten.KeyArrowRight, ebiten.KeyD, ebiten.KeyL:
		return KeyTypeRight // 方向キー: 右
	case ebiten.KeyEnter, ebiten.KeySpace:
		return KeyTypeDecision // 決定キー
	default:
		return KeyTypeNone // その他のキーまたは未入力
	}
}

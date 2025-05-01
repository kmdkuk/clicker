package game

type Decider interface {
	Decide(page, cursor int) (bool, string)
}

type DefaultDecider struct {
	gameState GameState
}

func NewDefaultDecider(gameState GameState) Decider {
	return &DefaultDecider{
		gameState: gameState,
	}
}

func (d *DefaultDecider) Decide(page, cursor int) (bool, string) {
	// マニュアルワークの選択
	if cursor == 0 {
		d.gameState.ManualWork()
		return true, ""
	}

	// 建物またはアップグレードの処理
	adjustedCursor := cursor - 1

	switch page {
	case 0: // 建物ページ
		return d.gameState.PurchaseBuildingAction(adjustedCursor)

	case 1: // アップグレードページ
		return d.gameState.PurchaseUpgradeAction(adjustedCursor)

	default:
		return false, "Invalid page selection"
	}
}

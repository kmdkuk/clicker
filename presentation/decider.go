package presentation

import "fmt"

type Decider interface {
	Decide(page, cursor int) (bool, string)
}

type DefaultDecider struct {
	ManualWorkUseCase ManualWorkUseCase
	BuildingUseCase   BuildingUseCase
	UpgradeUseCase    UpgradeUseCase
}

func NewDecider(manualWorkUseCase ManualWorkUseCase, buildingUseCase BuildingUseCase, upgradeUseCase UpgradeUseCase) Decider {
	return &DefaultDecider{
		ManualWorkUseCase: manualWorkUseCase,
		BuildingUseCase:   buildingUseCase,
		UpgradeUseCase:    upgradeUseCase,
	}
}

func (d *DefaultDecider) Decide(page, cursor int) (bool, string) {
	// マニュアルワークの選択
	if cursor == 0 {
		d.ManualWorkUseCase.ManualWorkAction()
		return true, ""
	}

	// 建物またはアップグレードの処理
	adjustedCursor := cursor - 1

	fmt.Printf("Page: %d, Cursor: %d\n", page, adjustedCursor)
	switch page {
	case 0: // 建物ページ
		return d.BuildingUseCase.PurchaseBuildingAction(adjustedCursor)

	case 1: // アップグレードページ
		return d.UpgradeUseCase.PurchaseUpgradeAction(adjustedCursor)

	default:
		return false, "Invalid page selection"
	}
}

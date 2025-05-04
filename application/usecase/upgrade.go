package usecase

import (
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/infrastructure/state"
)

func NewUpgradeUseCase(gameState state.GameState) *UpgradeUseCase {
	return &UpgradeUseCase{
		gameState: gameState,
	}
}

type UpgradeUseCase struct {
	gameState state.GameState
}

func (u *UpgradeUseCase) GetUpgrades() []dto.Upgrade {
	upgrades := make([]dto.Upgrade, len(u.gameState.GetUpgrades()))
	for i, upgrade := range u.gameState.GetUpgrades() {
		upgrades[i] = dto.Upgrade{
			Name:        upgrade.Name,
			IsPurchased: upgrade.IsPurchased,
			IsReleased:  upgrade.IsReleased(u.gameState),
			Cost:        upgrade.Cost,
		}
	}
	return upgrades
}

func (u *UpgradeUseCase) PurchaseUpgradeAction(upgradeIndex int) (bool, string) {
	if upgradeIndex < 0 || upgradeIndex >= len(u.gameState.GetUpgrades()) {
		return false, "Invalid upgrade selection!"
	}

	upgrade := &u.gameState.GetUpgrades()[upgradeIndex]

	if upgrade.IsPurchased {
		return false, "Upgrade already purchased!"
	}

	if !upgrade.IsReleased(u.gameState) {
		return false, "Upgrade not available yet!"
	}

	if u.gameState.GetMoney() < upgrade.Cost {
		return false, "Not enough money for upgrade!"
	}

	if err := u.gameState.SetUpgradesIsPurchased(upgradeIndex, true); err != nil {
		return false, "Failed to purchase upgrade!"
	}
	u.gameState.UpdateMoney(-upgrade.Cost)

	return true, "Upgrade purchased successfully!"
}

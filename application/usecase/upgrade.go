package usecase

import (
	"fmt"
	"sort"

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
			ID:          upgrade.ID,
			Name:        upgrade.Name,
			IsPurchased: upgrade.IsPurchased,
			IsReleased:  upgrade.IsReleased(u.gameState),
			Cost:        upgrade.Cost,
		}
	}
	return upgrades
}

func (u *UpgradeUseCase) GetUpgradesIsReleasedCostSorted() []dto.Upgrade {
	upgrades := u.GetUpgrades()
	upgradesIsRelease := make([]dto.Upgrade, 0)
	for _, upgrade := range upgrades {
		if upgrade.IsReleased {
			upgradesIsRelease = append(upgradesIsRelease, upgrade)
		}
	}

	sort.SliceStable(upgradesIsRelease, func(i, j int) bool {
		return upgradesIsRelease[i].Cost < upgradesIsRelease[j].Cost
	})

	return upgradesIsRelease
}

func (u *UpgradeUseCase) findUpgradeWithID(id string) (int, error) {
	upgrades := u.gameState.GetUpgrades()
	for i, upgrade := range upgrades {
		if upgrade.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("upgrade with ID %s not found", id)
}

func (u *UpgradeUseCase) PurchaseUpgradeAction(cursor int) (bool, string) {
	upgrades := u.GetUpgradesIsReleasedCostSorted()
	if cursor < 0 || cursor >= len(upgrades) {
		return false, "Invalid upgrade selection!"
	}

	upgrade := &upgrades[cursor]

	if upgrade.IsPurchased {
		return false, "Upgrade already purchased!"
	}

	if !upgrade.IsReleased {
		return false, "Upgrade not available yet!"
	}

	if u.gameState.GetMoney() < upgrade.Cost {
		return false, "Not enough money for upgrade!"
	}

	index, err := u.findUpgradeWithID(upgrade.ID)
	if err != nil {
		return false, "Failed to find upgrade"
	}
	if err := u.gameState.SetUpgradesIsPurchased(index, true); err != nil {
		return false, "Failed to purchase upgrade!"
	}
	u.gameState.UpdateMoney(-upgrade.Cost)

	return true, "Upgrade purchased successfully!"
}

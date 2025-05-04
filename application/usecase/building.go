package usecase

import (
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/infrastructure/state"
)

func NewBuildingUseCase(gameState state.GameState) *BuildingUseCase {
	return &BuildingUseCase{
		gameState: gameState,
	}
}

type BuildingUseCase struct {
	gameState state.GameState
}

func (b *BuildingUseCase) GetBuildings() []dto.Building {
	buildings := make([]dto.Building, len(b.gameState.GetBuildings()))
	for i, building := range b.gameState.GetBuildings() {
		genRate := building.BaseGenerateRate
		if building.IsUnlocked() {
			genRate = building.TotalGenerateRate(b.gameState.GetUpgrades())
		}
		buildings[i] = dto.Building{
			Name:              building.Name,
			IsUnlocked:        building.IsUnlocked(),
			Count:             building.Count,
			Cost:              building.Cost(),
			TotalGenerateRate: genRate,
		}
	}
	return buildings
}

func (b *BuildingUseCase) GetBuildingsIsUnlockedWithMaskedNextLock() []dto.Building {
	buildings := b.GetBuildings()
	buildingsInMaskedUnlock := make([]dto.Building, 0)
	unlockIndex := -1
	if len(buildings) == 0 {
		return buildingsInMaskedUnlock
	}
	for i := range buildings {
		building := buildings[i]
		if building.IsUnlocked {
			unlockIndex = i
		} else {
			building.Name = "???"
		}
		buildingsInMaskedUnlock = append(buildingsInMaskedUnlock, building)
	}
	if unlockIndex+2 >= len(buildings) {
		return buildingsInMaskedUnlock
	}
	return buildingsInMaskedUnlock[:unlockIndex+2]
}

func (b *BuildingUseCase) PurchaseBuildingAction(buildingIndex int) (bool, string) {
	buildings := b.gameState.GetBuildings()
	if buildingIndex < 0 || buildingIndex >= len(buildings) {
		return false, "Invalid building selection!"
	}

	building := &buildings[buildingIndex]
	cost := building.Cost()

	if b.gameState.GetMoney() < cost {
		if building.IsUnlocked() {
			return false, "Not enough money to purchase!"
		}
		return false, "Not enough money to unlock!"
	}

	building.Count++
	if err := b.gameState.SetBuildingCount(buildingIndex, building.Count); err != nil {
		return false, "Failed to update building count!"
	}
	b.gameState.UpdateMoney(-cost)

	return true, "Building purchased successfully!"
}

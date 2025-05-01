package input

type GameStateWriter interface {
	ManualWork()
	PurchaseBuildingAction(cursor int) (bool, string)
	PurchaseUpgradeAction(cursor int) (bool, string)
}

package input

type GameStateWriter interface {
	ManualWorkAction()
	PurchaseBuildingAction(cursor int) (bool, string)
	PurchaseUpgradeAction(cursor int) (bool, string)
}

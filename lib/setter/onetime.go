package setter

type Setter interface {
	// Update前に呼び出される
	OnBeforeUpdate()
	// Update語に呼び出される
	OnAfterUpdate(index int)
}

type Unit interface {
	// このUnitを取得すべき対象かを表す
	Check() bool
	// このUnitを登録したときに呼び出される
	OnStarted()
	// このUnitの登録が解除されたときに呼び出される
	OnClosed()
}

type OnetimeUnit struct {
}

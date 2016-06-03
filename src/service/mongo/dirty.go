package mongo

type Dirty interface {
	CRUD(interface{}) bool
	Inspect() string
}

package service

type MongoCRUD int

const (
	// basic database operation : CRUD
	MongoActionUnknown MongoCRUD = 0
	MongoActionCreate  MongoCRUD = 1
	MongoActionRead    MongoCRUD = 2
	MongoActionUpdate  MongoCRUD = 3
	MongoActionDelete  MongoCRUD = 4
)

type MongoDirty struct {
	Action MongoCRUD
	dirty
}

type dirty interface {
	Create() bool
	Read() bool
	Update() bool
	Delete() bool
    Inspect() string
}

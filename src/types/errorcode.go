package types

const (
    OK              = int32(0)
    ERR_INVALID_REQ = int32(1)
    ERR_GROUP_NOT_FOUND = int32(2)
)

var ErrDesc = map[int32]string{
    OK:              "success",
    ERR_INVALID_REQ: "invalid req message type",
    ERR_GROUP_NOT_FOUND: "boradcast group not found",
}
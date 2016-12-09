package types

const (
    OK              = int32(0)
    ERR_INVALID_REQ = int32(1)
)

var ErrDesc = map[int32]string{
    OK:              "success",
    ERR_INVALID_REQ: "invalid req message type",
}
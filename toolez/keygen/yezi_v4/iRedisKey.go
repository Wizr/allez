package yezi_v4

type iRedisKey interface {
	iSerializable
	Get() error
	Set() error
	GetKey() string
}

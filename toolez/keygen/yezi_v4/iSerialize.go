package yezi_v4
type iSerializable interface {
	Serialize()
	Deserialize() error
}


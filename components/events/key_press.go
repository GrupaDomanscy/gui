package events

type KeyPressEvent struct {
	PressedChars []int32
	DownKeys     []int32
	PressedKeys  []int32
	ReleasedKeys []int32
}

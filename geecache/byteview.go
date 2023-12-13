package geecache

/*缓存值的抽象与封装*/
type ByteView struct {
	b []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{b}
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b) // 防止被外界对象修改
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

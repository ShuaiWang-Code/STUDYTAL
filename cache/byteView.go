package cache

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte // byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等
}

// Len returns the view's length
// 被缓存对象必须实现 Value 接口
func (v ByteView) Len() int {
	return len(v.b) // 返回其所占的内存大小
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

// b 是只读的，使用 ByteSlice() 调用该方法返回一个拷贝，防止缓存值被外部程序修改。
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

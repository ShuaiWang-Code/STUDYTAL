package cache

import (
	"reflect"
	"testing"
)

// 测试接口型函数；标准库中用得也不少，net/http 的 Handler 和 HandlerFunc 就是一个典型。
// 借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter
func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

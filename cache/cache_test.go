package cache

import (
	"fmt"
	"log"
	"testing"
)

// mock数据源
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	// 数据源数量
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// 在缓存为空的情况下，能够通过回调函数获取到源数据。
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		} // load from callback function
		// 在缓存已经存在的情况下，是否直接从缓存中获取
		// 使用 loadCounts 统计某个键调用回调函数的次数，如果次数大于1，则表示调用了多次回调函数，没有缓存
		// 从 mainCache 中查找缓存，如果存在则返回缓存值并打印hit
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	// 如果设定内存大小为13，则第一个key被淘汰，cache tom miss
	// if view, err := gee.Get("Tom"); err != nil || view.String() != "630" {
	// 	t.Fatal("failed to get value of Tom")
	// } // load from callback function
	// if _, err := gee.Get("Tom"); err != nil || loadCounts["Tom"] > 1 {
	// 	t.Fatalf("cache %s miss", "Tom")
	// } // cache hit

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

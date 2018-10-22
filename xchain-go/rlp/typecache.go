// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rlp

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	typeCacheMutex sync.RWMutex                  // 读写锁，用来在多线程的时候保护typeCache这个Map
	typeCache      = make(map[typekey]*typeinfo) // 核心数据结构，保存了类型->编解码器函数
)

/*
存储了编\解码器函数
*/
type typeinfo struct { // 存储了编\解码器函数
	decoder
	writer
}

// represents struct tags
type tags struct {
	// rlp:"nil" controls whether empty input results in a nil pointer.
	// rlp：“nil”控制空输入是否产生nil指针。
	nilOK bool
	// rlp:"tail" controls whether this field swallows additional list
	// elements. It can only be set for the last field, which must be
	// of slice type.
	// rlp：“tail”控制该字段是否吞下其他列表元素。 它只能设置为最后一个字段，该字段必须是切片类型。
	tail bool
	// rlp:"-" ignores fields.
	ignored bool
}

/*
typecache的typekey
*/
type typekey struct {
	reflect.Type
	// the key must include the struct tags because they
	// might generate a different decoder.
	//key必须包含tags结构，因为它们可能会生成不同的解码器
	tags
}

type decoder func(*Stream, reflect.Value) error

type writer func(reflect.Value, *encbuf) error

/*
cachedTypeInfo
首先加读锁，根据typekey读取typeCache中的信息，解锁；
读到信息，返回；
否则，调用cachedTypeInfo1方法，将新的typekey写入typeCache
*/
func cachedTypeInfo(typ reflect.Type, tags tags) (*typeinfo, error) {
	typeCacheMutex.RLock()                // 加读锁来保护
	info := typeCache[typekey{typ, tags}] // 读到信息
	typeCacheMutex.RUnlock()              //解锁
	if info != nil {
		return info, nil
	}
	// not in the cache, need to generate info for this type.
	//否则加写锁 调用cachedTypeInfo1函数创建并返回
	//这里需要注意的是在多线程环境下有可能多个线程同时调用到这个地方
	//所以当你进入cachedTypeInfo1方法的时候需要判断一下是否已经被别的线程先创建成功了。
	typeCacheMutex.Lock()
	defer typeCacheMutex.Unlock()
	return cachedTypeInfo1(typ, tags)
}

/*
cachedTypeInfo1
将新的typekey写入typeCache；
先判断是否有值了，（可能有其他并发进程已经写入。所以避免重复写入），有值则直接返回
没有值的话，新建一个typecach对象，typeCache[key] = new(typeinfo)
调用genTypeInfo方法写入这个类型。
如果调用genTypeInfo方法报错了，则将创建的typeCache[key]删掉。
否则将*info值赋值给*typeCache[key]
*/
func cachedTypeInfo1(typ reflect.Type, tags tags) (*typeinfo, error) {
	key := typekey{typ, tags}
	info := typeCache[key]
	//先判断一下是否已经被创建成功了，判断是否有值了
	if info != nil {
		// another goroutine got the write lock first
		return info, nil
	}
	// put a dummmy value into the cache before generating.
	// if the generator tries to lookup itself, it will get
	// the dummy value and won't call itself recursively.
	//在生成之前将虚拟值放入缓存中。
	//如果生成器尝试查找自身，它将获取虚拟值，并且不会递归调用自身。
	typeCache[key] = new(typeinfo)
	info, err := genTypeInfo(typ, tags)
	if err != nil {
		// remove the dummy value if the generator fails
		delete(typeCache, key)
		return nil, err
	}
	//todo:不懂为什么再赋一次指针
	//都是指针变量，所以把实际值
	*typeCache[key] = *info
	return typeCache[key], err
}

type field struct {
	index int
	info  *typeinfo
}

/*
structFields函数遍历所有的字段，然后针对每一个字段调用cachedTypeInfo1。
可以看到这是一个递归的调用过程。
上面的代码中有一个需要注意的是f.PkgPath == "" 这个判断针对的是所有导出的字段，
所谓的导出的字段就是说以大写字母开头命令的字段。
*/
func structFields(typ reflect.Type) (fields []field, err error) {
	for i := 0; i < typ.NumField(); i++ {
		if f := typ.Field(i); f.PkgPath == "" { // exported
			tags, err := parseStructTag(typ, i)
			if err != nil {
				return nil, err
			}
			if tags.ignored {
				continue
			}
			info, err := cachedTypeInfo1(f.Type, tags)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field{i, info})
		}
	}
	return fields, nil
}

func parseStructTag(typ reflect.Type, fi int) (tags, error) {
	f := typ.Field(fi)
	var ts tags
	for _, t := range strings.Split(f.Tag.Get("rlp"), ",") {
		switch t = strings.TrimSpace(t); t {
		case "":
		case "-":
			ts.ignored = true
		case "nil":
			ts.nilOK = true
		case "tail":
			ts.tail = true
			if fi != typ.NumField()-1 {
				return ts, fmt.Errorf(`rlp: invalid struct tag "tail" for %v.%s (must be on last field)`, typ, f.Name)
			}
			if f.Type.Kind() != reflect.Slice {
				return ts, fmt.Errorf(`rlp: invalid struct tag "tail" for %v.%s (field type is not slice)`, typ, f.Name)
			}
		default:
			return ts, fmt.Errorf("rlp: unknown struct tag %q on %v.%s", t, typ, f.Name)
		}
	}
	return ts, nil
}

//产生一个新的编\解码器
func genTypeInfo(typ reflect.Type, tags tags) (info *typeinfo, err error) {
	info = new(typeinfo)
	//解码器
	if info.decoder, err = makeDecoder(typ, tags); err != nil {
		return nil, err
	}
	//编码器
	if info.writer, err = makeWriter(typ, tags); err != nil {
		return nil, err
	}
	return info, nil
}

func isUint(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use slist file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"container/list"
	"sync"
)

//Item used for flowquantity
type Item struct {
	Timestamp int64   `json:"timestamp"`
	InSpeed   float64 `json:"inSpeed"`
	OutSpeed  float64 `json:"outSpeed"`
}

//SafeLinkedList used for flowQuantity
type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

//NewSafeLinkedList 新创建SafeLinkedList容器
func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

//PushFront push item
func (slist *SafeLinkedList) PushFront(v interface{}) *list.Element {
	slist.Lock()
	defer slist.Unlock()
	return slist.L.PushFront(v)
}

//Front get first item
func (slist *SafeLinkedList) Front() *list.Element {
	slist.RLock()
	defer slist.RUnlock()
	return slist.L.Front()
}

//PopBack get and remove last item
func (slist *SafeLinkedList) PopBack() *list.Element {
	slist.Lock()
	defer slist.Unlock()

	back := slist.L.Back()
	if back != nil {
		slist.L.Remove(back)
	}

	return back
}

//Back get last element
func (slist *SafeLinkedList) Back() *list.Element {
	slist.Lock()
	defer slist.Unlock()

	return slist.L.Back()
}

//Len get len of list
func (slist *SafeLinkedList) Len() int {
	slist.RLock()
	defer slist.RUnlock()
	return slist.L.Len()
}

//PopAllStale delete time less than timestamp items
func (slist *SafeLinkedList) PopAllStale(timestamp int64) {
	slist.Lock()
	defer slist.Unlock()

	p := slist.L.Back()
	for p != nil {
		p = p.Prev()
		if p.Next().Value.(*Item).Timestamp < timestamp {
			slist.L.Remove(p.Next())
		} else {
			break
		}
	}
}

//FetchAllMatch get time bigger than timestamp items
func (slist *SafeLinkedList) FetchAllMatch(timestamp int64) []*list.Element {
	slist.Lock()
	defer slist.Unlock()

	ret := make([]*list.Element, 0)

	p := slist.L.Front()
	for p != nil {
		if p.Value.(*Item).Timestamp > timestamp {
			ret = append(ret, p.Value.(*list.Element))
			p = p.Next()
		} else {
			break
		}
	}
	return ret
}

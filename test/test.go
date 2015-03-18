// test
package main

import (
	"fmt"
	"sync"
)

type Set struct {
	m map[int]bool
	sync.RWMutex
}

func New() *Set {
	return &Set{
		m: map[int]bool{},
	}
}

func (s *Set) Add(item int) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item int) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) List() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func main() {
	// 初始化
	s := New()

	s.Add(1)
	s.Add(1)
	s.Add(2)

	s.Clear()
	if s.IsEmpty() {
		fmt.Println("0 item")
	}

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(2)

	if s.Has(2) {
		fmt.Println("2 does exist")
	}

	s.Remove(2)
	s.Remove(3)
	fmt.Println("list of all items", s.List())
}

//func main() {
//	file, err := os.Open("d:\\log")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer file.Close()

//	r := bufio.NewReader(file)

//	for i := 0; i < 2; i++ {
//		_, err := r.ReadBytes('\n')
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}

//	var delim = []byte(`","`)
//	for {
//		line, err := r.ReadBytes('\n')
//		if err != nil {
//			if err != io.EOF {
//				fmt.Println(err)
//				return
//			}
//			if len(line) == 0 {
//				return
//			}
//		}

//		index := bytes.IndexByte(line, ' ')
//		if index == -1 {
//			fmt.Println("沒有找到前綴")
//			return
//		}
//		prefix := line[:index]
//		fmt.Println("前綴", string(prefix))

//		contentIndex1 := bytes.IndexByte(line, '{')
//		if contentIndex1 == -1 {
//			fmt.Println("沒有找到'{'")
//			return
//		}

//		contentIndex2 := bytes.IndexByte(line, '}')
//		if contentIndex2 == -1 {
//			fmt.Println("沒有找到'}'")
//			return
//		}
//		if contentIndex1+1 == contentIndex2 {
//			continue
//		}

//		bytesArr := bytes.Split(line[contentIndex1+2:contentIndex2-2], delim)
//		for _, bs := range bytesArr {
//			fmt.Println(string(bs))
//		}
//	}
//}

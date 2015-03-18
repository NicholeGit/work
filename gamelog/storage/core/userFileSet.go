package core

import "sync"

type UserFileSet struct {
	m map[string]bool
	sync.RWMutex
}

func NewUserFileSet() *UserFileSet {
	return &UserFileSet{
		m: map[string]bool{},
	}
}

func (s *UserFileSet) Add(item string) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *UserFileSet) Remove(item string) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *UserFileSet) Has(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *UserFileSet) Len() int {
	return len(s.List())
}

func (s *UserFileSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
}

func (s *UserFileSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *UserFileSet) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

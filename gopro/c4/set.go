package main

import "sync"


type Set struct {
	m map[int]bool
	sync.RWMutex
}

// 工厂模式构造一个set
func SetFactory() *Set{
	return &Set{
		m:map[int]bool{},
	}
}

func(s *Set)Add(items ...int ){
	s.Lock()
	defer s.Unlock()

	if len(items) ==0 {
		return
	}
	for _,v:= range items{
		s.m[v]= true
	}
}

func(s *Set)Remove(items ...int ){
	s.Lock()
	defer s.Unlock()

	if len(items) == 0 {
		return
	}

	for _,v:=range items{
		delete(s.m,v)
	}
}

func(s *Set)Has(items int )bool{
	s.RLock()
	defer s.RUnlock()

	_,ok:=s.m[items]
	return ok
}
func(s *Set)Len( )int{
	return len(s.m)
}
func(s *Set)Clear(){
	s.Lock()
	defer s.Unlock()
	s.m = map[int]bool{}
}
func(s *Set)IsEmpty()bool{
	if s.Len() == 0 {
		return true
	}
	return false
}
func(s *Set)List()[]int{
	s.RLock()
	defer s.RUnlock()
	list:=[]int{}

	if s.IsEmpty(){
		return list
	}

	for k,_:=range s.m{
		list= append(list,k)
	}
	return list
}

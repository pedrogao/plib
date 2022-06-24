package regex

import (
	"sync/atomic"

	"github.com/liyue201/gostl/ds/set"
)

var (
	idCount uint32 = 0
)

type StateType int

const (
	GENERAL StateType = iota + 1
	END

	EPSILON = "Epsilon"
	CHAR    = "char"
	CHARSET = "charSet"
)

type State struct {
	Id   int
	Type StateType
	next map[string]*set.Set
}

type StateOption func(state *State)

func NewState(options ...StateOption) *State {
	s := &State{
		Id:   int(atomic.AddUint32(&idCount, 1)),
		Type: GENERAL,
		next: map[string]*set.Set{},
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *State) AddNext(edge string, nfaState *State) {
	st := s.next[edge]
	if st == nil {
		st = set.New()
		s.next[edge] = st
	}
	s.next[edge].Insert(nfaState)
}

func (s *State) SetType(st StateType) {
	s.Type = st
}

func (s *State) IsEnd() bool {
	return s.Type == END
}

type NFAState struct {
	*State
}

func NewNFAState(options ...StateOption) *NFAState {
	state := NewState(options...)
	return &NFAState{state}
}

func (s *NFAState) AddNext(edge string, nfaState *NFAState) {
	st := s.next[edge]
	if st == nil {
		st = set.New()
		s.next[edge] = st
	}
	s.next[edge].Insert(nfaState)
}

type DFAState struct {
	*State
	nfaStates   *set.Set
	allStateIds string
}

func NewDFAState(allStateIds string, states *set.Set,
	options ...StateOption) *DFAState {
	state := NewState(options...)

	s := &DFAState{
		State:       state,
		nfaStates:   states,
		allStateIds: allStateIds,
	}

	for iter := states.Begin(); iter.IsValid(); iter.Next() {
		// 如果有任意节点是终止态,新建的DFA节点就是终止态
		if (iter.Value().(*NFAState)).Type == END {
			s.Type = END
		}
	}

	return s
}

func (s *DFAState) Equals(other *DFAState) bool {
	return s.allStateIds == other.allStateIds
}

func (s *DFAState) AddNext(edge string, dfaState *DFAState) {
	st := s.next[edge]
	if st == nil {
		st = set.New()
		s.next[edge] = st
	}
	s.next[edge].Insert(dfaState)
}

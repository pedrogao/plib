package regex

import (
	"strings"

	"github.com/liyue201/gostl/ds/set"
)

// NFAGraph nfa graph
type NFAGraph struct {
	start,
	end *NFAState
}

// NewNFAGraph create nfa graph
func NewNFAGraph(start, end *NFAState) *NFAGraph {
	return &NFAGraph{
		start: start,
		end:   end,
	}
}

// |
func (g *NFAGraph) addParallelGraph(graph *NFAGraph) {
	start := NewNFAState()
	end := NewNFAState()
	start.AddNext(EPSILON, g.start)
	start.AddNext(EPSILON, graph.start)
	g.end.AddNext(EPSILON, end)
	graph.end.AddNext(EPSILON, end)
	g.start = start
	g.end = end
}

//
func (g *NFAGraph) addSeriesGraph(graph *NFAGraph) {
	g.end.AddNext(EPSILON, graph.start)
	g.end = graph.end
}

// * 重复0-n次
func (g *NFAGraph) repeatStar() {
	g.repeatPlus()
	g.addSToE()
}

// ? 重复0次
func (g *NFAGraph) addSToE() {
	g.start.AddNext(EPSILON, g.end)
}

// + 重复1-n次
func (g *NFAGraph) repeatPlus() {
	start := NewNFAState()
	end := NewNFAState()
	start.AddNext(EPSILON, g.start)
	g.end.AddNext(EPSILON, end)
	g.end.AddNext(EPSILON, g.start)
	g.start = start
	g.end = end
}

// DFAGraph dfa graph
type DFAGraph struct {
	start    *DFAState
	stateMap map[string]*DFAState
}

// NewDFAGraph create dfa graph
func NewDFAGraph(start *DFAState) *DFAGraph {
	return &DFAGraph{
		start:    start,
		stateMap: map[string]*DFAState{},
	}
}

// Get state
func (g *DFAGraph) Get(states *set.Set) *DFAState {
	var builder strings.Builder
	for iter := states.Begin(); iter.IsValid(); iter.Next() {
		builder.WriteByte('#')
		builder.WriteRune(rune(iter.Value().(*NFAState).Id))
	}
	key := builder.String()
	state, ok := g.stateMap[key]
	if ok {
		return state
	}
	dfaState := NewDFAState(key, states)
	g.stateMap[key] = dfaState
	return dfaState
}

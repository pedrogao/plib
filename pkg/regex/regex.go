package regex

import (
	"github.com/liyue201/gostl/ds/queue"
	"github.com/liyue201/gostl/ds/set"
)

type Regex struct {
	graph    *NFAGraph
	dfaGraph *DFAGraph
}

func NewRegex(graph *NFAGraph) *Regex {
	return &Regex{graph: graph}
}

func (r *Regex) Match(text string) bool {
	start := r.graph.start
	return r.match(text, 0, start)
}

func (r *Regex) match(text string, pos int,
	curState *NFAState) bool {
	if pos == len(text) {
		if curState.IsEnd() {
			return true
		}
		stateSet := curState.next[EPSILON]
		if stateSet == nil {
			return false
		}
		for nextState := stateSet.Begin(); nextState.IsValid(); nextState.Next() {
			if r.match(text, pos, nextState.Value().(*NFAState)) {
				return true
			}
		}
	}

	for edge, stateSet := range curState.next {
		// 这个if和else的先后顺序决定了是贪婪匹配还是非贪婪匹配
		if EPSILON == edge {
			// 如果是DFA模式,不会有EPSILON边,所以不会进这
			for nextState := stateSet.Begin(); nextState.IsValid(); nextState.Next() {
				if r.match(text, pos, nextState.Value().(*NFAState)) {
					return true
				}
			}
		} else {
			matchStrategy := Manager.Get(edge)
			if !matchStrategy.Match(text[pos], edge) {
				continue
			}
			// 遍历匹配策略
			for nextState := stateSet.Begin(); nextState.IsValid(); nextState.Next() {
				if r.match(text, pos+1, nextState.Value().(*NFAState)) {
					return true
				}
			}
		}
	}
	return false
}

func (r *Regex) IsDFAMatch(text string) bool {
	return r.isDFAMatch(text, 0, r.dfaGraph.start)
}

func (r *Regex) isDFAMatch(text string, pos int, startState *DFAState) bool {
	curState := startState
	for pos < len(text) {
		canContinue := false
		for edge, stateSet := range curState.next {
			matchStrategy := Manager.Get(edge)
			if matchStrategy.Match(text[pos], edge) {
				next := stateSet.First().Value()
				if next == nil {
					break
				}
				curState = next.(*DFAState)
				pos++
				canContinue = true
				break
			}
		}
		if !canContinue {
			return false
		}
	}
	return curState.IsEnd()
}

func Compile(regex string) *Regex {
	if regex == "" {
		return nil
	}
	graph := regex2nfa(regex)
	graph.end.SetType(END)
	return NewRegex(graph)
}

func regex2nfa(regex string) *NFAGraph {
	reader := NewReader(regex)
	var nfaGraph *NFAGraph

	for reader.hasNext() {
		ch := reader.Next()
		edge := ""
		switch ch {
		// 子表达式特殊处理
		case '(':
			subRegex := reader.getSubRegex(reader)
			newNFAGraph := regex2nfa(subRegex)
			checkRepeat(reader, newNFAGraph)
			if nfaGraph == nil {
				nfaGraph = newNFAGraph
			} else {
				nfaGraph.addSeriesGraph(newNFAGraph)
			}
		// 或表达式特殊处理
		case '|':
			remainRegex := reader.getRemainRegex(reader)
			newNFAGraph := regex2nfa(remainRegex)
			if nfaGraph == nil {
				nfaGraph = newNFAGraph
			} else {
				nfaGraph.addParallelGraph(newNFAGraph)
			}
		case '[':
			edge = getCharSetMatch(reader)
			// 暂时未支持零宽断言
		case '^':
			break
		// 暂未支持
		case '$':
			break
		case '.':
			edge = "."
			// 处理特殊占位符
		case '\\':
			nextCh := reader.Next()
			switch nextCh {
			case 'd':
				edge = "\\d"
			case 'D':
				edge = "\\D"
			case 'w':
				edge = "\\w"
			case 'W':
				edge = "\\W"
			case 's':
				edge = "\\s"
			case 'S':
				edge = "\\S"
			// 转义后的字符匹配
			default:
				edge = string(nextCh)
			}
		default:
			// 处理普通字符
			edge = string(ch)
		}
		if edge != "" {
			start := NewNFAState()
			end := NewNFAState()
			start.AddNext(edge, end)
			newNFAGraph := NewNFAGraph(start, end)
			checkRepeat(reader, newNFAGraph)
			if nfaGraph == nil {
				nfaGraph = newNFAGraph
			} else {
				nfaGraph.addSeriesGraph(newNFAGraph)
			}
		}
	}

	return nfaGraph
}

func checkRepeat(reader *Reader, newNFAGraph *NFAGraph) {
	nextCh := reader.Peek()
	switch nextCh {
	case '*':
		newNFAGraph.repeatStar()
		reader.Next()
	case '+':
		newNFAGraph.repeatPlus()
		reader.Next()
	case '?':
		newNFAGraph.addSToE()
		reader.Next()
	case '{':
		// 暂未支持{}指定重复次数
	default:
		return
	}
}

func getCharSetMatch(reader *Reader) string {
	charset := ""
	for ch := reader.Next(); ch != ']'; {
		charset += string(ch)
	}
	return charset
}

func nfa2dfa(nfaGraph *NFAGraph) *DFAGraph {
	startStates := set.New()
	// 用NFA图的起始节点构造DFA的起始节点
	subset := getNextEStates(nfaGraph.start, set.New())
	for subIter := subset.Begin(); subIter.IsValid(); subIter.Next() {
		startStates.Insert(subIter.Value())
	}
	if startStates.Size() == 0 {
		startStates.Insert(nfaGraph.start)
	}
	dfaGraph := NewDFAGraph(nil)
	start := dfaGraph.Get(startStates)
	dfaGraph.start = start
	queue := queue.New()
	finishedStates := set.New()
	// BFS的方式从已找到的起始节点遍历并构建DFA
	queue.Push(start)

	for !queue.Empty() {
		// 对当前节点已添加的边做去重,不放到queue和next里.
		addedNextStates := set.New()
		curState := queue.Pop().(*DFAState)
		for iter := curState.nfaStates.Begin(); iter.IsValid(); iter.Next() {
			nfaState := iter.Value().(*NFAState)
			nextStates := set.New()
			finishedEdges := set.New()
			finishedEdges.Insert(EPSILON)
			for edge, _ := range nfaState.next {
				if finishedEdges.Contains(edge) {
					continue
				}
				finishedEdges.Insert(edge)
				efinishedState := set.New()

				for subiter := curState.nfaStates.Begin(); subiter.IsValid(); subiter.Next() {
					state := subiter.Value().(*NFAState)
					edgeStates, ok := state.next[edge]
					if ok {
						nextStates = nextStates.Union(edgeStates)
					}

					for ssubiter := edgeStates.Begin(); ssubiter.IsValid(); ssubiter.Next() {
						eState := ssubiter.Value().(*NFAState)
						// 添加E可达节点
						if efinishedState.Contains(eState) {
							continue
						}
						nextStates = nextStates.Union(getNextEStates(eState, efinishedState))
						efinishedState.Insert(eState)
					}
				}
				// 将NFA节点列表转化为DFA节点，如果已经有对应的DFA节点就返回，否则创建一个新的DFA节点
				nextDFAstate := dfaGraph.Get(nextStates)
				if !finishedStates.Contains(nextDFAstate) && !addedNextStates.Contains(nextDFAstate) {
					queue.Push(nextDFAstate)
					addedNextStates.Insert(nextDFAstate) // 对queue里的数据做去重
					curState.AddNext(edge, nextDFAstate)
				}
			}
		}
		finishedStates.Insert(curState)
	}
	return dfaGraph
}

func getNextEStates(curState *NFAState, stateSet *set.Set) *set.Set {
	if _, ok := curState.next[EPSILON]; !ok {
		return set.New()
	}
	res := set.New()
	states := curState.next[EPSILON]
	for iter := states.Begin(); iter.IsValid(); iter.Next() {
		state := iter.Value().(*NFAState)
		if stateSet.Contains(state) {
			continue
		}
		res.Insert(state)
		subset := getNextEStates(state, stateSet)
		for subIter := subset.Begin(); subIter.IsValid(); subIter.Next() {
			res.Insert(subIter.Value())
		}
		stateSet.Insert(state)
	}
	return res
}

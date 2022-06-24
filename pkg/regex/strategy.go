package regex

type (
	MatchStrategy interface {
		Match(c byte, edge string) bool
	}

	BaseMatchStrategy struct {
		IsReversed bool
	}

	CharMatchStrategy struct {
		*BaseMatchStrategy
	}

	CharSetMatchStrategy struct {
		*BaseMatchStrategy
	}

	EpsilonMatchStrategy struct {
		*BaseMatchStrategy
	}

	DotMatchStrategy struct {
		*BaseMatchStrategy
	}

	SpaceMatchStrategy struct {
		*BaseMatchStrategy
	}

	// WMatchStrategy 匹配 \w和\W
	WMatchStrategy struct {
		*BaseMatchStrategy
	}

	// DigitalMatchStrategy 匹配数字
	DigitalMatchStrategy struct {
		*BaseMatchStrategy
	}

	MatchStrategyManager struct {
		matchStrategyMap map[string]MatchStrategy
	}
)

var Manager = &MatchStrategyManager{
	matchStrategyMap: map[string]MatchStrategy{
		"\\d":   &DigitalMatchStrategy{&BaseMatchStrategy{false}},
		"\\D":   &DigitalMatchStrategy{&BaseMatchStrategy{true}},
		"\\w":   &WMatchStrategy{&BaseMatchStrategy{false}},
		"\\W":   &WMatchStrategy{&BaseMatchStrategy{true}},
		"\\s":   &SpaceMatchStrategy{&BaseMatchStrategy{false}},
		"\\S":   &SpaceMatchStrategy{&BaseMatchStrategy{true}},
		".":     &DotMatchStrategy{&BaseMatchStrategy{false}},
		EPSILON: &EpsilonMatchStrategy{&BaseMatchStrategy{false}},
		CHAR:    &CharMatchStrategy{&BaseMatchStrategy{false}},
		CHARSET: &CharSetMatchStrategy{&BaseMatchStrategy{false}},
	},
}

func (m *MatchStrategyManager) Get(key string) MatchStrategy {
	// 特殊字符的匹配
	if s, ok := m.matchStrategyMap[key]; ok {
		return s
	}
	if len(key) == 1 {
		return m.matchStrategyMap[CHAR]
	}
	return m.matchStrategyMap[CHARSET]
}

func (s *EpsilonMatchStrategy) Match(c byte, edge string) bool {
	return true
}

func (s *WMatchStrategy) Match(c byte, edge string) bool {
	res := c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' || c >= '0' && c <= '9'
	if s.IsReversed {
		return !res
	}
	return res
}

func (s *SpaceMatchStrategy) Match(c byte, edge string) bool {
	res := c == '\f' || c == '\n' ||
		c == '\r' || c == '\t' || c == ' '
	if s.IsReversed {
		return !res
	}

	return res
}

func (s *DotMatchStrategy) Match(c byte, edge string) bool {
	return c != '\n' && c != '\r'
}

func (s *DigitalMatchStrategy) Match(c byte, edge string) bool {
	res := c >= '0' && c <= '9'
	if s.IsReversed {
		return !res
	}
	return res
}

func (s *CharSetMatchStrategy) Match(c byte, edge string) bool {
	ret := false
	not := edge[0] == '^'
	for i, item := range edge {
		if not {
			continue
		}
		if item == '-' {
			return c >= edge[i-1] && c <= edge[i+1]
		}
		if c == byte(item) {
			ret = true
			break
		}
	}
	if not {
		return !ret
	}
	return ret
}

func (s *CharMatchStrategy) Match(c byte, edge string) bool {
	return edge[0] == c
}

func (s *BaseMatchStrategy) Match(c byte, edge string) bool {
	return false
}

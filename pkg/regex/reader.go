package regex

type Reader struct {
	runes []rune
	cur   int
}

func NewReader(source string) *Reader {
	return &Reader{runes: []rune(source)}
}

func (r *Reader) Peek() rune {
	if r.cur == len(r.runes) {
		return 0
	}
	return r.runes[r.cur]
}

func (r *Reader) Next() rune {
	if r.cur == len(r.runes) {
		return 0
	}
	ret := r.runes[r.cur]
	r.cur++
	return ret
}

func (r *Reader) hasNext() bool {
	return r.cur < len(r.runes)
}

func (r *Reader) getSubRegex(other *Reader) string {
	cntParam := 1
	regex := ""
	for other.hasNext() {
		ch := other.Next()
		if ch == '(' {
			cntParam++
		} else if ch == ')' {
			cntParam--
			if cntParam == 0 {
				break
			}
		}
		regex += string(ch)
	}
	return regex
}

func (r *Reader) getRemainRegex(other *Reader) string {
	regex := ""
	for other.hasNext() {
		regex += string(other.Next())
	}
	return regex
}

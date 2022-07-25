package py

func Optimize(ir []*SSA) []*SSA {
	ret := make([]*SSA, 0)

	fetch := func(n int) (string,
		any, any) {
		if n < len(ir) {
			return ir[n].Action,
				ir[n].Arg1, ir[n].Arg2
		}
		return "", nil, nil
	}

	pushSSA := func(a string, b, c any) {
		ssa := &SSA{
			Action: a,
			Arg1:   b,
			Arg2:   c,
		}
		ret = append(ret, ssa)
	}

	index := 0
	for index < len(ir) {
		op1, a1, b1 := fetch(index)
		op2, a2, b2 := fetch(index + 1)
		op3, a3, _ := fetch(index + 2)

		if op1 == "mov" && a1 == b1 {
			index += 1
			continue
		}

		if op1 == "mov" && op2 == "mov" && a1 == b2 {
			index += 2
			pushSSA("mov", a2, b1)
			continue
		}

		if op1 == "push" && op2 == "pop" {
			index += 2
			pushSSA("mov", a2, a1)
			continue
		}

		if op1 == "push" && op3 == "pop" &&
			op2 != "push" && op2 != "pop" {
			if a2 != a3 {
				index += 3
				pushSSA("mov", a3, a1)
				pushSSA(op2, a2, b2)
				continue
			}
		}

		index++
		pushSSA(op1, a1, b1)
	}

	return ret
}

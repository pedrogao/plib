package rt

//go:nosplit
//go:noescape
//goland:noinspection GoUnusedParameter
func __isspace(ch byte) (ret byte)

// < 10000
//go:nosplit
//go:noescape
//goland:noinspection GoUnusedParameter
func __u32toa_small(out *byte, val uint32) (ret int)

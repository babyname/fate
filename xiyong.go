package fate

func filterXiYong(yong string, cs ...*Character) (b bool) {
	for _, c := range cs {
		if c.WuXing == yong {
			return true
		}
	}
	return false
}

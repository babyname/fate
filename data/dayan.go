package data

var DaYanList [81]DaYan

type DaYan struct {
	Number  int
	Sex     bool //male(0),female(1)
	Lucky   string
	NineSky string
	Comment string
}

func init() {
	DaYanList = [81]DaYan{
		0: {
			Number:  0,
			Sex:     false,
			Lucky:   "吉",
			NineSky: "太极之数",
			Comment: "太极之数，万物开泰，生发无穷，利禄亨通。",
		},
	}
}

package room

type Deck struct {
	Name    string
	Cards   []string
	Numeric bool
}

var (
	Fibonacci = Deck{Name: "fibonacci", Numeric: true, Cards: []string{"0", "1", "2", "3", "5", "8", "13", "21", "?", "coffee"}}
	TShirt    = Deck{Name: "tshirt", Cards: []string{"XS", "S", "M", "L", "XL", "?", "coffee"}}
)

func DeckByName(name string) (Deck, bool) {
	switch name {
	case Fibonacci.Name:
		return Fibonacci, true
	case TShirt.Name:
		return TShirt, true
	}
	return Deck{}, false
}

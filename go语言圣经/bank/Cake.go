package bank


type Cake struct{ state string }

func baker(cooked chan<- *Cake) {
	for {
		//Cakes会被严格地顺序访问，先是baker gorouine，然后是icer gorouine：
		cake := new(Cake)
		cake.state = "cooked"
		cooked <- cake // baker never touches this cake again
	}
}

func icer(iced chan<- *Cake, cooked <-chan *Cake) {
	for cake := range cooked {
		cake.state = "iced"
		iced <- cake // icer never touches this cake again
	}
}
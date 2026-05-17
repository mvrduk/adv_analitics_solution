package ads

type (
	User struct {
		Country string
		Browser string
	}

	Campaign struct {
		ClickUrl  string
		Price     float64
		Targeting Targeting
	}

	Targeting struct {
		Browser string
		Country string
	}

	filterFunc func(in *[]Campaign, u *User) (out *[]Campaign)
)

func filterBrowser(in []*Campaign, u *User) []*Campaign {
	for i := len(in) - 1; i >= 0; i-- {
		if len(in[i].Targeting.Browser) == 0 {
			continue
		}

		if in[i].Targeting.Browser == u.Browser {
			continue
		}

		in[i] = in[0]
	}

	return in
}

func filterByCountry(in []*Campaign, u *User) []*Campaign {
	for i := len(in) - 1; i >= 0; i-- {
		if len(in[i].Targeting.Country) == 0 {
			continue
		}

		if in[i].Targeting.Country == u.Country {
			continue
		}

		in[i] = in[0]
	}

	return in
}

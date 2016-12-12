package main

type currencyList struct {
	Currencies currencies `json:"currencies"`
}

type currencies []string

func getCurList() currencyList {
	currencies := currencies{
		"DZD",
		"AUD",
		"VEF",
		"BHD",
		"BWP",
		"BRL",
		"BND",
		"CAD",
		"CLP",
		"CNY",
		"COP",
		"CZK",
		"DKK",
		"HUF",
		"ISK",
		"INR",
		"IDR",
		"IRR",
		"ILS",
		"JPY",
		"KZT",
		"KRW",
		"KWD",
		"LYD",
		"MYR",
		"MUR",
		"MXN",
		"PEN",
		"NPR",
		"NZD",
		"NOK",
		"UYU",
		"PHP",
		"PKR",
		"PLN",
		"QAR",
		"RUB",
		"SAR",
		"SGD",
		"SZAR",
		"LKR",
		"SEK",
		"CHF",
		"TND",
		"THB",
		"TTD",
		"AED",
		"GBP",
		"USD",
		"EUR",
		"OMR",
	}

	return currencyList{currencies}
}

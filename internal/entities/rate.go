package entities

type Rate struct {
	RAW     map[string]map[string]RateValues
	DISPLAY map[string]map[string]RateValues
}

type RateValues struct {
	CHANGE24HOUR    string
	CHANGEPCT24HOUR string
	OPEN24HOUR      string
	VOLUME24HOUR    string
	VOLUME24HOURTO  string
	LOW24HOUR       string
	HIGH24HOUR      string
	PRICE           string
	SUPPLY          string
	MKTCAP          string
}

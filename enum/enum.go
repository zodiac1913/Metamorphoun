package enum

type PathLocation struct {
	Fonts           string
	Config          string
	Favorites       string
	FavWithQuote    string
	FavWithoutQuote string
	Quotes          string
	ConfigFile      string
	Pictures        string
	Logs            string
	Executable      string //not necessarily the same as BinFileLoc as this is for config creation
	BinFileLoc      string
}

var PathLoc = PathLocation{
	Fonts:           "fonts",
	Config:          "config",
	Favorites:       "favorites",
	FavWithQuote:    "favwithquote",
	FavWithoutQuote: "favwithoutquote",
	Quotes:          "quotes",
	ConfigFile:      "configfile",
	Pictures:        "pictures",
	Logs:            "logs",
	Executable:      "executable", //not necessarily the same as BinFileLoc as this is for config creation
	BinFileLoc:      "binfileloc",
}

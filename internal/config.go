package internal

import "flag"

var Config struct {
	DataDir string
}

func AddGlobalFlags() {
	flag.StringVar(&Config.DataDir, "datadir", "data/", "Path to directory with data")
}

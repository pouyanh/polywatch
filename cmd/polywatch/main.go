package main

import (
	"github.com/pouyanh/polywatch"
	_ "github.com/pouyanh/polywatch/config/viper"
)

func main() {
	if err := polywatch.Start(); err != nil {
		panic(err)
	}
}

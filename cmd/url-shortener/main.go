package main

import (
	"fmt"
	"github/closidx/url-shortener/internal/config"

)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// init logger

	// init storage

	// init router

	// run server
}
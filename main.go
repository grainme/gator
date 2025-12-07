package main

import (
	"fmt"
	"log"

	"github.com/grainme/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	cfg.CurrentUserName = "marouane"
	cfg.SetUser()

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	fmt.Println(cfg.CurrentUserName, cfg.DbUrl)
}

package main

import (
	"log"
	"my_blog_backend/config"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}
}

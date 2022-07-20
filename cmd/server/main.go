package main

import (
	"abc/lobo"
	"flag"
	"fmt"
)

func main() {
	endpoint := []string{"52.11.26.186:2379", "52.11.26.186:22379", "52.11.26.186:3237"}
	port := flag.Int("p", 8888, "Port Number")
	flag.Parse()
	s := lobo.NewServer()

	s.Register(new(UserServiceImpl))

	if err := s.Serve(*port, endpoint); err != nil {
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}

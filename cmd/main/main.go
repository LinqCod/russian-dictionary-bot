package main

import (
	"fmt"
	"russian-dictionary-bot/internal/dict"
)

func main() {

	c := make(chan *dict.Result)
	go dict.ParseWordData("манжета", c)

	res := <-c
	fmt.Println(res)
	//TODO: add results to database
}

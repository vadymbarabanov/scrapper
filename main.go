package main

import (
	"fmt"

	"github.com/vadymbarabanov/scrapper/mydict"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	baseWord := "first"
	dictionary := mydict.Dictionary{}

	err := dictionary.Add(baseWord, "def1")
	checkErr(err)

	word, err2 := dictionary.Search(baseWord)
	checkErr(err2)
	fmt.Println(word)

	err3 := dictionary.Update(baseWord, "def2")
	checkErr(err3)

	word2, err4 := dictionary.Search(baseWord)
	checkErr(err4)
	fmt.Println(word2)

}

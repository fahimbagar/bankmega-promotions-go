package main

import "fmt"

func main() {
	categories := LoadCategory()
	for _, category := range categories {
		if err := category.GetPage(); err != nil {
			panic(err)
		}
		fmt.Println(category)
	}
}
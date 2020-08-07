package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	categories := LoadCategory()
	for i, _ := range categories {
		wg.Add(1)
		go func(i int) {
			err := categories[i].GetPage()
			if err != nil {
				panic(err)
			}
			err = categories[i].GetPageContent()
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	for categoryIndex, _ := range categories {
		wg.Add(1)
	 var contentGroup sync.WaitGroup
		go func(i int) {
			for contentIndex, _ := range *categories[i].Contents {
				contentGroup.Add(1)
				go (*categories[i].Contents)[contentIndex].GetPromotions(&contentGroup)
				if contentIndex + 1 % categories[i].Pages == 0 {
					contentGroup.Wait()
				}
			}
			contentGroup.Wait()
			wg.Done()
		}(categoryIndex)
	}

	wg.Wait()

	var content = map[string]interface{}{}
	for _, category := range categories {
		content[category.Id] = category.Contents
	}

	file, _ := json.MarshalIndent(content, "", " ")
	_ = ioutil.WriteFile("solution.json", file, 0644)
}
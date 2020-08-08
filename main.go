package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

// Main function
func main() {
	var wg sync.WaitGroup
	// Load every category
	categories := LoadCategory()
	// Iterate categories
	for i, _ := range categories {
		wg.Add(1)
		go func(i int) {
			// Load every page
			err := categories[i].GetPage()
			if err != nil {
				panic(err)
			}
			// Load every content url
			err = categories[i].GetPageContent()
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	// Iterate category
	for categoryIndex, _ := range categories {
		wg.Add(1)
	 var contentGroup sync.WaitGroup
		// Crawl every content crawler from each category parallel
		go func(i int) {
			// Iterate content url from each category
			for contentIndex, _ := range *categories[i].Contents {
				// Add wait group for each page
				contentGroup.Add(1)
				// Spawn routine to get promotions from content url
				go (*categories[i].Contents)[contentIndex].GetPromotions(&contentGroup)
				// Wait if routine is spawn more than pages to make sure we don't act like DDOS
				if contentIndex + 1 % categories[i].Pages == 0 {
					contentGroup.Wait()
				}
			}
			// Wait for all left routines
			contentGroup.Wait()
			// State if category content crawling is finished
			wg.Done()
		}(categoryIndex)
	}

	wg.Wait()

	// Print it like the format
	var content = map[string]interface{}{}
	for _, category := range categories {
		content[category.Id] = category.Contents
	}

	file, _ := json.MarshalIndent(content, "", " ")
	_ = ioutil.WriteFile("solution.json", file, 0644)
}
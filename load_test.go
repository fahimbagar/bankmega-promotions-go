package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLoadCategory(t *testing.T) {
	categories := LoadCategory()
	assert.Greater(t, len(categories), 0, "categories should not be empty")
}

func TestCategory_GetPageAndContent(t *testing.T) {
	category := &Category{
		Id: "daily_needs",
		Url: "ajax.promolainnya.php?product=0&subcat=5",
	}
	err := category.GetPage()
	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, category.Pages, "pages from category should not be nil")
	assert.Greater(t, category.Pages, 0, "pages from category should not 0")

	err = (*category).GetPageContent()
	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, category.Contents, "contents from category should not be nil")
	assert.Greater(t, category.Pages, 0, "pages from category should not 0")
	assert.NotNil(t, (*category.Contents)[0].Url, "contents url should not be nil")
	assert.NotNil(t, (*category.Contents)[0].Image, "contents url should not be nil")
	assert.Empty(t, (*category.Contents)[0].Area, "contents url should be empty")
}

func TestContent_GetPromotions(t *testing.T) {
	content := Content{
		Url: "https://www.bankmega.com/promo_detail.php?id=1992",
	}
	var wg sync.WaitGroup
	wg.Add(1)
	content.GetPromotions(&wg)
	assert.NotNil(t, content.Title, "contents from category should not be nil")
	assert.NotNil(t, content.Image, "contents from category should not be nil")
	assert.NotNil(t, content.Area, "contents from category should not be nil")
	assert.NotNil(t, content.Information, "contents from category should not be nil")
}

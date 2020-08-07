package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Category struct {
	Id       string     `json:"id,omitempty"`
	Pages    int        `json:"pages,omitempty"`
	Url      string     `json:"url,omitempty"`
	Contents *[]Content `json:"contents,omitempty"`
}

type Content struct {
	Title       string `json:"title,omitempty"`
	Image       string `json:"image,omitempty"`
	Url         string `json:"url,omitempty"`
	Area        string `json:"area,omitempty"`
	Period      string `json:"period,omitempty"`
	Information string `json:"information,omitempty"`
	File        string `json:"file,omitempty"`
}

var (
	BASE_URL = "https://www.bankmega.com"
	PROMO = "promolainnya.php"
)

func LoadCategory() (categories []Category) {
	doc, err := htmlquery.LoadURL(fmt.Sprintf("%s/%s", BASE_URL, PROMO))
	if err != nil {
		log.Println(fmt.Sprintf("load category: %s", err))
	}

	categoryNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[1]/script/text()]")

	lines := strings.Split(categoryNode.Data, "\n")
	for i := 0; i < len(lines); i++ {
		categoryRegex := regexp.MustCompile(`\$\("#([\w_]+)"\)`)
		urlRegex := regexp.MustCompile(`\.load\("(.*?subcat.*?)"\);`)
		categoryMatch := categoryRegex.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if len(categoryMatch) > 1 {
			i = i + 1
			urlMatch := urlRegex.FindStringSubmatch(strings.TrimSpace(lines[i]))
			if len(urlMatch) > 1 {
				categories = append(categories, Category{Id:categoryMatch[1], Url: urlMatch[1]})
			}
		}
	}

	return categories
}

func (category *Category) GetPage() (err error) {
	doc, err := htmlquery.LoadURL(fmt.Sprintf("%s/%s", BASE_URL, category.Url))
	if err != nil {
		log.Println(fmt.Sprintf("get pages: %s", err))
		return err
	}

	pagingNode := htmlquery.FindOne(doc, "//*[@id=\"paging1\"]")
	pages, err := strconv.Atoi(strings.TrimPrefix(htmlquery.SelectAttr(pagingNode, "title"), "Page 1 of "))
	if err != nil {
		log.Println(fmt.Sprintf("get pages: %s", err))
		return err
	}
	category.Pages = pages

	return err
}

func (category *Category) GetPromotions() (err error) {

}
package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Redirect    string `json:"redirect,omitempty"`
}

var (
	BASE_URL = "https://www.bankmega.com"
	PROMO = "promolainnya.php"
)

func LoadCategory() (categories []Category) {
	doc, err := htmlquery.LoadURL(fmt.Sprintf("%s/%s", BASE_URL, PROMO))
	if err != nil {
		fmt.Println(fmt.Sprintf("load category: %s", err))
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
		fmt.Println(fmt.Sprintf("get pages: %s", err))
		return err
	}

	pagingNode := htmlquery.FindOne(doc, "//*[@id=\"paging1\"]")
	pages, err := strconv.Atoi(strings.TrimPrefix(htmlquery.SelectAttr(pagingNode, "title"), "Page 1 of "))
	if err != nil {
		fmt.Println(fmt.Sprintf("get pages: %s", err))
		return err
	}
	category.Pages = pages
	category.Contents = new([]Content)

	return err
}

func (category *Category) GetPageContent() (err error) {
	for i := 1; i <= category.Pages; i++ {
		doc, err := htmlquery.LoadURL(fmt.Sprintf("%s/%s&page=%d", BASE_URL, category.Url, i))
		if err != nil {
			fmt.Println(fmt.Sprintf("get pages: %s", err))
			return err
		}
		promotionNode := htmlquery.Find(doc, "//*[@id=\"promolain\"]/li/a")
		for _ , n := range promotionNode{
			imageNode := htmlquery.FindOne(n, "//*[@id=\"imgClass\"]")
			url := htmlquery.SelectAttr(n, "href")
			if !strings.HasPrefix(url, "http") {
				url = fmt.Sprintf("%s/%s", BASE_URL, url)
			}
			image := htmlquery.SelectAttr(imageNode, "src")
			if !strings.HasPrefix(image, "http") {
				image = fmt.Sprintf("%s/%s", BASE_URL, image)
			}
			content := Content{
				Url: url,
				Image: image,
			}
			*category.Contents = append(*category.Contents, content)
		}
	}
	return nil
}

func (content *Content) GetPromotions(contentGroup *sync.WaitGroup) {
	for {
		doc, err := htmlquery.LoadURL(content.Url)
		if err != nil {
			fmt.Println(fmt.Sprintf("get pages: %s", err))
			<-time.After(1 * time.Second)
			continue
		}
		titleNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[4]/h3")
		content.Title = htmlquery.InnerText(titleNode)
		areaNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[5]/b")
		content.Area = htmlquery.InnerText(areaNode)
		periodNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[6]")
		content.Period =
			strings.ReplaceAll(
			strings.Join(
				strings.Fields(
					strings.TrimSpace(
						htmlquery.InnerText(periodNode))), " "), "Periode Promo : ", "")
		imageNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[7]/img")
		if imageNode == nil {
			redirectNode := htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[7]/a")
			imageNode = htmlquery.FindOne(doc, "//*[@id=\"contentpromolain2\"]/div[7]/a/img")
			redirect := htmlquery.SelectAttr(redirectNode, "href")
			if !strings.HasPrefix(redirect, "http") {
				redirect = fmt.Sprintf("%s/%s", BASE_URL, redirect)
			}
			content.Redirect = redirect
		}
		image := htmlquery.SelectAttr(imageNode, "src")
		if !strings.HasPrefix(image, "http") {
			image = fmt.Sprintf("%s/%s", BASE_URL, image)
		}
		content.Image = image
		contentGroup.Done()
		return
	}
}
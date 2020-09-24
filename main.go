package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// Scrape uriを設定するとそれをもとにスクレイピングを行う
type Scrape struct {
	uri string
}

// Document スクレイピングするためのドキュメントを返す
func (s *Scrape) Document() *goquery.Document {
	response, err := http.Get(s.URI())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	// HTMLを読み込む
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// HostURL URLからHostを取得する
func (s *Scrape) HostURL() string {
	u, err := url.Parse(s.URI())
	if err != nil {
		log.Fatal(err)
	}

	return u.Scheme + "://" + u.Host
}

// URI uriを取得
func (s *Scrape) URI() string {
	return s.uri
}

// Scraping スクレイピング処理
func (s *Scrape) Scraping() []string {
	// aタグだけを抽出してスクレイピング処理
	doc := s.Document()
	host := s.HostURL()
	var paths []string
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, ok := sel.Attr("href")
		if ok {
			var p string
			if strings.Index(href, "http") == 0 {
				p = href
			} else {
				p = host + href
			}

			// 重複チェック
			if !stringContains(paths, p) {
				paths = append(paths, p)
			}
		}
	})
	return paths
}

// saveXlsx 文字列配列をエクセルファイルにする
func saveXlsx(filename string, array []string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 新しいエクセルファイルを作成して保存
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Links")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Save(filename + ".xlsx")

	// 各セルに値を入れる
	for _, v := range array {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = v
	}
}

// stringContains 重複チェック
func stringContains(array []string, check string) bool {
	for _, v := range array {
		if check == v {
			return true
		}
	}
	return false
}

func main() {
	// コマンド引数の取得
	uri := flag.String("uri", "", "please scraping uri")
	filename := flag.String("fn", "new", "please save filename")
	flag.Parse()

	if len(*uri) == 0 {
		fmt.Println("please scraping uri")
		return
	}

	// ScrapeのURIを指定する
	scrape := Scrape{*uri}

	// スクレイピングを実行してaタグのリンク配列を取得
	paths := scrape.Scraping()

	if paths != nil {
		// 取得したリンク配列をソートする
		sort.Sort(sort.StringSlice(paths))

		// リンク一覧をコンソールに表示
		for _, path := range paths {
			fmt.Println(path)
		}

		// エクセルファイルで保存
		saveXlsx(*filename, paths)
	} else {
		fmt.Println("Not Found Links")
	}
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

func main() {
	content, err := readPdf2("./AVENIR du 21-05-2017.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	return
}

type PageLines struct {
	PageID      int
	ColumnTexts []string
	LineTexts   []string
}

func readPdf1(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.ReadFrom(r.GetPlainText())
	return buf.String(), nil
}

func readPdf2(path string) (string, error) {
	pagesL := []PageLines{}

	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()


	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		pageLines := new(PageLines)
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		pageLines.PageID = pageIndex

		var lastTextStyle pdf.Text

		texts := p.Content().Text

		for _, text := range texts {
			if text.FontSize > 7.4400 {
				continue
			}

			if isSameSentence(text.S, lastTextStyle.S) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				switch text.Font {
				case "Arial-BoldMT":
					line := fmt.Sprintf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
					lastTextStyle = text
					pageLines.ColumnTexts = append(pageLines.ColumnTexts, line)
				case "ArialMT":
					line := fmt.Sprintf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
					lastTextStyle = text
					pageLines.LineTexts = append(pageLines.LineTexts, line)
				}
			}
			pagesL = append(pagesL, *pageLines)
		}
		if pageLines.PageID == 1 {
			log.Printf("for [%d] - len(pageLines.LineTexts) is: %d\nContent is: %s\n------\n", pageLines.PageID, len(pageLines.ColumnTexts), pageLines.ColumnTexts)
		}

	}

	return "", nil
}

func isSameSentence(a, b string) bool {
	length := index(a, b)
	if length == len(a) {
		return true
	}
	return false
}

func index(s1, s2 string) int {
	res := 0
	for i, w := 0, 0; i < len(s2); i += w {
		if i >= len(s1) {
			return res
		}
		runeValue1, width := utf8.DecodeRuneInString(s1[i:])
		runeValue2, width := utf8.DecodeRuneInString(s2[i:])
		if runeValue1 != runeValue2 {
			return res
		}
		w = width
		res = i + w
	}
	return res
}

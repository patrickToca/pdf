package main

import (
	"bytes"
	"fmt"
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
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text.S, lastTextStyle.S) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
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

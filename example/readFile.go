package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

var (
	ref             Lines = make(Lines, 0)
	sentenceCounter uint32
)

func main() {
	content, err := readPdf2("./AVENIR du 21-05-2017.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(len(content))
	return
}

type Reference struct {
	ID         string
	PageNumber int
	LineNumber uint32
	LineText   []string
}

// Lines is a map with a string Key="PageNumber-LineNumber"
type Lines map[string][]*Reference

func (line Lines) AddReference(pageNumber int, lineNumber uint32, lineText []string) error {
	ID := fmt.Sprintf("%d-%d", pageNumber, lineNumber)
	log.Printf("%s\n", ID)

	err := errors.New("Duplicate? This Line already exists!")

	ref, exist := line[ID]
	if exist {
		return err
	}
	line[ID] = append(ref, &Reference{
		ID:         ID,
		PageNumber: pageNumber,
		LineNumber: lineNumber,
		LineText:   lineText,
	})
	return nil
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

	reference := new(Reference)
	lines := make(Lines)

	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {

		reference.PageNumber = pageIndex

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
				reference.LineNumber = atomic.AddUint32(&sentenceCounter, 1)
				reference.LineText = append(reference.LineText, lastTextStyle.S)
				//fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
				if err := lines.AddReference(reference.PageNumber, reference.LineNumber, reference.LineText); err != nil {
					log.Printf("lines.AddReference() error: %s\n", err)
				}
			}
		}
	}
	fmt.Printf("This Line is: %s\n", lines["12-11367"])
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

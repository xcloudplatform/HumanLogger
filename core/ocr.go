package core

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"os"
	"path/filepath"
	"strings"
)

type Ocr struct {
	client     *gosseract.Client
	queue      chan string
	resultChan chan Result
}

type Result struct {
	Text string
	File string
}

func NewOcr() *Ocr {
	return &Ocr{
		client:     gosseract.NewClient(),
		queue:      make(chan string),
		resultChan: make(chan Result),
	}
}

func (o *Ocr) Recognize(imagePath string) {
	o.queue <- imagePath
}

func (o *Ocr) ProcessQueue() {
	for imagePath := range o.queue {
		o.client.SetImage(imagePath)
		text, _ := o.client.Text()
		result := Result{
			Text: text,
			File: imagePath,
		}

		ocrPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + ".txt"
		if err := os.WriteFile(ocrPath, []byte(text), 0644); err != nil {
			fmt.Printf("error attempting to save OCR text to %s\n", ocrPath)

		} else {
			fmt.Printf("OCR text saved to %s\n", ocrPath)
		}

		o.resultChan <- result
	}
}

func (o *Ocr) GetResult() <-chan Result {
	return o.resultChan
}

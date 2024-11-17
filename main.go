package main

import (
	"fmt"
	"log"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	pdfFiles := []string{
		"./pdf-files/A-Philosophy-of-Software-Design.pdf",
		// "./pdf-files/the-art-of-writing-code.pdf",
	}

	for _, pdfFile := range pdfFiles {
		pdfType, err := JudgePDFType(instance, pdfFile)
		if err != nil {
			log.Fatalf("判断 PDF 类型失败: %v", err)
		}
		fmt.Printf("文件%s的类型是：%s\n\n\n", pdfFile, pdfType)
	}
}

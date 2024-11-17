package main

import (
	"fmt"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/enums"
	"github.com/klippa-app/go-pdfium/requests"
	"golang.org/x/exp/rand"
)

const (
	SELECT_SAMPLE_PAGE_NUM        = 10  // 随机选取的页数
	NORMAL_PAGE_RATIO             = 0.8 // 普通页占比
	NORMAL_PAGE_TEXT_FRAGMENT_NUM = 5   // 普通页的文字元素数量
)

type PDFType int

const (
	PDF_TYPE_UNKNOWN PDFType = 0
	PDF_TYPE_NORMAL  PDFType = 1
	PDF_TYPE_SCAN    PDFType = 2
)

func (pdfType PDFType) String() string {
	switch pdfType {
	case PDF_TYPE_NORMAL:
		return "normal"
	case PDF_TYPE_SCAN:
		return "scan"
	default:
		return "unknown"
	}
}

// 新增函数：从 0 到 n-1 的序列中随机选择 m 个数字
func SelectRandomNumbers(n int, m int) []int {
	if m > n {
		m = n // 如果 m 大于 n，则将 m 设置为 n
	}

	// 创建一个包含 0 到 n-1 的序列
	sequence := make([]int, n)
	for i := 0; i < n; i++ {
		sequence[i] = i
	}

	// 随机打乱序列
	rand.Seed(uint64(time.Now().UnixNano())) // 设置随机种子
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})

	// 返回前 m 个随机选择的数字
	return sequence[:m]
}

func JudgePDFType(instance pdfium.Pdfium, inputPath string) (PDFType, error) {
	pdfDoc, err := instance.FPDF_LoadDocument(&requests.FPDF_LoadDocument{
		Path:     &inputPath,
		Password: nil,
	})
	if err != nil {
		return PDF_TYPE_UNKNOWN, fmt.Errorf("无法加载 PDF 文档: %v", err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: pdfDoc.Document,
	})

	pageCountRes, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: pdfDoc.Document,
	})
	if err != nil {
		return PDF_TYPE_UNKNOWN, fmt.Errorf("无法获取 PDF 文档页数: %v", err)
	}

	var totalPageNum = pageCountRes.PageCount
	var normalPageNum = 0

	// 随机选取SELECT_SAMPLE_PAGE_NUM页
	randomPageIndexes := SelectRandomNumbers(totalPageNum, SELECT_SAMPLE_PAGE_NUM)

	for _, i := range randomPageIndexes {
		pdfPage, err := instance.FPDF_LoadPage(&requests.FPDF_LoadPage{
			Document: pdfDoc.Document,
			Index:    i,
		})
		if err != nil {
			return PDF_TYPE_UNKNOWN, fmt.Errorf("无法加载 PDF 文档页面: %v", err)
		}

		// 遍历页面中的对象
		objectCountRes, err := instance.FPDFPage_CountObjects(&requests.FPDFPage_CountObjects{
			Page: requests.Page{ByReference: &pdfPage.Page},
		})
		if err != nil {
			return PDF_TYPE_UNKNOWN, fmt.Errorf("无法获取 PDF 文档页面对象数量: %v", err)
		}

		var onePageTextFragmentNum = 0
		for j := 0; j < objectCountRes.Count; j++ {
			objRes, err := instance.FPDFPage_GetObject(&requests.FPDFPage_GetObject{
				Page:  requests.Page{ByReference: &pdfPage.Page},
				Index: j,
			})
			if err != nil {
				return PDF_TYPE_UNKNOWN, fmt.Errorf("无法获取 PDF 文档页面对象: %v", err)
			}

			objTypeRes, err := instance.FPDFPageObj_GetType(&requests.FPDFPageObj_GetType{
				PageObject: objRes.PageObject,
			})
			if err != nil {
				return PDF_TYPE_UNKNOWN, fmt.Errorf("无法获取 PDF 文档页面对象类型: %v", err)
			}

			if objTypeRes.Type == enums.FPDF_PAGEOBJ_TEXT {
				onePageTextFragmentNum++
			}

			if onePageTextFragmentNum >= NORMAL_PAGE_TEXT_FRAGMENT_NUM {
				break
			}

		}
		// fmt.Printf("第%d页的总元素数:%d text元素数：%d\n", i, objectCountRes.Count, onePageTextFragmentNum)

		if onePageTextFragmentNum >= NORMAL_PAGE_TEXT_FRAGMENT_NUM {
			normalPageNum++
		}
	}

	if float64(normalPageNum)/float64(len(randomPageIndexes)) >= NORMAL_PAGE_RATIO {
		return PDF_TYPE_NORMAL, nil
	}

	return PDF_TYPE_SCAN, nil
}

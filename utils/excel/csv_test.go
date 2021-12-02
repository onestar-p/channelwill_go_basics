package excel

import (
	"fmt"
	"testing"
)

func TestCsv(t *testing.T) {

	filePath := "./csv_test.csv"
	csv := NewCsv(filePath)
	header := []string{
		"编号", "姓名", "年龄",
	}
	csv.SetHeader(header)
	datas := [][]string{
		{"123", "Golang", "18"},
		{"123", "Golang", "18"},
	}
	csv.AppendDatas(datas...)
	if err := csv.Export(); err != nil {
		panic(err)
	}

	for i := 1; i < 100; i++ {
		data := []string{fmt.Sprintf("%d", i), fmt.Sprintf("Golang%d", i), "18"}
		csv.AppendDatas(data)
	}

	if err := csv.AdditionalExport(); err != nil {
		panic(err)
	}

}

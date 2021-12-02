package excel

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Csv struct {
	filePath string // 文件保存路径
	header   []string
	datas    [][]string
}

func NewCsv(filePath string) *Csv {
	return &Csv{
		filePath: filePath,
	}
}

func (c *Csv) SetHeader(header []string) {
	c.header = header
}

func (c *Csv) AppendDatas(data ...[]string) {
	c.datas = append(c.datas, data...)
}

// 导出
// @params filePath
func (c *Csv) Export() error {
	file, err := os.Create(c.filePath)
	if err != nil {
		return fmt.Errorf("cannot Create file path err: %v", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write(c.header); err != nil {
		return err
	}
	if err := writer.WriteAll(c.datas); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// 追加数据
func (c *Csv) AdditionalExport() error {
	if _, err := os.Stat(c.filePath); err != nil {
		if ok := os.IsNotExist(err); ok {
			return err
		}
	}

	file, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("cannot Create file path err: %v", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(c.datas); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

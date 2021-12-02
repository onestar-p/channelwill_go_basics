package excel

type ExcelInterfaced interface {
	SetHeader(header []string)
	AppendData(data []string)
}

// 表格导出
type Excel struct {
}

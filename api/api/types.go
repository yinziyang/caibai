package handler

// Response 结构定义
type Response struct {
	JsonResult  string      `json:"JsonResult"`
	JsonMessage JsonMessage `json:"JsonMessage"`
	JsonData    []JsonData  `json:"JsonData"`
}

type JsonMessage struct {
	MessageIndex string `json:"MessageIndex"`
	Remark       string `json:"Remark"`
	MessageInfo  string `json:"MessageInfo"`
}

type JsonData struct {
	SQLBuilderID string    `json:"SQLBuilderID"`
	Field        []Field   `json:"FIELD"`
	Row          []RowData `json:"ROW"`
}

type Field struct {
	AttrName  string `json:"attrname"`
	FieldType string `json:"fieldtype"`
	Width     string `json:"WIDTH"`
}

type RowData struct {
	FKindName  string `json:"FKIND_NAME"`
	FPriceBase string `json:"FPRICE_BASE"`
	FNewTime   string `json:"FNEWTIME"`
	FTopRemark string `json:"FTOP_REMARK"`
	FRemark    string `json:"FREMARK"`
}

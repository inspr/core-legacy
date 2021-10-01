package model

type Request struct {
	Op1 int `json:"op1"`
	Op2 int `json:"op2"`
}

type Response struct {
	Result int `json:"result"`
}

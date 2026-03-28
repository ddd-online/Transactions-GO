package models

type Result struct {
	Code int         `json:"code"`           // return code
	Msg  string      `json:"msg,omitempty"`  // message
	Data interface{} `json:"data,omitempty"` // data object
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
		Data: nil,
	}
}

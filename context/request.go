package context

import "encoding/json"

type Request struct {
	// 请求资源
	Uri string `json:"uri"`
	// 请求内容,json字符串
	Body interface{} `json:"body"`
}


func (req Request) jsonMarshal() ([]byte, error) {
	return json.Marshal(req.Body)
}

func NewRequest(msg []byte) (req Request, err error) {
	err = json.Unmarshal(msg, &req)
	return
}

package context

import "encoding/json"

// Request
type Request struct {
	// 请求资源
	Uri string `json:"uri"`
	// 请求头，用于存放如token此类的数据
	Header map[string]string `json:"header"`
	// 请求内容,json字符串
	Body interface{} `json:"body"`
}

// Unmarshal use pointer to get body
func (req Request) Unmarshal(v interface{}) error {
	// marshal json
	b, err := json.Marshal(req.Body)
	if err != nil {
		return err
	}
	// unmarshal json
	return json.Unmarshal(b, v)
}

//
func newRequest(msg []byte) (req Request, err error) {
	err = json.Unmarshal(msg, &req)
	return
}

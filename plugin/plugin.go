package plugin

import (
	"context"
	"encoding/json"
)

func MockPluginCall(ctx context.Context, req *CallPluginToolReq) (*CallPluginToolResp, error) {
	if req.PluginId == 0 {
		var input map[string]interface{}
		_ = json.Unmarshal(req.Request, &input)
		for k, v := range req.Secrets {
			input[k] = v
		}
		output, _ := json.Marshal(input)
		return &CallPluginToolResp{Response: output}, nil

	}
	return &CallPluginToolResp{Response: req.Request}, nil
}

type CallPluginToolReq struct {
	PluginId int64             `thrift:"plugin_id,1,required" frugal:"1,required,i64" json:"plugin_id"`
	ToolId   int64             `thrift:"tool_id,2,required" frugal:"2,required,i64" json:"tool_id"`
	Secrets  map[string]string `thrift:"secrets,3,optional" frugal:"3,optional,map<string:string>" json:"secrets,omitempty"`
	Request  []byte            `thrift:"request,4,required" frugal:"4,required,binary" json:"request"`
}

type CallPluginToolResp struct {
	Code     int64  `thrift:"code,1,required" frugal:"1,required,i64" json:"code"`
	Msg      string `thrift:"msg,2,required" frugal:"2,required,string" json:"msg"`
	Response []byte `thrift:"response,3,required" frugal:"3,required,binary" json:"response"`
	HttpCode int64  `thrift:"http_code,4,required" frugal:"4,required,i64" json:"http_code"`
}

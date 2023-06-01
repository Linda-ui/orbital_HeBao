package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client/genericclient"
)

type requiredParams struct {
	Method    string `json:"method,required"`
	BizParams string `json:"biz_params,required"`
}

var SvcMap = make(map[string]genericclient.Client)

func Gateway(ctx context.Context, c *app.RequestContext) {

}

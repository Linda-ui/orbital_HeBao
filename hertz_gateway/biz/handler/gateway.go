package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type requiredParams struct {
	Method    string `json:"method,required"`
	BizParams string `json:"biz_params,required"`
}

var SvcMap = make(map[string]genericclient.Client)

type error struct {
	msg string
}

func Gateway(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")
	cli, ok := SvcMap[svcName]
	if !ok {
		c.JSON(http.StatusOK, error{msg: "bad request"})
		return
	}

	var params requiredParams
	if err := c.BindAndValidate(&params); err != nil {
		hlog.Error(err)
		c.JSON(http.StatusOK, error{msg: "server method not found"})
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%s/%s", svcName, params.Method),
		bytes.NewBuffer([]byte(params.BizParams)),
	)
	if err != nil {
		hlog.Warnf("new http request failed: %v", err)
		c.JSON(http.StatusOK, error{msg: "request server fail"})
		return
	}

	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		hlog.Errorf("convert request failed: %v", err)
		c.JSON(http.StatusOK, error{msg: "server handle fail"})
		return
	}

	resp, err := cli.GenericCall(ctx, "", customReq)
	respMap := make(map[string]interface{})
	if err != nil {
		hlog.Errorf("GenericCall err:%v", err)
		bizErr, ok := kerrors.FromBizStatusError(err)
		if !ok {
			c.JSON(http.StatusOK, error{msg: "server handle fail"})
			return
		}
		respMap["err_code"] = bizErr.BizStatusCode()
		respMap["err_message"] = bizErr.BizMessage()
		c.JSON(http.StatusOK, respMap)
		return
	}

	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		c.JSON(http.StatusOK, error{msg: "server handle fail"})
		return
	}

	realResp.Body["err_code"] = 0
	realResp.Body["err_message"] = "ok"
	c.JSON(http.StatusOK, realResp.Body)
}

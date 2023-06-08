package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/errors"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/idl_mapping"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type requiredParams struct {
	BizParams string `json:"biz_params"`
}

var SvcMap = idl_mapping.DynamicMap{}

func Gateway(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")
	methodName := c.Param("method")

	cli, ok := SvcMap.GetClient(svcName)
	if !ok {
		c.JSON(http.StatusOK, errors.New(errors.Err_ServerNotFound))
		return
	}

	var params requiredParams
	if err := c.BindAndValidate(&params); err != nil {
		hlog.Errorf("binding error: %v", err)
		c.JSON(http.StatusOK, errors.New(errors.Err_BadRequest))
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%s/%s", svcName, methodName),
		bytes.NewBuffer([]byte(params.BizParams)),
	)
	if err != nil {
		hlog.Errorf("new http request failed: %v", err)
		c.JSON(http.StatusOK, errors.New(errors.Err_RequestServerFail))
		return
	}

	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		hlog.Errorf("convert request failed: %v", err)
		c.JSON(http.StatusOK, errors.New(errors.Err_BadRequest))
		return
	}

	resp, err := cli.GenericCall(ctx, "", customReq)
	respMap := make(map[string]interface{})
	if err != nil {
		hlog.Errorf("generic call err: %v", err)
		bizErr, ok := kerrors.FromBizStatusError(err)
		if !ok {
			c.JSON(http.StatusOK, errors.New(errors.Err_ServerMethodNotFound))
			return
		}
		respMap["err_code"] = bizErr.BizStatusCode()
		respMap["err_message"] = bizErr.BizMessage()
		c.JSON(http.StatusOK, respMap)
		return
	}

	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		c.JSON(http.StatusOK, errors.New(errors.Err_ServerHandleFail))
		return
	}

	realResp.Body["err_code"] = 0
	realResp.Body["err_message"] = "success"
	c.JSON(http.StatusOK, realResp.Body)
}

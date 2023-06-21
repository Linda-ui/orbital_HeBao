package handler

import (
	//"bytes"
	"context"
	"encoding/json"

	//"fmt"
	"net/http"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/errors"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/idl_mapping"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	//"github.com/cloudwego/kitex/pkg/generic"
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

	hlog.Info(methodName)
	// req is of correct format to JSONThriftGeneric client
	req := params.BizParams
	hlog.Info(req)
	resp, err := cli.GenericCall(ctx, methodName, req)
	hlog.Info(err)

	//respMap is for when resp is unavailable
	respMap := make(map[string]interface{})
	if err != nil {
		hlog.Errorf("generic call err: %v", err)
		bizErr, ok := kerrors.FromBizStatusError(err)
		if !ok {
			// at here

			c.JSON(http.StatusOK, errors.New(errors.Err_ServerMethodNotFound))
			return
		}
		respMap["err_code"] = bizErr.BizStatusCode()
		respMap["err_message"] = bizErr.BizMessage()
		c.JSON(http.StatusOK, respMap)
		return
	}

	realResp, ok := resp.([]byte)
	if !ok {
		c.JSON(http.StatusOK, errors.New(errors.Err_ServerHandleFail))
		return
	}
	jsonResp := make(map[string]string)
	error := json.Unmarshal(realResp, &jsonResp)
	if error != nil {
		c.JSON(http.StatusOK, errors.New(errors.Err_ResponseUnableParse))
	}

	jsonResp["err_code"] = "0"
	jsonResp["err_message"] = "success"
	c.JSON(http.StatusOK, jsonResp)
}

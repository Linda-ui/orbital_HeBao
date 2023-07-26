package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/errors"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/entity"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type apiGateway struct {
	idlMapManager entity.MapManager
	params        requiredParams
}

type requiredParams struct {
	BizParams string `json:"biz_params,required"`
}

func NewGateway(idlMapManager entity.MapManager) *apiGateway {
	return &apiGateway{
		idlMapManager: idlMapManager,
	}
}

func (gateway *apiGateway) Handler(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")
	methodName := c.Param("method")

	errSender := errors.New()

	cli, ok := gateway.idlMapManager.GetClient(svcName)
	if !ok {
		c.JSON(http.StatusOK, errSender.JSONEncode(entity.Err_ServerNotFound))
		return
	}

	if err := c.BindAndValidate(&gateway.params); err != nil {
		hlog.Errorf("binding error: %v", err)
		c.JSON(http.StatusOK, errSender.JSONEncode(entity.Err_BadRequest))
		return
	}

	req := gateway.params.BizParams
	// req is of type string. It is a valid type to be passed in to the GenericCall method.
	resp, err := cli.GenericCall(ctx, methodName, req)

	respMap := make(map[string]interface{})
	if err != nil {
		splitMsg := strings.SplitN(err.Error(), ": ", 2)
		respMap["error_category"] = splitMsg[0]
		respMap["error_details"] = splitMsg[1]
		c.JSON(http.StatusOK, respMap)
		return
	}

	realResp, ok := resp.(string)
	if !ok {
		c.JSON(http.StatusOK, errSender.JSONEncode(entity.Err_ServerHandleFail))
		return
	}

	// Unmarshalling the response to map to append extra data.
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(realResp), &jsonMap)
	if err != nil {
		c.JSON(http.StatusOK, errSender.JSONEncode(entity.Err_ResponseUnableParse))
	}
	jsonMap["err_code"] = 0
	jsonMap["err_message"] = "success"

	c.JSON(http.StatusOK, jsonMap)
}

namespace go echo

struct EchoReq {
  1: string msg
}

struct EchoResp {
  1: string msg
}

service EchoSvc {
  EchoResp EchoMethod(1: EchoReq req)( api.post = '/echo/echomethod', api.param = 'true' )
}
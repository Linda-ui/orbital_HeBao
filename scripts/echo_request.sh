# Normal request 1
# returns {"err_code":0,"err_message":"success","msg":"hello"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"msg\":\"hello\"}"}' \
    -w '\n' 

# Normal request 2
# returns {"err_code":0,"err_message":"success","msg":"goodbye"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"msg\":\"goodbye\"}"}' \
    -w '\n' 

# Bad request 1: request with INVALID json body
# returns {"err_code":10001,"err_message":"bad request"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"XXX\":\"hello\"}"' \
    -w '\n' 

# Bad request 2: request with EMPTY json body
# returns {"err_code":10001,"err_message":"bad request"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -w '\n' 

# Server not found. The server echoXXX does not exist.
# returns {"err_code":10002,"err_message":"server not found"}
curl -X POST http://localhost:8080/gateway/echoXXX/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"msg\":\"hello\"}"}' \
    -w '\n' 

# Server method not found. The method EchoMethodXXX does not exist.
# returns {"err_code":10003,"err_message":"server method not found"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethodXXX \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"msg\":\"hello\"}"}' \
    -w '\n' 

# Problem to be solved:
# the name of the fields in the request body is not validated due to the implementation
# of Kitex's JSON mapping generic call feature. A mechanism needs to be incorporated in 
# the gateway to validate the request body's field names.

# expected: {"err_code":10001,"err_message":"bad request"}
# got: {"err_code":0,"err_message":"success","msg":""} (empty string is the default value for msg)

# curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
#     -H "Content-Type: application/json" \
#     -d '{"biz_params":"{\"XXX\":\"hello\"}"}' \
#     -w '\n' 

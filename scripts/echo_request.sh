# Normal request 1
# returns {"err_code":0,"err_message":"success","msg":"hello"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n' 

# Normal request 2
# returns {"err_code":0,"err_message":"success","msg":"Lorem ipsum dolor sit amet, consectetur adipiscing elituiofdsnfmdlauidansdfjusarhenwavjchx zmxc vkjfdhwjgndmVcuohqewLorem ipsum dolor sit amet, consectetur adipiscing elituiofdsnfmdlauidansdfjusarhenwavjchx zmxc vkjfdhwjgndmVcuohqew"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -w '\n' \
    -d '{"msg": "Lorem ipsum dolor sit amet, consectetur adipiscing elituiofdsnfmdlauidansdfjusarhenwavjchx zmxc vkjfdhwjgndmVcuohqewLorem ipsum dolor sit amet, consectetur adipiscing elituiofdsnfmdlauidansdfjusarhenwavjchx zmxc vkjfdhwjgndmVcuohqew"}'

# Bad request: request with INVALID json body
# returns {"err_code":10001,"err_message":"bad request"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"msg": "hello}' \
    -w '\n' 

# Server not found. The server echoXXX does not exist.
# returns {"err_code":10002,"err_message":"server not found"}
curl -X POST http://localhost:8080/gateway/echoXXX/EchoMethod \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n' 

# Server method not found. The method EchoMethodXXX does not exist.
# returns {"err_code":10003,"err_message":"server method not found"}
curl -X POST http://localhost:8080/gateway/echo/EchoMethodXXX \
    -H "Content-Type: application/json" \
    -d '{"msg":"hello"}' \
    -w '\n' 

# Normal request 1: both positive numbers
# returns {"err_code":0,"err_message":"success","sum":6}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"firstNum\":2,\"secondNum\":4}"}' \
    -w '\n'

# Normal request 2: with a negative number
# returns {"err_code":0,"err_message":"success","sum":-6000}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"biz_params":"{\"firstNum\":-10000,\"secondNum\":4000}"}' \
    -w '\n'

# Problem to be solved:
# the value and type of the fields in the request body is not validated due to the implementation
# of Kitex's JSON mapping generic call feature. A mechanism needs to be incorporated in 
# the gateway to validate the request body's field names.

# expected: {"err_code":10001,"err_message":"bad request"}
# got: {"err_code":0,"err_message":"success","sum":4000} (0 is the default value for firstNum)

# curl -X POST http://localhost:8080/gateway/sum/SumMethod \
#     -H "Content-Type: application/json" \
#     -d '{"biz_params":"{\"firstNum\":\"XXX\",\"secondNum\":4000}"}' \
#     -w '\n'
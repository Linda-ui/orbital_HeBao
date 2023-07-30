# Normal request 1: both positive numbers
# returns {"err_code":0,"err_message":"success","sum":6}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"firstNum": 4, "secondNum": 2}' \
    -w '\n'

# Normal request 2: with a negative number
# returns {"err_code":0,"err_message":"success","sum":-6000}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"firstNum": -10000, "secondNum": 4000}' \
    -w '\n'

# Normal request 3: with a decimal number
# returns {"err_code":0,"err_message":"success","sum":14}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"firstNum": 5.4, "secondNum": 9}' \
    -w '\n'

# Bad request: integer overflow
# returns {"error_category":"remote or network error[remote]","error_details":"biz error: Overflow: sum is too large to be represented."}
curl -X POST http://localhost:8080/gateway/sum/SumMethod \
    -H "Content-Type: application/json" \
    -d '{"firstNum": 1, "secondNum": 9223372036854775807}' \
    -w '\n'

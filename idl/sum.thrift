namespace go sum

struct SumReq {
  1: i64 firstNum
  2: i64 secondNum
}

struct SumResp {
  1: i64 sum
}

service SumSvc {
  SumResp SumMethod(1: SumReq req)
}
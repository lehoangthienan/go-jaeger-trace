# EXAMPLE JAEGER TRACE

![jaeger-go-trace](https://dytvr9ot2sszz.cloudfront.net/wp-content/uploads/2020/11/1200x628_Golang_Jaeger-Tracing_2-min-1024x536.jpg)

## SETUP

```
make setup
make init
``

## HOW TO START
```
make dev-service-a
make dev-service-b
http://localhost:3000/go
```


## JAEGER CLIENT
- Call any http request from service-a or service-b
```
http://localhost:16686
```
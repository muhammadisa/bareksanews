<h1 align="center">Welcome to bareksanews backend source code 👋</h1>
<p>
  <img alt="documentation: yes" src="https://img.shields.io/badge/Documentation-Yes-green.svg" />
  <img alt="maintained: yes" src="https://img.shields.io/badge/Maintained-Yes-green.svg" />
</p>


>This project is using go-kit as the standard library for microservice style and adapting clean architecture like uncle bob does you can check for detail here is the link
>
>uncle bob article about clean architecture https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html



# About

### Bareksanews Service

This service used for news, tag, topic management for newses about bareksa

### Tech Stack

- programming language : golang : https://golang.org/
- standard lib : go-kit : https://github.com/go-kit/kit.git
- circuit breaker : hystrix : https://github.com/afex/hystrix-go
- proxy gen : grpc-gateway : https://github.com/grpc-ecosystem/grpc-gateway.git
- tracer engine : zipkin : https://zipkin.io/
- rpc : grpc : https://grpc.io/
- unit test : testify : https://github.com/stretchr/testify
- web token : jwt : https://jwt.io/
- database sql : mysql : https://www.mysql.com/
- cache : redis : https://redis.io/

### To Start

- move to prerequisite directory in this project
- then docker-compose up -d
- then go run main.go server live in 8010 port
- postman collection inside prerequisite directory

###### Bareksanews Service
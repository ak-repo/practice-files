# microservice-demo
microservice using gRPC and HTTP



https://github.com/ak-repo/microservice-demo.git



├── docker-compose.yml
├── go.mod
├── go.work
├── pkg
│   ├── config
│   │   └── config.go
│   ├── go.mod
│   ├── go.sum
│   ├── jwt
│   │   └── jwt.go
│   └── logger
│       └── logger.go
├── README.md
└── services
    ├── auth-service
    │   ├── cmd
    │   │   └── main.go
    │   ├── Dockerfile
    │   ├── go.mod
    │   └── go.sum
    ├── order-service
    │   ├── cmd
    │   │   └── main.go
    │   ├── Dockerfile
    │   └── go.mod
    └── product-service
        ├── cmd
        │   └── main.go
        ├── Dockerfile
        ├── go.mod
        └── go.sum

12 directories, 20 files
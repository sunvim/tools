package constant

const DockerFile = `
FROM ankrnetwork/alpine:v1.0.0 AS builder
LABEL stage=builder
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /root/
RUN mkdir config
COPY --from=builder /go/src/app/app .
CMD ["/root/app"]
`

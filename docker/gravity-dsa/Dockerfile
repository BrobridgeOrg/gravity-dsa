FROM golang:1.15.3-alpine3.12 AS builder
WORKDIR $GOPATH/src/github.com/BrobridgeOrg/gravity-dsa/
COPY . .

RUN apk add --update git
RUN CGO_ENABLED=0 GO111MODULE=on go build -installsuffix cgo -o /gravity-dsa

FROM alpine:3.12
COPY --from=builder /gravity-dsa /usr/local/bin/gravity-dsa
COPY ./config/config.toml /usr/local/bin/config.toml
EXPOSE 44444
CMD ["/usr/local/bin/gravity-dsa"]

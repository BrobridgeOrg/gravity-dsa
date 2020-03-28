FROM golang:1.12.1-alpine AS builder
WORKDIR $GOPATH/src/github.com/BrobridgeOrg/gravity-dsa/
COPY . .

RUN apk add --update git
RUN CGO_ENABLED=0 GO111MODULE=on go build -installsuffix cgo -o /gravity-dsa

FROM alpine:3.8
COPY --from=builder /gravity-dsa /usr/local/bin/gravity-dsa
COPY ./config/config.toml /usr/local/bin/config.toml
EXPOSE 44444
CMD ["/usr/local/bin/gravity-dsa"]

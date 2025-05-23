# build stage
FROM golang:1.23.9 AS builder

WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && make clean test build


# final stage
FROM alpine
LABEL name=firewall-policy-api
LABEL url=https://github.com/jeessy2/firewall-policy-api

WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
COPY --from=builder /app/firewall-policy-api /app/firewall-policy-api
EXPOSE 80
ENTRYPOINT ["/app/firewall-policy-api"]

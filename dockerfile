# 1、构建阶段
FROM golang:1.24.5 AS builder
 
WORKDIR /go/src/app
COPY . /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -o beianApp ./cmd/web 

# 2、运行阶段
FROM scratch

LABEL version="1.1" \
description="多域名备案展示系统" \
maintainer="teamlet@golangle.net"

WORKDIR /
COPY --from=builder /go/src/app/conf /conf
COPY --from=builder /go/src/app/template /template
COPY --from=builder /go/src/app/log /log
COPY --from=builder /go/src/app/beianApp .

VOLUME /conf
VOLUME /log
VOLUME /template

EXPOSE 8901
CMD ["./beianApp"]

FROM --platform=$TARGETPLATFORM golang:alpine as builder

WORKDIR /build

COPY . .

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

RUN go mod tidy
RUN go build -o bin/app learning

FROM --platform=$TARGETPLATFORM alpine

WORKDIR /chat

COPY --from=builder /build/bin ./
COPY --from=builder /build/dist ./dist

RUN ls -al

EXPOSE 3000

ENTRYPOINT ["./app"]

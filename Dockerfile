FROM golang:1.23.5-alpine3.21 AS builder
WORKDIR /usr/local/app/
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o main .

FROM scratch
COPY --from=builder /usr/local/app/main ./main
CMD ["./main"]

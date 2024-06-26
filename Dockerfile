FROM golang:alpine as builder

WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/app ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/app /app/asyncSender

EXPOSE 8080
CMD ["./asyncSender"]
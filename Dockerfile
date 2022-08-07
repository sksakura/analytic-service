FROM golang:1.17-buster as builder

WORKDIR /app

COPY cmd cmd
COPY config config
COPY internal internal
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o analytic-service ./cmd/app/main.go

FROM alpine:3.15.4
WORKDIR /app

COPY --from=builder /app/analytic-service /app/analytic-service
CMD ["/app/analytic-service"]
ENV DB_CONNECTION_STRING="postgres://team21:mNgdxITbhVGd@91.185.93.23:5432/postgres"
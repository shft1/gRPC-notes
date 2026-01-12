FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client ./cmd/client/main.go
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/client /service-client
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/service-client"]

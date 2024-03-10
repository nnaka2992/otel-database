FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o start_server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/start_server /start_server
CMD ["/start_server"]

FROM cgr.dev/chainguard/go:latest-dev AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o libreria-mariela-api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/libreria-mariela-api .
COPY --from=builder /app/assets ./assets
ENV TZ=America/Argentina/Buenos_Aires
EXPOSE 8080
USER nonroot:nonroot
CMD ["./libreria-mariela-api"]
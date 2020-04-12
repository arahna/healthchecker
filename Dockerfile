FROM golang:1.14.2 as builder
LABEL maintainer="Julia N."
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

######## Start a new stage #######
FROM alpine:3.11.5
RUN adduser -D myuser
USER myuser

WORKDIR /app/
COPY --from=builder /app/main .

EXPOSE 8000
CMD ["./main"]
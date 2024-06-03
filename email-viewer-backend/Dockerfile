FROM golang:1.22.3 as builder

WORKDIR /app

COPY src/ ./
COPY go.mod ./
COPY go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o email-viewer-backend .

FROM debian:buster-slim
WORKDIR /bin  

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/email-viewer-backend .
EXPOSE 3000

# Command to run the executable
CMD ["./email-viewer-backend"]

FROM golang:1.22.3 as builder

WORKDIR /app

COPY src/ ./
COPY go.mod ./
COPY go.sum ./
#compile app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o indexer

FROM debian:buster-slim
WORKDIR /bin

#make binary accesible
COPY --from=builder /app/indexer .


# Download and unpack the Enron dataset
# ADD http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz .
# RUN tar -xzf enron_mail_20110402.tgz && rm enron_mail_20110402.tgz

CMD ["./indexer"]



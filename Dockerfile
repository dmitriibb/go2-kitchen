FROM golang:1.19
LABEL authors="dmitrii"

WORKDIR /app
COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /go2-kitchen

EXPOSE 9020
EXPOSE 9120

CMD ["/go2-kitchen"]

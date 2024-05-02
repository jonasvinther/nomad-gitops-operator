FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

ENV GO111MODULE=on

WORKDIR /app

ADD . .

RUN go mod download
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/nomoporator

RUN adduser -S scratchuser
RUN chown scratchuser /go/bin/nomoporator

FROM scratch
COPY --from=builder /go/bin/nomoporator /nomoporator
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser
ENTRYPOINT ["/nomoporator"]

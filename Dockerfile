# build stage
FROM golang AS builder

ENV GO111MODULE=on

WORKDIR /app/

COPY poller/go.mod .
COPY poller/go.sum .

RUN go mod download

COPY poller/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o poller

# final stage
FROM scratch
COPY --from=builder app/poller /app/
ENTRYPOINT ["/app/poller"]
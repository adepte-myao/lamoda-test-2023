FROM golang:1.21-alpine AS builder
WORKDIR /lamoda-test-2023
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
WORKDIR /lamoda-test-2023/cmd/app
RUN CGO_ENABLED=0 go build -o /bin/app
WORKDIR /lamoda-test-2023

FROM alpine:latest
RUN apk --update add ca-certificates
COPY --from=builder /bin/app /bin/app
COPY --from=builder /lamoda-test-2023/deploy/docker-compose/.env /deploy/docker-compose/.env
COPY --from=builder /lamoda-test-2023/configs/default.toml /configs/default.toml

CMD [ "/bin/app" ]
FROM golang:1.24.5-alpine AS build-stage

WORKDIR /home/app

COPY . .

RUN go build -o rate-limiter ./cmd/server/

FROM alpine:latest AS run-stage

COPY --from=build-stage /home/app /home/app

WORKDIR /home/app

EXPOSE 9240

CMD ["./rate-limiter"]
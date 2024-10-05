FROM golang:1.22.1 AS build

WORKDIR /app

COPY . .

RUN go build -o moodle_attendo cmd/moodle_attendo/main.go

FROM alpine:3.20.2

RUN apk add --no-cache tzdata

ENV TZ=Asia/Jakarta

RUN apk add --no-cache libc6-compat && apk add --no-cache chromium

WORKDIR /app

COPY --from=build /app/moodle_attendo /app/moodle_attendo

ENTRYPOINT ["/app/moodle_attendo"]

CMD []

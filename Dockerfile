FROM golang:alpine3.16 AS build
LABEL stage=build
WORKDIR /app
COPY . ./

RUN apk add build-base
RUN go build cmd/main.go

FROM alpine:3.16 AS runner
WORKDIR /app
LABEL authors="@Subudei"
COPY --from=build /app/main /app/main
COPY config.json /app/config.json
CMD ["/app/main"]
FROM golang:1.25 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:3.22.1 AS run
COPY --from=build /app/server /server
EXPOSE 8081
CMD ["/server"]

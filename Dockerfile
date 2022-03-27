FROM golang:1.18-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /frya-temp

FROM alpine
WORKDIR /
COPY --from=build /frya-temp /frya-temp
EXPOSE 3000
#USER nonroot:nonroot
ENTRYPOINT ["/frya-temp"]

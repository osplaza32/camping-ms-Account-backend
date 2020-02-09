
#build stage
FROM golang:alpine
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
CMD ["make all"]

RUN go mod download
COPY . .
RUN go build -o app  cmd/main.go
EXPOSE 8080
EXPOSE 50051
CMD ["/app/app"]
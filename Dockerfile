# build stage
FROM golang AS build-env
WORKDIR /go/src/github.com/ichsanrp/tax-calculator
COPY . /go/src/github.com/ichsanrp/tax-calculator
RUN go get -u
# Testing Apps
RUN go get -d -v gopkg.in/DATA-DOG/go-sqlmock.v1
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# final stage
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build-env /go/src/github.com/ichsanrp/tax-calculator/ .
CMD ["./app"]
EXPOSE 8080
FROM golang:1.14

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["DBTest]
FROM golang:1.8

WORKDIR $GOPATH/src/vue/src
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN ls -al

CMD ["src"]
FROM golang
ADD . /go/src/github.com/projectweekend/cta-bus-predictions
WORKDIR /go/src/github.com/projectweekend/cta-bus-predictions
RUN go get && go build

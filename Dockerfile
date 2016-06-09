FROM golang
 
ADD . /go/src/go-gateway
RUN go install go-gateway
ENTRYPOINT /go/bin/go-gateway
 
EXPOSE 80

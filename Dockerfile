FROM golang:1.17
WORKDIR /tasmota-http-go/
ADD . /tasmota-http-go/
RUN go build .

FROM scratch  
WORKDIR /
COPY --from=0 /tasmota-http-go/tasmota-http-go ./
CMD ["./tasmota-http-go"]  

FROM golang:1.12.1-alpine3.9 as builder
ARG LD_FLAGS

COPY ./ /go/src/github.com/cnative/servicebuilder/
WORKDIR /go/src/github.com/cnative/servicebuilder/

RUN apk add git make
RUN make install-deptools clean build 

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/cnative/servicebuilder/bin/servicebuilder /usr/bin
ENTRYPOINT ["/usr/bin/servicebuilder"]
CMD ["-h"]

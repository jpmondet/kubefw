FROM golang:1.11.2-alpine3.8 as build
#RUN adduser -D -g '' fwuser # Need to be privileged for iptables & co...
WORKDIR /go/src/github.com/jpmondet/kubefw
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/kubefw cmd/kubefw/kubefw.go
FROM alpine:3.9.3 as run-image
#COPY --from=build /etc/passwd /etc/passwd # Same here, need to be privileged 
COPY --from=build /go/bin/kubefw /
RUN apk add --no-cache \
      iptables \
      ip6tables \
      ipset \
      iproute2 \
      ipvsadm \
      conntrack-tools \
      curl \
      bash
CMD ["/kubefw"]

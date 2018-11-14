FROM golang:1.11.2-alpine3.8 as build
#RUN adduser -D -g '' fwuser # Need to be privileged for iptables & co...
WORKDIR /go/src/github.com/jpmondet/kubefw
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/kubefw cmd/kubefw/kubefw.go
FROM gcr.io/distroless/base
#COPY --from=build /etc/passwd /etc/passwd # Same here, need to be privileged 
COPY --from=build /go/bin/kubefw /
USER fwuser
CMD ["/kubefw"]

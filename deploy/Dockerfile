FROM alpine:latest
RUN apk --no-cache add ca-certificates

#RUN git clone git@github.com:thisrepo/thisrepo.git
#RUN go build thisrepo/cmd/node/main.go -o nodeserver

ADD nodeserver /
ADD /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
CMD ["/nodeserver"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates

#RUN git clone git@github.com:uxff/flexdrive.git
#RUN cd flexdrive
#RUN go build -o flexdrive cmd/node/main.go

ADD flexdrive /
ADD pkg/app/admin/view /pkg/app/admin/view
ADD pkg/app/customer/view /pkg/app/customer/view
ADD static /

ADD /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

CMD ["/flexdrive"]

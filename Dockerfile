# build
FROM index.alauda.cn/alaudaorg/gobuild:1.8-alpine as builder

WORKDIR $GOPATH/src/monkey
ENV MONKEY=$GOPATH/src/monkey
ENV COMPONENT="monkey"

COPY . $GOPATH/src/monkey
RUN make build

# prod
# base image contains migrate command
FROM index.alauda.cn/alaudaorg/alaudabase:alpine-supervisor-migrate-1

WORKDIR /
EXPOSE 80
ENV COMPONENT="monkey"
CMD ["monkey.sh"]

COPY --from=builder /go/src/monkey/*.sh /
COPY --from=builder /go/src/monkey/migrations /monkey/migrations
COPY --from=builder /go/src/monkey/conf/supervisord.conf /etc/supervisord.conf
COPY --from=builder /go/src/monkey/kill_supervisor.py /monkey/

RUN chmod +x /*.sh

COPY --from=builder /go/src/monkey/bin/monkey /monkey/monkey
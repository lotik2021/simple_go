FROM alpine:3.8
ENV TZ=Europe/Moscow
RUN apk --update add ca-certificates tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone
WORKDIR /app
COPY configs /app/configs
COPY docs /app/docs
COPY migrations /app/migrations
COPY maasapi /app
EXPOSE 8081
ENTRYPOINT ["./maasapi"]

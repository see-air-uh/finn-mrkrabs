FROM alpine:latest 

RUN mkdir /app 

COPY mrkrabs /app

CMD ["/app/mrkrabs"]
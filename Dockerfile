FROM alpine:3.8

RUN apk add --no-cache ca-certificates curl

ADD ./shutdown-deferrer /shutdown-deferrer
ADD ./pre-shutdown-hook /pre-shutdown-hook

ENTRYPOINT ["/shutdown-deferrer"]

CMD ["daemon"]

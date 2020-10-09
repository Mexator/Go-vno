# Build context should be project root
FROM alpine
COPY build/client /
RUN apk add --update fuse
ENTRYPOINT ["sh", "-c", "while :; do sleep 1; done"]

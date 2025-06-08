FROM alpine:3.12

RUN apk --no-cache add ca-certificates openssl && update-ca-certificates

WORKDIR /opt/sen

COPY sen-graphql-go ./bin/sen-graphql-go
COPY ./etc ./etc

ENTRYPOINT ["/opt/sen/bin/sen-graphql-go"]

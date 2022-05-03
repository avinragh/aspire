FROM golang:1.18.1-alpine3.15
RUN apk --update add curl
RUN apk add --no-cache git
ENV CGO_ENABLED=0
COPY go.mod go.sum /go/src/aspire/
WORKDIR /go/src/aspire
RUN ls
RUN go mod download
COPY . /go/src/aspire
RUN curl -o /usr/local/bin/swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64
RUN chmod +x /usr/local/bin/swagger
RUN /usr/local/bin/swagger generate model --spec=swagger.yml
RUN GOOS=linux go build -a -installsuffix cgo -o build/aspire
RUN chmod +x entrypoint.sh
EXPOSE 8080 8080
ENTRYPOINT ["sh","entrypoint.sh"]


# FROM alpine
# RUN apk add --no-cache ca-certificates && update-ca-certificates
# COPY --from=builder /go/src/aspire/build/aspire /aspire/aspire
# COPY --from=builder /go/src/aspire/entrypoint.sh /aspire/entrypoint.sh
# WORKDIR /aspire
# RUN chmod +x entrypoint.sh
# EXPOSE 8080 8080
# ENTRYPOINT ["sh","entrypoint.sh"]

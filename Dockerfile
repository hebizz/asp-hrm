#FROM golang:alpine as build
#RUN apk --no-cache add ca-certificates
#
#
#FROM scratch
#
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#EXPOSE 9006
#WORKDIR /app
#COPY asp-hrm /app/
#COPY config.yaml /app/
#ENTRYPOINT ["./asp-hrm"]

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 9006
WORKDIR /app
COPY asp-hrm /app/
ENTRYPOINT ["./asp-hrm"]

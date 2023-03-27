FROM golang:1.20-alpine as builder
WORKDIR /

COPY . ./
RUN go mod download


RUN go build -o /service-ticket

FROM alpine
COPY --from=builder /service-ticket .

EXPOSE 80
CMD [ "/service-ticket" ]
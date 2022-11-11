FROM golang:1.19

WORKDIR /app

COPY main.go ./

RUN go build -o secret main.go

EXPOSE 8080

CMD [ "./secret" ]

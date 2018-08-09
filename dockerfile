FROM golang:latest

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o cmd .

CMD ["/app/cmd"]
#CMD ["/app/cmd", "-set=CF1,1.33", "-items=CF1,CF1"]

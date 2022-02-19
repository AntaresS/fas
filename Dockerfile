FROM golang:1.16-alpine

WORKDIR /fas
COPY . ./
RUN go mod download
RUN go build -o ./fas

EXPOSE 9527

CMD [ "./fas" ]
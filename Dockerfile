FROM golang:alpine
ADD . ./
RUN go build -o goscrutinio
ENTRYPOINT [ "./goscrutinio" ]

FROM golang:1.18
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o workerbee .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/workerbee ./
CMD ["./workerbee"]
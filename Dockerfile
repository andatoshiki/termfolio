FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ssh-portfolio .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates openssh-keygen

WORKDIR /app
COPY --from=build /out/ssh-portfolio /app/ssh-portfolio
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 2222
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/ssh-portfolio"]

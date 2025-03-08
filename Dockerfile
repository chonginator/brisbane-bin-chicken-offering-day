FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY templates templates

COPY static static

COPY brisbane-bin-chicken-day .

EXPOSE 8080

CMD ["./brisbane-bin-chicken-day"]
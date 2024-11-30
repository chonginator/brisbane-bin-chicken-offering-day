FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY templates templates

COPY brisbane-bin-chicken-day .

CMD ["./brisbane-bin-chicken-day"]
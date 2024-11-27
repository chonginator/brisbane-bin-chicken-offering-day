FROM debian:stable-slim
WORKDIR /app
COPY templates templates
COPY brisbane-bin-chicken-day .
CMD ["./brisbane-bin-chicken-day"]
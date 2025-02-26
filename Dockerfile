FROM golang:1.23.2

ENV MTGJSON_API_CONFIG_LOCATION "/root/.config/mtgjson-api/config.json"

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/mtgjson-api .

RUN mkdir -p /var/log/mtgjson-api && chown -R 1000:1000 /var/log/mtgjson-api

CMD ["mtgjson-api", "run", "--config", "$MTGJSON_API_CONFIG_LOCATION"]

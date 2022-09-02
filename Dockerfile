FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod download -x && go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN make build && \
  curl -Lo wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
  chmod +x ./wait-for-it.sh

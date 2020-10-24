FROM golang:latest

COPY . .

RUN go get github.com/julienschmidt/httprouter \
    && GO111MOD=on github.com/aws/aws-sdk-go github.com/jackc/pgx/v4 github.com/jmoiron/sqlx  github.com/spf13/viper 

CMD go run .
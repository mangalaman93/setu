## setu

Find slots available in your center for vaccination

## How to Run

### Setup Config

```
cp .env.sample .env
# add values to `.env` file
source .env
```

```
go get github.com/mangalaman93/setu
cd $GOPATH/src/github.com/mangalaman93/setu
go build -o setu
./setu
```

This will run a background process and poll once every 5 min for empty slots.
The logs can be found at `$HOME/setu/logs/` and `pid` can be found in the file
`$HOME/setu/setu.pid`.

If a Sendgrid API key and other details are setup, it will email you when
empty slots are available for vaccination for 18+.

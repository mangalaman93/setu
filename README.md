## setu

Find slots available in your center for vaccination

## How to Run

### Setup Config

```
cp .env.sample .env
# add values to `.env` file
source .env
```

### Run the Daemon

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

## Get District and Center Information

* https://github.com/bhattbhavesh91/cowin-vaccination-slot-availability/blob/main/district_mapping%20v1.csv
* https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id={}&date={16-05-2021}

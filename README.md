# Run all tests

go test ./...

go install bitbucket.org/rwirdemann/kontrol

# Query accounts
```
curl -s http://localhost:8991/kontrol/accounts | python -m json.tool
curl -s http://localhost:8991/kontrol/accounts/AN/bookings | python -m json.tool
```

## Build for different Linux
```
env GOOS=linux GOARCH=amd64 go build bitbucket.org/rwirdemann/kontrol
scp kontrol root@94.130.79.196:~
```

# Run all tests

go test ./...

go install bitbucket.org/rwirdemann/kontrol

# Query account
```
curl -s http://localhost:8991/kontrol/accounts/AN/bookings | py -m json.tool
```

## Build for different Linux
```
env GOOS=linux GOARCH=amd64 go build bitbucket.org/rwirdemann/kontrol
scp kontrol root@94.130.79.196:~
```

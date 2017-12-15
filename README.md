# Kontrol Backend

## Todos
- port und filename über flags
- Monats Report
- GET /accounts soll keine bookings liefern

# Setup instructions

cd $PROJECTROOT

## Run main

go run kontrol/main.go

## Run all tests

go test ./...

## Build and install

go install bitbucket.org/rwirdemann/kontrol

## Regenerate HTML assets

Only necessary after html/index.html was changed

```
go-bindata -pkg html -o html/assets.go html/
```

## Build for different Linux
```
env GOOS=linux GOARCH=amd64 go build bitbucket.org/rwirdemann/kontrol
scp kontrol root@94.130.79.196:~
```

# Query accounts
```
curl -s http://localhost:8991/kontrol/accounts | python -m json.tool
curl -s http://localhost:8991/kontrol/accounts/AN | python -m json.tool
curl -s http://localhost:8991/kontrol/accounts/AN?year=2107&month=12 | python -m json.tool
```

# Rules

R1: AR = Ausgangsrechnungen

R1#S1Partner: Leistung wurde von Partner erbracht
- Partner bekommt 70% seiner Nettoposition
- Kommitment bekommt 25% der Partnernettoposition
- Vertrieb bekommt 5% der Partnernettoposition

R1#S2#Extern: Leistung wurde von Extern erbracht
- Kommitment bekommt 95% der Extern-Nettoposition
- Vertrieb bekommt 5% der Partner-Nettoposition

R1#S3#Employee: Leistung wurde vom Angestellten erbracht
- Kommitment bekommt 95% der Extern-Nettoposition
- Vertrieb bekommt 5% der Partner-Nettoposition
- 100% der Nettoposition weden auf das Angestelltenkonto verbucht

R2: ER = Eingangsrechnung
- 100% des Nettobetrags werden gegen das Kommitment-Konto gebucht

R3: GV = Geschäftsführerentnahmen
- 100% der Entnahme werden gegen das Kommitment-Konto gebucht
# Kontrol Backend

## Todos
- plausi check
- figo api einbinden: anke.nehrenberg@kommitment.biz

## API

### GET http://localhost:8991/kontrol/version

Aktuelle Version.

### GET http://localhost:8991/kontrol/bankaccount

Das Bankkonto inkl. Buchungen.

```
{
    "Owner": {
        "Name": "Kommitment GmbH & Co. KG",
        "Type": "bank"
    },
    "Bookings": [
        {
            "Typ": "ER",
            "Amount": 830.29,
            "Text": "hauptsache.net, Büro- und Konfimiete",
            "Year": 2017,
            "Month": 1
        },
        ...
    ],
    "Saldo": 18281.85
}
```

Aktuelle Version.

### GET http://localhost:8991/kontrol/accounts

Alle "virtuellen" Konten.

```
{
    "Accounts": [
        {
            "Owner": {
                "Id": "AN",
                "Name": "Anke Nehrenberg",
                "Type": "partner"
            },
            "Saldo": 0
        },
        ...
    ]
}
```

### GET http://localhost:8991/kontrol/accounts/AN

Ein einzelnes Konto inkl. Buchungen.

```
{
    "Owner": {
        "Id": "RW",
        "Name": "Ralf Wirdemann",
        "Type": "partner"
    },
    "Bookings": [
        {
            "Typ": "Nettoanteil",
            "Amount": 7559.999999999999,
            "Text": "RN_20170131-picue#NetShare#RW",
            "Year": 2017,
            "Month": 1
        },
        ...
    ],
    "Saldo": 18281.85
}
```

### GET http://localhost:8991/kontrol/accounts/collective

Die Buchungen des Bankkontos.

## CLI

```
cli --account=RW
```

## Run, Test, Build and Deploy

```
make run

Startet lokalen Server auf Port 8891.
```

```
make test

Führt alle Tests aus.
```

```
make build

Erzeugt das Binary kontrol-main im aktuellen Verzeichnis.
```

```
make install

Erzeugt das Binary kontrol und cli in $GOPATH/bin.
```

```
./deploy.sh 

Erzeugt Binary, Deployment auf 94.130.79.196, Neustart des Backend.
```
    
## Rules

### R1: AR = Ausgangsrechnungen

#### R1#S1Partner: Leistung wurde von Partner erbracht
- Partner bekommt 70% seiner Nettoposition
- Kommitment bekommt 25% der Partnernettoposition
- Vertrieb bekommt 5% der Partnernettoposition

#### R1#S2#Extern: Leistung wurde von Extern erbracht
- Kommitment bekommt 95% der Extern-Nettoposition
- Vertrieb bekommt 5% der Partner-Nettoposition

#### R1#S3#Employee: Leistung wurde vom Angestellten erbracht
- Kommitment bekommt 95% der Extern-Nettoposition
- Vertrieb bekommt 5% der Partner-Nettoposition
- 100% der Nettoposition weden auf das Angestelltenkonto verbucht

#### R2: ER = Eingangsrechnung
- 100% des Nettobetrags werden gegen das Kommitment-Konto gebucht

### R3: GV = Geschäftsführerentnahmen
- 100% der Entnahme werden gegen das Partner-Konto gebucht

### R4: IS = Interne Stunden
- 100% werden auf das Partner-Konto gebucht
- 100% werden gegen das Kommitment-Konto gebucht

### Offen R5: GWSteuer = Gewerbesteuer

### R6: SV-Beitrag
- 100% werden gegen das Kommitment-Konto gebucht

### Offen R7: LNSteuer

### Offen R8: Mitarbeiter Bonuszahlungen
# Kontrol Backend

## Todos

* Deployment automatisieren
* Excel Dartei regelmäßg abholen und parsen
* Vertriebsprovision für Angestellte auf deren Accounts buchen

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
            "CostCenter": "RW"
            "Amount": 830.29,
            "Text": "hauptsache.net, Büro- und Konfimiete",
            "Year": 2017,
            "Month": 1,
            "FileCreated": "2017-02-06T00:00:00Z"
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

Parameter: 
- cs="BW"   # Filter auf Costcenter (optional)

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
            "CostCenter": "JM"
            "Amount": 7559.999999999999,
            "Text": "RN_20170131-picue#NetShare#RW",
            "Year": 2017,
            "Month": 1,
            "FileCreated": "2017-02-06T00:00:00Z"
        },
        ...
    ],
    "Saldo": 18281.85
}
```

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
Die Regeln beschreiben, wie eine Buchung im Spreadsheet verbucht werden. Debei bedeutet R# eine Regel, S# eine Buchungsposition. "R1#S1Partner" ist also Regel 1, Buchungsposition S1 für Partner. "Vertrieb" ebntspricht der Spalte "Cost Center" im Spreadsheet.

Alle Beträge im Buchungssheet sind vorzeichenklos.

gegen: Betrag * -1
auf  : Betrag * 1

### R1: AR = Ausgangsrechnungen
- alle Ausgangsrechnungen werden netto auf das Bankkonto gebucht

#### R1#S1Partner: Leistung wurde von Partner erbracht
- Partner bekommt 70% seiner Nettoposition
- Kommitment bekommt 25% der Partnernettoposition
- Vertrieb bekommt 5% der Partnernettoposition

#### R1#S2#Extern: Leistung wurde von Extern erbracht
- Vertrieb bekommt 5% der Nettoposition
- Kommitment bekommt 95% des Nettorechnungsbetrags

#### R1#S3#Employee: Leistung wurde vom Angestellten erbracht
- Vertrieb bekommt 5% der Nettoposition
- Kommitment bekommt den 95% der Nettoposition

#### R2: ER = Eingangsrechnung
- 100% werden netto gegen das Bankkonto gebucht
- 100% des Nettobetrags werden gegen das Kommitment-Konto gebucht

### R3: GV = Partnerentnahme
- 100% werden gegen das Bankkonto gebucht
- 100% der Entnahme werden gegen das Partner-Konto gebucht

### R4: IS = Interne Stunden
- werden nicht auf das Bankkonto gebucht
- 100% werden auf das Partner-Konto gebucht
- 100% werden gegen das Kommitment-Konto gebucht

### R5: GWSteuer = Gewerbesteuer
- 100% werden auf das Bankkonto gebucht
- 100% werden gegen das Kommitment-Konto gebucht. Diese Regel ist noch unscharf:
  eigentlich müssen die 100% aufgeteilt werden auf: 70% auf Partner, 25% auf 
  Kommitment und 5% auf Dealbringer

### R6: SV-Beitrag
- 100% werden gegen das Bankkonto gebucht
- 100% werden gegen das Kommitment-Konto gebucht
- Kostenstelle: Kürzel des Angestellten

### Offen R7: LNSteuer

### R8: Gehalt Angestellter
- 100% Brutto gegen Bankkonto
- 100% Brutto gegen Kommitmentkonto
- Kostenstelle: Kürzel des Angestellten

### Offen R8: Mitarbeiter Bonuszahlungen

# Client
```
cli -h
Usage of cli:
  -account string
    	fetches given account
  -bank
    	fetches bank account
  -vsaldo
    	saldo sum virtual accounts
```

# Project todolist

Mein Todo-Listen Projekt für meine Bewerbung. Von Daniel Baliakin

## Start

Um diese Anwendung zu starten braucht man:
- Eine postgreSQL Instanz
- Eine .env Datei in der root directory mit folgendem Inhalt:
            
            PORT= // dein Port
            APP_ENV=local
            DB_HOST= // der Host auf welcher die PostgreSQL Datenbank läuft
            DB_PORT= // der Port auf welchen die PostgreSQL Datenbank hört
            DB_DATABASE= // der Name der Datenbank in PostgreSQL
            DB_USERNAME= // dein Benutzername in PostgreSQL
            DB_PASSWORD= // dein Passwort in PostgreSQL
            JWT_SECRET= // einen Geheimschlüssel zur Generierung von JSON Web Tokens (er sollte lang genug sein (> 256 bit))

Starten der Anwendung im Terminal in der root directory
"""bash
go run cmd/api/main.go
"""

Jetzt nur noch den Endpoint .../login im Browser abfragen und schon kann man sich registrieren!

## Features

- Accounts registrieren und anmelden
- Kategorien hinzufügen oder löschen
- Todos hinzufügen oder löschen
- Todos verschieben, sowohl untereinander als auch zwischen Kategorien
- Kategorien verschieben
- Todos als erledigt markieren

## Entstehung des Projekts

Benutzte Tools: 

- go-blueprint (https://github.com/Melkeydev/go-blueprint):
    Erstellung des Boilerplate Codes für einen Go Webserver mit Gin Gonic und PostgreSQL Unterstützung

- ChatGPT
    Erstellung der html und css Dateien sowie ein paar Funktionen in JavaScript

- golang-jwt (github.com/golang-jwt/jwt/v5):
    Authentifiezierung mittels JSON Web Tokens


## Retrospekte Verbesserungvorschläge

- Repository Pattern benutzen. Meine database Datei ist schon ziemlich vollgepackt geworden.
- Möglicherweise eine andere Methode zur Interaktion mit der Datenbank.
  Die SQL Queries waren relativ fehleranfällig während der Entwicklung und wahrscheinlich sind sie auch nicht wirklich sicher gegen
  SQL Injections, zum Beispiel.
- Mehr mit Context arbeiten. Scheint eine gute Praktik zu sein da man zum Beispiel das handeln eines abgebrochenen Requests ebenfalls abbrechen kann oder 
  Interprozesskommunkationen nutzen kann
- Tests schreiben. Damit habe ich leider wenig Erfahrung und ich bin leider auch nicht mehr zeitlich dazu gekommen. :(


## Anmerkungen

- Die air, docker-compose und Makefile Dateien waren bereits im Boilerplate mit dabei. Allerdings habe ich diese nicht benutzt und kein Funktional vorgesehen.
- Der JWT Token lebt eine Stunde. Es können Fehler auftreten sobald dieser abgelaufen ist. Ein erneutes Einloggen über /login löst aber das Problem!
- Ich habe davor noch nie mit JavaScript oder Go programmiert.
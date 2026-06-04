# Albumer

Prywatny katalog kolekcji CD, winyli, kaset i ksiazek. Aplikacja jest pisana jako jedna binarka Go z widokami SSR, malym wbudowanym JS do filtrowania oraz baza Neo4j jako osobna usluga.

Frontend jest renderowany przez backend Go, ale trzymany w normalnych plikach HTML/CSS/JS:

- `internal/web/templates/*.gohtml`
- `internal/web/assets/app.css`
- `internal/web/assets/app.js`

Szablony uzywaja rozszerzenia `.gohtml`, zeby GoLand rozpoznawal Go HTML templates i podpowiadal dane widoku. Pliki sa osadzane w binarce przez `go:embed`, wiec deploy dalej wymaga tylko aplikacji Go. Stary statyczny prototyp (`index.html`, `css/`, `js/`, `svg/`, PWA/service worker) zostal usuniety.

## Wymagania

- Go 1.22+
- Neo4j 5+

## Konfiguracja

Aplikacja czyta konfiguracje ze zmiennych srodowiskowych:

```sh
ADDR=:4000
NEO4J_URI=neo4j://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=neo4j
APP_SECRET=change-me
UPLOAD_DIR=./uploads
ADMIN_PASSWORD_HASH='bcrypt-hash'
```

Przy starcie aplikacja laduje tez lokalny plik `.env`, jesli istnieje. Zmienne ustawione juz w systemie maja pierwszenstwo przed wartosciami z `.env`.

Jesli `ADMIN_PASSWORD_HASH` jest puste, lokalne haslo developerskie to `admin`. Na serwerze ustaw hash bcrypt.

## Uruchomienie lokalne

```sh
go mod tidy
go run ./cmd/albumer
```

Po starcie wejdz na `http://localhost:4000`.

## Model danych

- `Person` reprezentuje artyste albo autora.
- `Group` reprezentuje zespol i ma relacje `HAS_MEMBER` do osob.
- `MusicItem` reprezentuje fizyczna pozycje muzyczna. Pole `medium` jest wymagane i przyjmuje `CD`, `VINYL` albo `CASSETTE`.
- `Book` reprezentuje ksiazke.

Ten sam album na CD i winylu to dwie osobne pozycje `MusicItem`.
Pozycje tego samego albumu na roznych nosnikach mozna polaczyc relacja `SAME_RELEASE`. Formularz muzyki pozwala wybrac relacje recznie i proponuje istniejace pozycje po wpisaniu tytulu.

## Deploy jako systemd

Przykladowy unit jest w `deploy/albumer.service`. Przed startem ustaw prawdziwe sekrety i katalog uploadow.

```sh
go build -o albumer ./cmd/albumer
sudo install -m 755 albumer /opt/albumer/albumer
sudo install -m 644 deploy/albumer.service /etc/systemd/system/albumer.service
sudo systemctl daemon-reload
sudo systemctl enable --now albumer
```

## Testy

```sh
go test ./...
```

## Import starego JSON

Stary `json/data.json` mozna zaladowac do Neo4j jako kolekcje plyt CD:

```sh
go run ./cmd/import-json -file json/data.json
```

Importer tworzy osoby, zespoly (`Band` w Neo4j), relacje czlonkow zespolow oraz `MusicItem` z `medium=CD`. Rekordy w bazie dostaja nowe UUID; stare `id` z JSON sa uzywane tylko tymczasowo podczas importu do zbudowania relacji.

Jesli baza byla wypelniona przed dodaniem przyjaznego wyszukiwania, odpal jednorazowo:

```sh
go run ./cmd/backfill-search
```

To uzupelni pola `searchName` i `searchTitle`, dzieki czemu np. `ostr` znajdzie `O.S.T.R`.

Stare losowe ID w juz istniejacej bazie mozna zamienic na UUID:

```sh
go run ./cmd/migrate-ids
```

Skrypt zapisuje poprzednie ID w `legacyId`, wiec relacje zostaja zachowane.

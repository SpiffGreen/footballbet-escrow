
# FootballBet-Escrow API

FootballBet-Escrow API is a simple API for showcasing an escrow service.

## Tech Stack

GO, Gin, Gorm, Postgresql

[Soccerfootball Info API - RapidAPI](https://rapidapi.com/soccerfootball-info-soccerfootball-info-default/api/soccer-football-info/playground/apiendpoint_c6020c0a-9773-499f-bc5d-81c49ed25ee2)


## Features

- Auth
- Get games to bet
- Place bets
- Deposit and Withdraw from wallet


## Run Locally

Clone the project

```bash
  git clone https://github.com/spiffgreen/footballbet-escrow
```

Go to the project directory

```bash
  cd footballbet-escrow
```

Run application

```bash
  go run main.go
```
## Database Setup
Using SQLite. To setup the db and tables, simple run the command below

```sh
go run migrate/migrate.go
```

## Documentation
[Documentation](./docs/Football-escrow_postman.json)

## Contributing

Contributions are always welcome!

See `CONTRIBUTING.md` for ways to get started.

Please adhere to this project's `code of conduct`.

## Author
[Spiff Jekey-Green](https://www.github.com/spiffgreen)

## License

[MIT](https://choosealicense.com/licenses/mit/)


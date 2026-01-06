module github.com/corecollectives/mist/cli

go 1.25.1

require (
	github.com/corecollectives/mist v0.0.0
	github.com/mattn/go-sqlite3 v1.14.32
	golang.org/x/crypto v0.43.0
	golang.org/x/term v0.36.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/corecollectives/mist => ../server

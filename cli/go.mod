module github.com/corecollectives/mist/cli

go 1.25.1

require (
	github.com/corecollectives/mist v0.0.0
	golang.org/x/crypto v0.43.0
	golang.org/x/term v0.36.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/corecollectives/mist => ../server

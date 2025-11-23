module github.com/corecollectives/mist

go 1.25.1

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/websocket v1.5.3
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/rs/zerolog v1.34.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	golang.org/x/crypto v0.43.0
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/docker/docker/api => github.com/moby/moby/api v1.52.0-beta.2

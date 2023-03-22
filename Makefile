linux:
	GOOS=linux  GOAMD64=v3 CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o fontface main.go

mac:
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o fontface-mac main.go

windows:
	GOOS=windows GOAMD64=v3 CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o fontface.exe main.go

release: linux mac windows


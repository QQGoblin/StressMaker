# StressMaker

```bash
CGO_ENABLED=1 go build -ldflags '-extldflags "-static"' -o stress cmd/main.go
```
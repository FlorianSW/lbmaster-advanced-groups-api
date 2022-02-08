mkdir -p build/
rm -rf build/*

GOOS=windows GOARCH=amd64 go build -o build/advanced_groups_api_win64.exe handler/cmd/cmd.go
GOOS=linux GOARCH=amd64 go build -o build/advanced_groups_api_linux handler/cmd/cmd.go

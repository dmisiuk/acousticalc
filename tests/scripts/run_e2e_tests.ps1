# PowerShell script for running E2E tests on Windows
choco install asciinema -y
$env:Path += ";C:\Program Files\Asciinema\bin"
go build -o cmd/acousticalc/acousticalc.exe ./cmd/acousticalc
go test -v -timeout 120s ./tests/e2e/...
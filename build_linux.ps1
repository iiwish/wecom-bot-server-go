# PowerShell 脚本：交叉编译为 Linux amd64 可执行文件
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o app ./cmd/main.go

Write-Host "已生成 Linux 可执行文件：app"
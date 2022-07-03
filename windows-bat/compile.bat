@echo off

IF ! command -v go (
  echo "golang could not be found; terminating installation"
  echo "install golang at https:/\go.dev\"
  exit
)

@echo on

echo "go was found; attempting to install with go"
go build -o fit
copy fit C:\Windows\System32
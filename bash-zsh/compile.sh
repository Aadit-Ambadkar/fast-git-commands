if ! command -v go &> /dev/null
then
    echo "golang could not be found; terminating installation"
    echo "install golang at https://go.dev/"
    exit
fi

echo "go was found; attempting to install with go"
go build -o fit
sudo cp fit /usr/local/bin
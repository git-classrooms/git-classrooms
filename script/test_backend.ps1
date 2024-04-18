$command=$args[0]

if (!$command) {
    echo "Usage: $PSCommandPath one|air"
    exit 1
}

if ($command -eq "one") {
    go mod download
    go generate
    go test -v ./...
}

if ($command -eq "air") {
    air -c .air.test.toml
}

# go1.20

function check() {
  golint -set_exit_status -min_confidence 0.3 ./...
  gocyclo -avg -over 15 .
  golangci-lint run ./...
  gosec -quiet ./...
}

# linux
export GOOS=linux
export GOARCH=amd64
check
export GOARCH=386
check
export GOARCH=arm64
check

# windows
export GOOS=windows
export GOARCH=amd64
check
export GOARCH=386
check
export GOARCH=arm64
check

# darwin
export GOOS=darwin
export GOARCH=amd64
check
export GOARCH=arm64
check
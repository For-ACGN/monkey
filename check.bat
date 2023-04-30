rem go1.20

rem windows
set GOOS=windows
set GOARCH=amd64
call :check
set GOARCH=386
call :check
set GOARCH=arm64
call :check

rem linux
set GOOS=linux
set GOARCH=amd64
call :check
set GOARCH=386
call :check
set GOARCH=arm64
call :check

rem darwin
set GOOS=darwin
set GOARCH=amd64
call :check
set GOARCH=arm64
call :check

:check
 golint -set_exit_status -min_confidence 0.3 ./...
 gocyclo -avg -over 15 .
 golangci-lint run ./...
 gosec -quiet ./...
:END

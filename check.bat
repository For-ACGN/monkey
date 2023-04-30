@echo off

rem go1.20

rem windows
set GOOS=windows
set GOARCH=amd64
call :check %1
set GOARCH=386
call :check %1
set GOARCH=arm64
call :check %1

rem linux
set GOOS=linux
set GOARCH=amd64
call :check %1
set GOARCH=386
call :check %1
set GOARCH=arm64
call :check %1

rem darwin
set GOOS=darwin
set GOARCH=amd64
call :check %1
set GOARCH=arm64
call :check %1

exit /b

:check
  echo ------------------------------------------------
  echo check %GOOS% %GOARCH%

  if "%1" == "-golint" (
    golint -set_exit_status -min_confidence 0.3 ./...
    goto END
  )
  if "%1" == "-gocyclo" (
    gocyclo -avg -over 15 .
    goto END
  )
  if "%1" == "-cilint" (
    golangci-lint run ./...
    goto END
  )
  if "%1" == "-gosec" (
    gosec -quiet ./...
    goto END
  )

  golint -set_exit_status -min_confidence 0.3 ./...
  gocyclo -avg -over 15 .
  golangci-lint run ./...
  gosec -quiet ./...
  goto END

:END
  echo ------------------------------------------------
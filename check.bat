@echo off

rem go1.20

set exit_code=0

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

echo exit code: %exit_code%
exit /b %exit_code%

:check
  echo ================================================
  echo check %GOOS% %GOARCH%
  echo ------------------------------------------------

  if "%1" == "-golint" (
    golint -set_exit_status -min_confidence 0.3 ./...
    call :echo_line
    call :set_exit_code
    goto END
  )
  if "%1" == "-gocyclo" (
    gocyclo -avg -over 15 .
    call :echo_line
    call :set_exit_code
    goto END
  )
  if "%1" == "-cilint" (
    golangci-lint run ./...
    call :echo_line
    call :set_exit_code
    goto END
  )
  if "%1" == "-gosec" (
    gosec -quiet ./...
    call :echo_line
    call :set_exit_code
    goto END
  )

  echo ------------------------------------------------
  echo golint
  echo ------------------------------------------------
  golint -set_exit_status -min_confidence 0.3 ./...
  call :set_exit_code
  echo ------------------------------------------------
  echo:

  echo ------------------------------------------------
  echo gocyclo
  echo ------------------------------------------------
  gocyclo -avg -over 15 .
  call :set_exit_code
  echo ------------------------------------------------
  echo:

  echo ------------------------------------------------
  echo golangci-lint
  echo ------------------------------------------------
  golangci-lint run ./...
  call :set_exit_code
  echo ------------------------------------------------
  echo:

  echo ------------------------------------------------
  echo gosec
  echo ------------------------------------------------
  gosec -quiet ./...
  call :set_exit_code
  echo ------------------------------------------------

  call :echo_line
  goto END
rem end

:echo_line
  echo ================================================
  echo:
  goto END
rem end

:set_exit_code
  if not %ERRORLEVEL% == 0 (
    set exit_code=1
  ) else (
    echo pass
  )
  goto END
rem end

:END

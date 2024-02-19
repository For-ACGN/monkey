@echo off

if "%1" == "-help" (
  echo -golint
  echo -gocyclo
  echo -cilint
  echo -gosec
  goto :EOF
)

if "%2" == "-e" (
  set exit_on_error=1
)

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
set GOARCH=loong64
call :check %1

rem darwin
set GOOS=darwin
set GOARCH=amd64
call :check %1
set GOARCH=arm64
call :check %1

:exit_bat
if %exit_code% == 0 (
  echo all check passed
) else (
  echo exit code: %exit_code%
)
exit /b %exit_code%
rem END_exit_bat

:check
  echo ================================================
  echo check %GOOS% %GOARCH%
  echo ------------------------------------------------

  if "%1" == "-golint" (
    golint -set_exit_status -min_confidence 0.3 ./...
    call :set_exit_code
    call :echo_line
    goto :EOF
  )
  if "%1" == "-gocyclo" (
    gocyclo -avg -over 15 .
    call :set_exit_code
    call :echo_line
    goto :EOF
  )
  if "%1" == "-cilint" (
    golangci-lint run ./...
    call :set_exit_code
    call :echo_line
    goto :EOF
  )
  if "%1" == "-gosec" (
    gosec -quiet ./...
    call :set_exit_code
    call :echo_line
    goto :EOF
  )

  echo ------------------------------------------------
  echo golint
  echo ------------------------------------------------
  golint -set_exit_status -min_confidence 0.3 ./...
  call :set_exit_code
  echo ------------------------------------------------
  echo.

  echo ------------------------------------------------
  echo gocyclo
  echo ------------------------------------------------
  gocyclo -avg -over 15 .
  call :set_exit_code
  echo ------------------------------------------------
  echo.

  echo ------------------------------------------------
  echo golangci-lint
  echo ------------------------------------------------
  golangci-lint run ./...
  call :set_exit_code
  echo ------------------------------------------------
  echo.

  echo ------------------------------------------------
  echo gosec
  echo ------------------------------------------------
  gosec -quiet ./...
  call :set_exit_code
  echo ------------------------------------------------

  call :echo_line
  goto :EOF
rem END_check

:set_exit_code
  if not %ERRORLEVEL% == 0 (
    set exit_code=1
    rem if %exit_on_error% == 1 (
    rem   goto :exit_bat
    rem )
  ) else (
    echo pass
  )
  goto :EOF
rem END_set_exit_code

:echo_line
  echo ================================================
  echo.
  goto :EOF
rem END_echo_line

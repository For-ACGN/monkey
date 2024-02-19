@echo off

rem initialize environment variables
set exit_code=0
set exit_on_error=0

rem print help information
if "%1" == "-help" (
  echo -golint
  echo -gocyclo
  echo -cilint
  echo -gosec
  goto :EOF
)

rem process arguments
if "%2" == "-e" (
  set exit_on_error=1
)

rem -----------------Windows-----------------
set GOOS=windows
rem -----------------------------
set GOARCH=amd64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=386
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=arm64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------


rem ------------------Linux------------------
set GOOS=linux
rem -----------------------------
set GOARCH=amd64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=386
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=arm64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=loong64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------


rem -----------------Darwin------------------
set GOOS=darwin
rem -----------------------------
set GOARCH=amd64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------
set GOARCH=arm64
call :check %1
rem --------exit on error--------
if not exit_code == 0 (
  if %exit_on_error% == 1 (
    goto :exit_bat
  )
)
rem -----------------------------


rem end of script
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
  if %ERRORLEVEL% == 0 (
    echo pass
  ) else (
    set exit_code=1
  )
  goto :EOF
rem END_set_exit_code

:echo_line
  echo ================================================
  echo.
  goto :EOF
rem END_echo_line

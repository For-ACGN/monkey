@echo off

rem initialize environment variables
set is_print_help=0
set exit_code=0
set exit_on_error=0

call :main %1 %2
exit /b %exit_code%

:main
  rem process print help information
  call :is_print_help %1
  if %is_print_help% == 1 (
    call :print_help
    goto :EOF
  )

  rem process arguments
  if "%2" == "-e" (
    set exit_on_error=1
  )


  if %exit_code% == 0 (
    echo all check passed
  ) else (
    echo exit code: %exit_code%
  )
goto :EOF
rem END_main

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

:is_print_help
  if "%1" == "-help" (
    set is_print_help=1
  )
  if "%1" == "--help" (
    set is_print_help=1
  )
  if "%1" == "-h" (
    set is_print_help=1
  )
  if "%1" == "/?" (
    set is_print_help=1
  )
goto :EOF
rem END_is_print_help

:print_help
  echo Usage of lint:
  echo   -golint      only use golint to check code
  echo   -gocyclo     only use gocyclo to check code
  echo   -cilint      only use golangci-lint to check code
  echo   -gosec       only use gosec to check code
  echo   -e           interrupt script when detect error
  echo.
  echo example:
  echo   "lint -golint"    only use golint to check code
  echo   "lint -gosec -e"  only use gosec and exit on error
goto :EOF
rem END_print_help

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

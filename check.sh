#!/bin/bash

# go1.20

function main() {
  export exit_code=0

  # linux
  export GOOS=linux
  export GOARCH=amd64
  check $1
  export GOARCH=386
  check $1
  export GOARCH=arm64
  check $1

  # windows
  export GOOS=windows
  export GOARCH=amd64
  check $1
  export GOARCH=386
  check $1
  export GOARCH=arm64
  check $1

  # darwin
  export GOOS=darwin
  export GOARCH=amd64
  check $1
  export GOARCH=arm64
  check $1

  return $exit_code
}

function check() {
  echo ================================================
  echo check $GOOS $GOARCH
  echo ------------------------------------------------

  case $1 in
  -golint)
    golint -set_exit_status -min_confidence 0.3 ./...
    set_exit_code
    echo_line
  ;;

  -gocyclo)
    gocyclo -avg -over 15 .
    set_exit_code
    echo_line
  ;;

  -cilint)
    golangci-lint run ./...
    set_exit_code
    echo_line
  ;;

  -gosec)
    gosec -quiet ./...
    set_exit_code
    echo_line
  ;;

  *)
    echo ------------------------------------------------
    echo golint
    echo ------------------------------------------------
    golint -set_exit_status -min_confidence 0.3 ./...
    set_exit_code
    echo ------------------------------------------------
    echo

    echo ------------------------------------------------
    echo gocyclo
    echo ------------------------------------------------
    gocyclo -avg -over 15 .
    set_exit_code
    echo ------------------------------------------------
    echo

    echo ------------------------------------------------
    echo golangci-lint
    echo ------------------------------------------------
    golangci-lint run ./...
    set_exit_code
    echo ------------------------------------------------
    echo

    echo ------------------------------------------------
    echo gosec
    echo ------------------------------------------------
    gosec -quiet ./...
    set_exit_code
    echo ------------------------------------------------

    echo_line
  ;;
  esac
}

function echo_line() {
  echo ================================================
  echo
}

function set_exit_code() {
  if [$? != 0]
  then
    export exit_code=1
  else
    echo pass
  fi
}

main $1

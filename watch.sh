#!/usr/bin/env bash
while true; do
  sleep .25
  #echo === Scanning... ===
  if [[  ( ! -f ./go-lox || -n "$(find . -name '*.go' -newer ./go-lox -print | head -n 1)" ) ]] ; then
    clear

    # This is a terminal control sequence, proprietary to iTerm,
    # that clears the screen, including scrollback history. See
    # https://www.iterm2.com/documentation-escape-codes.html
    echo "]1337;ClearScrollback"

    true &&
      # echo Reformatting...
      # find . -name '*.go' -newer ./go-lox -print -exec go fmt '{}' \;
      echo Generating ...   &&
        (go generate ./...       || (echo Generating FAILED && false)) &&
      echo Tooling...       &&
          (go build cmd/generate-ast/generate-ast.go &&
           ./generate-ast -d ast || (echo Tooling FAILED && false)) &&
      echo Reformatting ... &&
        (go fmt ./...            || (echo Reformatting FAILED && false)) &&
      echo Building ...     &&
        (go build                || (echo Building FAILED && false)) &&
      echo Testing ...      &&
        (go test ./...           ||
          (echo Testing FAILED && (go test ./... | cmd/test-diff.sh) && false)
        )                   &&
      echo SUCCESS

      echo '-----'

   touch ./go-lox
  fi
done

#!/bin/bash
set -eu

touch go.mod

PROJECT_NAME=$(basename $(pwd | xargs dirname))
CURRENT_DIR=$(basename $(pwd))

CONTENT=$(cat <<-EOD
module github.com/rahulzzore/realworld-aws-lambda-dynamodb-go

require (
	github.com/aws/aws-lambda-go v1.6.0
	github.com/aws/aws-sdk-go v1.23.15
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gosimple/slug v1.7.0
	github.com/rainycape/unidecode v0.0.0-20150907023854-cb7f23ec59be // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
)
EOD
)

echo "$CONTENT" > go.mod

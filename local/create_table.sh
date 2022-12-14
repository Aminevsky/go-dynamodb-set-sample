#!/bin/sh

aws dynamodb create-table \
    --table-name BaseballTeams \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:4566 \
    --profile localstack
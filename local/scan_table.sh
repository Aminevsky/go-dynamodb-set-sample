#!/bin/sh

aws dynamodb scan \
    --table-name BaseballTeams \
    --endpoint-url http://localhost:4566 \
    --profile localstack
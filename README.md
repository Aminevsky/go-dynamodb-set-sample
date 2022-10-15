# DynamoDB Set Sample for Golang

This repository shows how to handle DynamoDB's Set using AWS SDK for Go.

## Memo

For this is my personal memo, I don't write `main` function. If you want to run, `go test`.

## Marshal

By default, [feature/dynamodb/attributevalue](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue) marshals slices to a List. If you want to marshal slices to a Set, you must specify in struct tag.

Field `Reserve` of struct `BaseballTeam` shows how to do it.

## Add

In DynamoDB, to add elements to a set we use `ADD` action.

Method `BaseballTeamRepository.AddReserve()` shows how to do it.

## Delete

In DynamoDB, to delete elements from a set we use `DELETE` action.

Method `BaseballTeamRepository.DeleteReserve()` shows how to do it.

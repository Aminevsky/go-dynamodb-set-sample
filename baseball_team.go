package main

type BaseballTeam struct {
	ID           string `dynamodbav:"id"`
	TeamName     string `dynamodbav:"team_name"`
	BattingOrder []int  `dynamodbav:"batting_order"`
	Reserve      []int  `dynamodbav:"reserve,numberset"`
}

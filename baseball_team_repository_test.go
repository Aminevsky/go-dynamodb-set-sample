package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
)

func createTestDynamoDBClient() (*dynamodb.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: "ap-northeast-1",
			}, nil
		})),
		config.WithSharedConfigProfile("localstack"),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(sdkConfig), nil
}

func TestBaseballTeamRepository_Create(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"
	repo := NewBaseballTeamRepository(client, tableName)

	team := BaseballTeam{
		ID:           "test001",
		TeamName:     "Team 1",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	err = repo.Create(context.TODO(), team)
	if err != nil {
		t.Fatalf("item create failed: %s", err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("item get failed: %s", err)
	}

	var item BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &item); err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}

	if item.ID != team.ID {
		t.Errorf("ID should be %s, got =  %s", team.ID, item.ID)
	}
	if item.TeamName != team.TeamName {
		t.Errorf("TeamName should be %s, got = %s", team.TeamName, item.TeamName)
	}

	if !reflect.DeepEqual(item.BattingOrder, team.BattingOrder) {
		t.Errorf("BattingOrder is mismatch, got = %v, want = %v", item.BattingOrder, team.BattingOrder)
	}
	if !reflect.DeepEqual(item.Reserve, team.Reserve) {
		t.Errorf("Reserve is mismatch, got = %v, want = %v", item.Reserve, team.Reserve)
	}
}

func TestBaseballTeamRepository_Get_Success(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	team := BaseballTeam{
		ID:           "test002",
		TeamName:     "Team 2",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}
	item, _ := attributevalue.MarshalMap(team)

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	repo := NewBaseballTeamRepository(client, tableName)
	got, err := repo.Get(context.TODO(), "test002")
	if err != nil {
		t.Fatalf("Get failed: %s", err)
	}

	if got.ID != team.ID {
		t.Errorf("ID should be %s, got = %s", team.ID, got.ID)
	}
	if got.TeamName != team.TeamName {
		t.Errorf("TeamName should be %s, got = %s", team.TeamName, got.TeamName)
	}
	if !reflect.DeepEqual(got.BattingOrder, team.BattingOrder) {
		t.Errorf("BattingOrder is mismatch, got = %v, want = %v", got.BattingOrder, team.BattingOrder)
	}
	if !reflect.DeepEqual(got.Reserve, team.Reserve) {
		t.Errorf("Reserve is mismatch, got = %v, want = %v", got.Reserve, team.Reserve)
	}
}

func TestBaseballTeamRepository_Get_NoResponse(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: "invalid"}},
		TableName: aws.String(tableName),
	})

	repo := NewBaseballTeamRepository(client, tableName)
	got, err := repo.Get(context.TODO(), "invalid")
	if err != nil {
		t.Fatalf("Get failed: %s", err)
	}

	want := &BaseballTeam{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response is mismatch. got = %v, want = %v", got, want)
	}
}

func TestBaseballTeamRepository_AddReserve(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	team := BaseballTeam{
		ID:           "test003",
		TeamName:     "Team 3",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}
	item, _ := attributevalue.MarshalMap(team)

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	repo := NewBaseballTeamRepository(client, tableName)
	err = repo.AddReserve(context.TODO(), team.ID, []int{7, 8, 9})
	if err != nil {
		t.Fatalf("AddReserve failed: %s", err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("item get failed: %s", err)
	}

	var got BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &got); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	want := []int{4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(got.Reserve, want) {
		t.Errorf("Reserve is mismatch: got = %v, want = %v", got.Reserve, want)
	}
}

func TestBaseballTeamRepository_DeleteReserve(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	team := BaseballTeam{
		ID:           "test004",
		TeamName:     "Team 4",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}
	item, _ := attributevalue.MarshalMap(team)

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	repo := NewBaseballTeamRepository(client, tableName)
	if err := repo.DeleteReserve(context.TODO(), team.ID, []int{4}); err != nil {
		t.Fatalf("DeleteReserve failed: %s", err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("item get failed: %s", err)
	}

	var got BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &got); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	want := []int{5, 6}
	if !reflect.DeepEqual(got.Reserve, want) {
		t.Errorf("Reserve is mismatch: got = %v, want = %v", got.Reserve, want)
	}
}

func TestBaseballTeamRepository_AddBattingOrder(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	team := BaseballTeam{
		ID:           "test005",
		TeamName:     "Team 5",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}
	item, _ := attributevalue.MarshalMap(team)

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	repo := NewBaseballTeamRepository(client, tableName)
	if err := repo.AddBattingOrder(context.TODO(), team.ID, []int{10, 11}); err != nil {
		t.Fatalf("AddBattingOrder failed: %s", err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("item get failed: %s", err)
	}

	var got BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &got); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	want := []int{1, 2, 3, 10, 11}
	if !reflect.DeepEqual(got.BattingOrder, want) {
		t.Errorf("BattingOrder is mismatch: got = %v, want = %v", got.BattingOrder, want)
	}
}

func TestBaseballTeamRepository_RemoveBattingOrder(t *testing.T) {
	client, err := createTestDynamoDBClient()
	if err != nil {
		t.Fatalf("db client create failed: %s", err)
	}

	tableName := "BaseballTeams"

	team := BaseballTeam{
		ID:           "test006",
		TeamName:     "Team 6",
		BattingOrder: []int{1, 2, 3},
		Reserve:      []int{4, 5, 6},
	}
	item, _ := attributevalue.MarshalMap(team)

	client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	repo := NewBaseballTeamRepository(client, tableName)
	if err := repo.RemoveBattingOrder(context.TODO(), team.ID, []int{1, 2}); err != nil {
		t.Fatalf("RemoveBattingOrder failed: %s", err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: team.ID}},
		TableName: aws.String(tableName),
	})
	if err != nil {
		t.Fatalf("item get failed: %s", err)
	}

	var got BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &got); err != nil {
		t.Errorf("unmarshal error: %s", err)
	}

	want := []int{1}
	if !reflect.DeepEqual(got.BattingOrder, want) {
		t.Errorf("BattingOrder is mismatch: got = %v, want = %v", got.BattingOrder, want)
	}
}

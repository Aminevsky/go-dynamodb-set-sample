package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

type BaseballTeamRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewBaseballTeamRepository(client *dynamodb.Client, tableName string) *BaseballTeamRepository {
	return &BaseballTeamRepository{
		client:    client,
		tableName: tableName,
	}
}

// Update updates an item by the team parameter.
func (r *BaseballTeamRepository) Update(ctx context.Context, team BaseballTeam) error {
	key, err := attributevalue.Marshal(team.ID)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       map[string]types.AttributeValue{"ID": key},
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  nil,
		ExpressionAttributeValues: nil,
		UpdateExpression:          nil,
	})
	if err != nil {
		return err
	}

	return nil
}

// Create creates an item by the team parameter.
func (r *BaseballTeamRepository) Create(ctx context.Context, team BaseballTeam) error {
	item, err := attributevalue.MarshalMap(team)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return fmt.Errorf("create failed: %w", err)
	}

	return nil
}

// Get gets an item by the id parameter.
// Reserve of the response should be unmarshalled to a slice.
func (r *BaseballTeamRepository) Get(ctx context.Context, id string) (*BaseballTeam, error) {
	key, err := attributevalue.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	res, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("get failed: %w", err)
	}

	var team BaseballTeam
	if err := attributevalue.UnmarshalMap(res.Item, &team); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &team, nil
}

// AddReserve adds addNumbers to a set "reserve" of the item specified by the id parameter.
func (r *BaseballTeamRepository) AddReserve(ctx context.Context, id string, addNumbers []int) error {
	key, err := attributevalue.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// each element must be string
	params := make([]string, 0, len(addNumbers))
	for _, n := range addNumbers {
		params = append(params, strconv.Itoa(n))
	}

	// value is number set, not number list
	update := expression.Add(expression.Name("reserve"), expression.Value(types.AttributeValueMemberNS{Value: params}))
	builder, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
		UpdateExpression:          builder.Update(),
	})
	if err != nil {
		return fmt.Errorf("add to reserve failed: %w", err)
	}

	return nil
}

// DeleteReserve deletes deleteNumbers from a set "reserve" of the item specified by the id parameter.
func (r *BaseballTeamRepository) DeleteReserve(ctx context.Context, id string, deleteNumbers []int) error {
	key, err := attributevalue.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// each element must be string
	params := make([]string, 0, len(deleteNumbers))
	for _, n := range deleteNumbers {
		params = append(params, strconv.Itoa(n))
	}

	// value is number set, not number list
	update := expression.Delete(expression.Name("reserve"), expression.Value(types.AttributeValueMemberNS{Value: params}))
	builder, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
		UpdateExpression:          builder.Update(),
	})
	if err != nil {
		return fmt.Errorf("delete from reserve failed: %s", err)
	}

	return nil
}

// AddBattingOrder adds addNumbers to a list "batting_order" of the item specified by the id parameter.
func (r *BaseballTeamRepository) AddBattingOrder(ctx context.Context, id string, addNumbers []int) error {
	key, err := attributevalue.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	update := expression.Set(
		expression.Name("batting_order"),
		expression.ListAppend(expression.Name("batting_order"), expression.Value(addNumbers)),
	)
	builder, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
		UpdateExpression:          builder.Update(),
	})
	if err != nil {
		return fmt.Errorf("add to batting_order failed: %s", err)
	}

	return nil
}

// RemoveBattingOrder removes removeNumbers from a list "batting_order" of the item specified by the id parameter.
func (r *BaseballTeamRepository) RemoveBattingOrder(ctx context.Context, id string, removeNumbers []int) error {
	key, err := attributevalue.MarshalMap(map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	fmt.Printf("[debug] removeNumbers=%v\n", removeNumbers)

	var update expression.UpdateBuilder
	for _, n := range removeNumbers {
		update = update.Remove(expression.Name("batting_order[" + strconv.Itoa(n) + "]"))
	}

	//update.Remove(expression.Name("batting_order" + strconv.Itoa(n)))

	builder, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("builder error: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  builder.Names(),
		ExpressionAttributeValues: builder.Values(),
		UpdateExpression:          builder.Update(),
	})
	if err != nil {
		return fmt.Errorf("remove from batting_order failed: %s", err)
	}

	return nil
}

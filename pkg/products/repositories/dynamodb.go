package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/products"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ProductRepository struct {
	DynamoDB  *dynamodb.DynamoDB
	tableName *string
}

func (repository *ProductRepository) Store(product *products.Product) error {
	item, err := dynamodbattribute.MarshalMap(product)

	if err != nil {
		return err
	}

	_, err = repository.DynamoDB.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(id)"),
		Item:                item,
		TableName:           repository.tableName,
	})

	return err
}

func (repository *ProductRepository) Update(product *products.Product) error {
	item, err := dynamodbattribute.MarshalMap(product)

	if err != nil {
		return err
	}

	_, err = repository.DynamoDB.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_exists(id)"),
		Item:                item,
		TableName:           repository.tableName,
	})

	return err
}
func (repository *ProductRepository) FindOne(ID *string) (*products.Product, error) {
	var item *products.Product
	output, err := repository.DynamoDB.GetItem(&dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: ID}},
		TableName: repository.tableName,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, item)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (repository *ProductRepository) batchRequest(key string, items []*string) ([]*products.Product, error) {
	list := make([]*products.Product, len(items))
	keys := make([]map[string]*dynamodb.AttributeValue, len(items))

	for index, item := range items {
		keys[index] = map[string]*dynamodb.AttributeValue{key: {S: item}}
	}

	output, err := repository.DynamoDB.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems:           map[string]*dynamodb.KeysAndAttributes{*repository.tableName: {Keys: keys}},
		ReturnConsumedCapacity: nil,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Responses[*repository.tableName], &list)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repository *ProductRepository) FindMany(ids []*string) ([]*products.Product, error) {
	return repository.batchRequest("id", ids)
}

func (repository *ProductRepository) FindByCategoryID(ID *string) ([]*products.Product, error) {
	var items []*products.Product
	output, err := repository.DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":categoryId": {S: ID}},
		KeyConditionExpression:    aws.String("categoryId = :categoryId"),
		IndexName:                 aws.String("categoryId-index"),
		TableName:                 repository.tableName,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repository *ProductRepository) Delete(ID *string) error {
	_, err := repository.DynamoDB.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: ID}},
		TableName: repository.tableName,
	})

	return err
}

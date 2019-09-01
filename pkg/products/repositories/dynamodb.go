package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/products"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBProductRepository struct {
	DynamoDB  *dynamodb.DynamoDB
	tableName *string
}

func NewDynamoDBProductRepository(db *dynamodb.DynamoDB) *DynamoDBProductRepository {
	return &DynamoDBProductRepository{
		DynamoDB:  db,
		tableName: aws.String("products"),
	}
}

func (repository *DynamoDBProductRepository) Store(product *products.Product) error {
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

func (repository *DynamoDBProductRepository) Update(id *string, product *products.Product) error {
	item, err := dynamodbattribute.MarshalMap(product)

	if err != nil {
		return err
	}

	_, err = repository.DynamoDB.PutItem(&dynamodb.PutItemInput{
		ConditionExpression:       aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: id}},
		Item:                      item,
		TableName:                 repository.tableName,
	})

	return err
}
func (repository *DynamoDBProductRepository) FindOne(ID *string) (*products.Product, error) {
	item := &products.Product{}
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

func (repository *DynamoDBProductRepository) batchRequest(key string, items []*string) ([]*products.Product, error) {
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

func (repository *DynamoDBProductRepository) FindMany(ids []*string) ([]*products.Product, error) {
	return repository.batchRequest("id", ids)
}

func (repository *DynamoDBProductRepository) FindByCategoryID(ID *string) ([]*products.Product, error) {
	items := make([]*products.Product, 0)
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

func (repository *DynamoDBProductRepository) Delete(ID *string) error {
	_, err := repository.DynamoDB.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: ID}},
		TableName: repository.tableName,
	})

	return err
}

func (repository *DynamoDBProductRepository) All() ([]*products.Product, error) {
	items := make([]*products.Product, 0)
	scanInput := &dynamodb.ScanInput{TableName: repository.tableName}
	output, err := repository.DynamoDB.Scan(scanInput)

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	if err != nil {
		return nil, err
	}

	return items, nil
}

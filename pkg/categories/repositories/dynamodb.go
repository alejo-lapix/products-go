package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/categories"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBCategoryRepository struct {
	DynamoDB  *dynamodb.DynamoDB
	tableName *string
}

func NewDynamoDBCategoryRepository(db *dynamodb.DynamoDB) *DynamoDBCategoryRepository {
	tableName := "categories"

	return &DynamoDBCategoryRepository{
		DynamoDB:  db,
		tableName: &tableName,
	}
}

func (repository *DynamoDBCategoryRepository) MainCategories(limit, offset int) ([]*categories.Category, error) {
	items := make([]*categories.Category, 0)
	output, err := repository.DynamoDB.Scan(&dynamodb.ScanInput{
		FilterExpression: aws.String("attribute_not_exists(parentCategoryId)"),
		TableName:        repository.tableName,
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

func (repository *DynamoDBCategoryRepository) SubCategories(categoryID *string) ([]*categories.Category, error) {
	var mainCategories []*categories.Category
	output, err := repository.DynamoDB.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    aws.String("parentCategoryId = :categoryId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":categoryId": {S: categoryID}},
		IndexName:                 aws.String("parentCategoryId"),
		TableName:                 repository.tableName,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, mainCategories)

	if err != nil {
		return nil, err
	}

	return mainCategories, nil
}

func (repository *DynamoDBCategoryRepository) All(cursor *string) ([]*categories.Category, error) {
	currentCategory := make([]*categories.Category, 0)
	scanInput := &dynamodb.ScanInput{
		Limit:     aws.Int64(15),
		TableName: repository.tableName,
	}

	if *cursor != "" {
		scanInput.FilterExpression = aws.String("#name > :name")
		scanInput.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{":name": {S: cursor}}
		scanInput.ExpressionAttributeNames = map[string]*string{"#name": aws.String("name")}
	}

	output, err := repository.DynamoDB.Scan(scanInput)

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &currentCategory)

	if err != nil {
		return nil, err
	}

	return currentCategory, nil
}

func (repository *DynamoDBCategoryRepository) Find(ID *string) (*categories.Category, error) {
	var currentCategory *categories.Category
	output, err := repository.DynamoDB.GetItem(&dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: ID}},
		TableName: repository.tableName,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, currentCategory)

	if err != nil {
		return nil, err
	}

	return currentCategory, nil
}

func (repository *DynamoDBCategoryRepository) Store(category *categories.Category) error {
	item, err := dynamodbattribute.MarshalMap(category)

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

func (repository *DynamoDBCategoryRepository) Remove(ID *string) error {
	_, err := repository.DynamoDB.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: ID}},
		TableName: repository.tableName,
	})

	return err
}

func (repository *DynamoDBCategoryRepository) Update(ID *string, category *categories.Category) error {
	item, err := dynamodbattribute.MarshalMap(category)

	if err != nil {
		return err
	}

	_, err = repository.DynamoDB.PutItem(&dynamodb.PutItemInput{
		ConditionExpression:       aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: ID}},
		Item:                      item,
		TableName:                 repository.tableName,
	})

	return err
}

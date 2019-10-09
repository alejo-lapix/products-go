package repositories

import (
	"fmt"
	"github.com/alejo-lapix/products-go/pkg/categories"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBCategoryRepository struct {
	DynamoDB            *dynamodb.DynamoDB
	tableName           *string
	parentCategoryTries int
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
	output, err := repository.DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":yes":     {S: aws.String("y")},
			":visible": {BOOL: aws.Bool(true)},
		},
		FilterExpression:       aws.String("visible = :visible"),
		IndexName:              aws.String("isMainCategory-index"),
		KeyConditionExpression: aws.String("isMainCategory = :yes"),
		TableName:              repository.tableName,
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

func (repository *DynamoDBCategoryRepository) Total() (int64, error) {
	output, err := repository.DynamoDB.Scan(&dynamodb.ScanInput{
		ReturnConsumedCapacity: aws.String("TOTAL"),
		Select:                 aws.String("COUNT"),
		TableName:              repository.tableName,
	})

	if err != nil {
		return 0, err
	}

	return *output.Count, nil
}

func (repository *DynamoDBCategoryRepository) SubCategories(categoryID *string) ([]*categories.Category, error) {
	mainCategories := make([]*categories.Category, 0)
	output, err := repository.DynamoDB.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    aws.String("parentCategoryId = :categoryId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":categoryId": {S: categoryID}},
		IndexName:                 aws.String("parentCategoryId-index"),
		TableName:                 repository.tableName,
	})

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &mainCategories)

	if err != nil {
		return nil, err
	}

	return mainCategories, nil
}

func (repository *DynamoDBCategoryRepository) All() ([]*categories.Category, error) {
	currentCategory := make([]*categories.Category, 0)
	scanInput := &dynamodb.ScanInput{TableName: repository.tableName, IndexName: aws.String("id-name-index")}
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

func (repository *DynamoDBCategoryRepository) FindMainCategory(childCategoryID *string) (*categories.Category, error) {
	category, err := repository.Find(childCategoryID)

	defer repository.resetRetries()

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, fmt.Errorf("the given category does not have a valid parent")
	}

	if category.ParentCategoryID == nil || *category.ParentCategoryID == "" {
		return category, nil
	}

	repository.parentCategoryTries++

	if repository.parentCategoryTries > 10 {
		return nil, fmt.Errorf("the parent category could not be find while looking at %d categories", repository.parentCategoryTries)
	}

	return repository.FindMainCategory(category.ParentCategoryID)
}

func (repository *DynamoDBCategoryRepository) resetRetries() {
	repository.parentCategoryTries = 0
}

func (repository *DynamoDBCategoryRepository) Find(ID *string) (*categories.Category, error) {
	currentCategory := &categories.Category{}
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

func (repository *DynamoDBCategoryRepository) FindMany(items []*string) ([]*categories.Category, error) {
	list := make([]*categories.Category, len(items))
	keys := make([]map[string]*dynamodb.AttributeValue, len(items))

	for index, item := range items {
		keys[index] = map[string]*dynamodb.AttributeValue{"id": {S: item}}
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

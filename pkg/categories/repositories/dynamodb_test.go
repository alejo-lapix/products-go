package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/categories"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"testing"
)

var tableName = "categories"

func TestDynamoDBCategoryRepository_All(t *testing.T) {
	tests := []struct {
		name    string
		want    []*categories.Category
		wantErr bool
	}{
		{
			name:    "Must Scan the table",
			want:    make([]*categories.Category, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  dynamoDB(),
				tableName: &tableName,
			}
			got, err := repository.All()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_Total(t *testing.T) {
	type args struct {
		cursor *string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Must Scan the table and get the row count",
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  dynamoDB(),
				tableName: &tableName,
			}
			got, err := repository.Total()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_Find(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		ID *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *categories.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			got, err := repository.Find(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_MainCategories(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*categories.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			got, err := repository.MainCategories(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("MainCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MainCategories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_Remove(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		ID *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			if err := repository.Remove(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_Store(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		category *categories.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			if err := repository.Store(tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_SubCategories(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		categoryID *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*categories.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			got, err := repository.SubCategories(tt.args.categoryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubCategories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoDBCategoryRepository_Update(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		ID       *string
		category *categories.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &DynamoDBCategoryRepository{
				DynamoDB:  tt.fields.DynamoDB,
				tableName: tt.fields.tableName,
			}
			if err := repository.Update(tt.args.ID, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDynamoDBCategoryRepository(t *testing.T) {
	type args struct {
		db *dynamodb.DynamoDB
	}
	tests := []struct {
		name string
		args args
		want *DynamoDBCategoryRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDynamoDBCategoryRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDynamoDBCategoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dynamoDB() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})

	if err != nil {
		panic(err)
	}

	return dynamodb.New(sess)
}

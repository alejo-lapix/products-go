package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/products"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

var tableName = "productos"
var ID = uuid.New().String()

func TestProductRepository_Store(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		product *products.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Insert an item",
			fields: fields{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			},
			args: args{product: &products.Product{
				ID:          aws.String(ID),
				Name:        aws.String("Name"),
				Price:       aws.Float64(5000),
				Description: aws.String("Description"),
				CategoryID:  aws.String("aaaaa"),
				Multimedia:  nil,
				UnitOfMeasurement: &products.UnitOfMeasurement{
					Quantity: aws.Float64(1),
					Unit:     aws.String("litro"),
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			if err := repository.Store(tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductRepository_Update(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		product *products.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Updates an element",
			args: args{product: &products.Product{
				ID:                aws.String(ID),
				Name:              aws.String("bbbbb"),
				Price:             aws.Float64(6000),
				Description:       aws.String("New Descriptoin"),
				CategoryID:        aws.String("bbbbb"),
				Multimedia:        nil,
				UnitOfMeasurement: nil,
				CreatedAt:         nil,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			if err := repository.Update(tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductRepository_FindByCategoryID(t *testing.T) {
	type fields struct {
		DynamoDB  *dynamodb.DynamoDB
		tableName *string
	}
	type args struct {
		id *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*products.Product
		wantErr bool
	}{
		{
			name:    "Cagetories",
			args:    args{id: aws.String("bbbbb")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			got, err := repository.FindByCategoryID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByCategoryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if false && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByCategoryID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_FindMany(t *testing.T) {
	type fields struct {
		DynamoDB *dynamodb.DynamoDB
	}
	type args struct {
		ids []*string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*products.Product
		wantErr bool
	}{
		{
			name:   "Get Many Products",
			fields: fields{DynamoDB: dynamoDBInstance()},
			args: args{
				ids: []*string{aws.String("aaaaa"), aws.String("bbbbb")},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			got, err := repository.FindMany(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if false && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMany() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_FindOne(t *testing.T) {
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
		want    *products.Product
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			got, err := repository.FindOne(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductRepository_Delete(t *testing.T) {
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
			repository := &ProductRepository{
				DynamoDB:  dynamoDBInstance(),
				tableName: &tableName,
			}
			if err := repository.Delete(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func dynamoDBInstance() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		panic(err.Error())
	}

	return dynamodb.New(sess)
}

package repositories

import (
	"github.com/alejo-lapix/products-go/pkg/categories"
	"reflect"
	"testing"
	"time"
)

type repository struct{}

func s(input string) *string {
	return &input
}

func tm(input time.Time) *time.Time {
	return &input
}

func commonElements() []*categories.Category {
	return []*categories.Category{
		{ID: s("abcd")},
		{ID: s("efgh")},
		{ID: s("ijkl")},
	}
}

func temporalElements() []*categories.Category {
	return []*categories.Category{
		{ID: s("mnop")},
		{ID: s("qrst")},
		{ID: s("uvxy")},
		{ID: s("z123")},
	}
}

func (repo repository) MainCategories(limit, offset int) ([]*categories.Category, error) {
	return commonElements(), nil
}

func (repo repository) SubCategories(categoryID *string) ([]*categories.Category, error) {
	return []*categories.Category{}, nil
}
func (repo repository) Find(ID *string) (*categories.Category, error) {
	return &categories.Category{}, nil
}
func (repo repository) Store(*categories.Category) error {
	return nil
}
func (repo repository) Remove(ID *string) error {
	return nil
}
func (repo repository) Update(ID *string, category *categories.Category) error {
	return nil
}
func (repo repository) All() ([]*categories.Category, error) {
	return []*categories.Category{}, nil
}
func (repo repository) Total() (int64, error) {
	return 3, nil
}

func Test_inMemory_Remember(t *testing.T) {
	type fields struct {
		elements   map[string]interface{}
		timeStamps map[string]*time.Time
	}
	type args struct {
		key      string
		seconds  int
		callback func() (interface{}, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Should call the function passed",
			fields: fields{
				elements:   map[string]interface{}{},
				timeStamps: map[string]*time.Time{},
			},
			args: args{
				key:      "empty",
				seconds:  1,
				callback: func() (interface{}, error) { return "cached", nil },
			},
			want:    "cached",
			wantErr: false,
		},
		{
			name: "Should take items from memory",
			fields: fields{
				elements:   map[string]interface{}{"not empty": "cached"},
				timeStamps: map[string]*time.Time{},
			},
			args: args{
				key:      "not empty",
				seconds:  1,
				callback: func() (interface{}, error) { return nil, nil },
			},
			want:    "cached",
			wantErr: false,
		},
		{
			name: "Should call the function if time passed",
			fields: fields{
				elements:   map[string]interface{}{"old": "cached"},
				timeStamps: map[string]*time.Time{"old": tm(time.Now().Add(time.Second * -5))},
			},
			args: args{
				key:      "old",
				seconds:  0,
				callback: func() (interface{}, error) { return "new cached", nil },
			},
			want:    "new cached",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &inMemory{
				elements:   tt.fields.elements,
				timeStamps: tt.fields.timeStamps,
			}
			got, err := driver.Remember(tt.args.key, tt.args.seconds, tt.args.callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

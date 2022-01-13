package adapters

import (
	"context"
	"github.com/obarbier/awesome-crypto/common"
	"github.com/obarbier/awesome-crypto/crypto_service_api/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
)

func TestMongoKeyRepository_Get(t *testing.T) {
	type fields struct {
		MongoConnection *common.MongoConnection
		db              *mongo.Database
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantK   *domain.KeyEntity
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MongoKeyRepository{
				MongoConnection: tt.fields.MongoConnection,
				db:              tt.fields.db,
			}
			gotK, err := m.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("Get() gotK = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}

func TestMongoKeyRepository_Save(t *testing.T) {
	type fields struct {
		MongoConnection *common.MongoConnection
		db              *mongo.Database
	}
	type args struct {
		ctx context.Context
		k   domain.KeyEntity
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
			m := MongoKeyRepository{
				MongoConnection: tt.fields.MongoConnection,
				db:              tt.fields.db,
			}
			if err := m.Save(tt.args.ctx, tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMongoKeyRepository(t *testing.T) {
	tests := []struct {
		name    string
		want    *MongoKeyRepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMongoKeyRepository()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMongoKeyRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMongoKeyRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}

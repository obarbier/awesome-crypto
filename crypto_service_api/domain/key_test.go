package domain

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyService_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockKeyRepo := NewMockKeyRepository(mockCtrl)

	type fields struct {
		repo KeyRepository
	}
	type args struct {
		ctx     context.Context
		keyName string
		ownerId string
		keyType string
		opt     []map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    Key
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Testing create key",
			fields: fields{repo: mockKeyRepo},
			mock: func() {
				mockKeyRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				ctx:     context.Background(),
				keyName: "test-1",
				ownerId: "abc",
				keyType: "rsa",
				opt:     make([]map[string]interface{}, 1),
			},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			k := KeyService{
				repo: tt.fields.repo,
			}
			got, err := k.Create(tt.args.ctx, tt.args.ownerId, tt.args.keyName, tt.args.keyType, tt.args.opt...)
			if !tt.wantErr(t, err, fmt.Sprintf("Create(%v, %v, %v, %v)", tt.args.ctx, tt.args.ownerId, tt.args.keyType, tt.args.opt)) {
				return
			}

			assert.Equalf(t, tt.want, got, "Create(%v, %v, %v, %v)", tt.args.ctx, tt.args.ownerId, tt.args.keyType, tt.args.opt)
		})
	}
}

func TestKeyService_GetKeyById(t *testing.T) {
	type fields struct {
		repo KeyRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Key
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KeyService{
				repo: tt.fields.repo,
			}
			got, err := k.GetKeyById(tt.args.ctx, tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("GetKeyById(%v, %v)", tt.args.ctx, tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetKeyById(%v, %v)", tt.args.ctx, tt.args.id)
		})
	}
}

func TestNewKeyService(t *testing.T) {
	type args struct {
		repo KeyRepository
	}
	tests := []struct {
		name string
		args args
		want KeyService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewKeyService(tt.args.repo), "NewKeyService(%v)", tt.args.repo)
		})
	}
}

package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Key interface {
	Type() string
	UnMarshal() ([]byte, error)
	Marshal(string) error
}

type KeyEntity struct {
	Id         string `json:"id" bson:"_id,omitempty"`
	KeyName    string `json:"keyName" bson:"keyName,omitempty"`
	KeyOwnerId string `json:"keyOwnerId" bson:"keyOwnerId,omitempty"`
	KeyType    string `json:"keyType" bson:"keyType,omitempty"`
	KeyString  string `json:"keyString" bson:"keyString,omitempty"`
}

//go:generate mockgen -destination=./mock_user_repo.go -package=domain github.com/obarbier/awesome-crypto/crypto_service_api/domain KeyRepository
type KeyRepository interface {
	Save(ctx context.Context, k KeyEntity) error
	Get(ctx context.Context, id string) (k *KeyEntity, err error)
}

var _ IKeyService = KeyService{}

type IKeyService interface {
	Create(ctx context.Context, ownerId, keyName, keyType string, opt ...map[string]interface{}) (Key, error)
	GetKeyById(ctx context.Context, id string) (Key, error)
}

type KeyService struct {
	repo KeyRepository
}

func NewKeyService(repo KeyRepository) KeyService {
	return KeyService{repo: repo}
}

func (k KeyService) Create(ctx context.Context, ownerId, keyName, keyType string, opt ...map[string]interface{}) (Key, error) {
	var key Key
	switch keyType {
	case "rsa", "RSA": // FIXME Create type for this
		var rsaErr error
		rsaKey, rsaErr := NewRsaKey(2048)
		if rsaErr != nil {
			return nil, rsaErr
		}
		key = rsaKey
	default:
		return nil, fmt.Errorf("not a valid KeyType")

	}
	// SAVE the keys
	keyByte, err := key.UnMarshal()
	if err != nil {
		return nil, err
	}
	// create a new Id
	id, uuidErr := uuid.NewUUID()
	if uuidErr != nil {
		return nil, uuidErr
	}
	entity := KeyEntity{
		Id:         id.String(),
		KeyName:    keyName,
		KeyOwnerId: ownerId,
		KeyType:    keyType,
		KeyString:  string(keyByte),
	}
	err = k.repo.Save(ctx, entity)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (k KeyService) GetKeyById(ctx context.Context, id string) (Key, error) {
	entity, err := k.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	var key Key
	switch entity.KeyType {
	case "rsa", "RSA": // FIXME Create type for this
		rsaKey := &RSA{}
		err := rsaKey.Marshal(entity.KeyString)
		if err != nil {
			return nil, err
		}
		key = rsaKey
	default:
		return nil, fmt.Errorf("not a valid KeyType")
	}
	return key, nil
}

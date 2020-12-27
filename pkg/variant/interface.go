//go:generate mockgen -source interface.go -destination variant_mock.go -package variant

package variant

import "github.com/markus-azer/products-service/pkg/entity"

//MessagesReader Reader interface
type messagesReader interface {
}

//MessagesWriter variant writer
type messagesWriter interface {
	SendMessage(m *entity.Message)
	SendMessages(messages []*entity.Message)
}

//MessagesRepository repository interface
type MessagesRepository interface {
	messagesReader
	messagesWriter
}

//StoreReader variant reader interface
type storeReader interface {
	FindOneByID(id entity.ID) (*entity.Variant, error)
	FindOneByAttribute(product entity.ID, attributes map[string]string) (*entity.Variant, error)
}

//StoreWriter variant writer interface
type storeWriter interface {
	StoreCommand(c *entity.Command) (*entity.ID, error)
	Create(variant *entity.Variant) (*entity.ID, error)
	UpdateOne(id entity.ID, variant *entity.UpdateVariant, version entity.Version) (int, error)
	DeleteOne(id entity.ID, version entity.Version) (int, error)
}

//StoreRepository product store repository interface
type StoreRepository interface {
	storeReader
	storeWriter
}

//Reader interface
type reader interface {
}

//Writer interface
type writer interface {
	Create(createVariantDTO CreateVariantDTO) (entity.ID, int32, []entity.ClientError)
	UpdateOne(id entity.ID, version int32, updateVariantDTO UpdateVariantDTO) (int32, []entity.ClientError)
	Delete(id entity.ID, version int32) *entity.ClientError
}

//UseCase use case interface
type UseCase interface {
	reader
	writer
}

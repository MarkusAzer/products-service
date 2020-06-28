//go:generate mockgen -source interface.go -destination brand_mock.go -package brand

package brand

import (
	"github.com/MarkusAzer/products-service/pkg/entity"
)

//messagesReader brand reader interface
type messagesReader interface {
	GetMessages() <-chan entity.Message
}

//MessagesRepository repository interface
type MessagesRepository interface {
	messagesReader
}

//StoreReader brand reader interface
type storeReader interface {
	FindOneByID(id entity.ID) (entity.Brand, error)
	FindOneByName(name string) (*entity.Brand, error)
}

//StoreWriter brand writer interface
type storeWriter interface {
	Create(b *entity.Brand)
	UpdateOne(id entity.ID, b *entity.Brand)
	DeleteOne(id entity.ID)
}

//StoreRepository brand store repository interface
type StoreRepository interface {
	storeReader
	storeWriter
}

//UseCase use case interface
type UseCase interface {
	storeReader
	storeWriter
}

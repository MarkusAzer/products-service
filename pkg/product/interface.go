//go:generate mockgen -source interface.go -destination product_mock.go -package product

package product

import "github.com/markus-azer/products-service/pkg/entity"

//MessagesReader Reader interface
type messagesReader interface {
}

//MessagesWriter product writer
type messagesWriter interface {
	SendMessage(m *entity.Message)
	SendMessages(messages []*entity.Message)
}

//MessagesRepository repository interface
type MessagesRepository interface {
	messagesReader
	messagesWriter
}

//StoreReader product reader interface
type storeReader interface {
	FindOneByID(id entity.ID) (*entity.Product, error)
}

//StoreWriter product writer interface
type storeWriter interface {
	StoreCommand(c *entity.Command) (*entity.ID, error)
	Create(p *entity.Product) (*entity.ID, error)
	UpdateOne(id entity.ID, p *entity.Product, v entity.Version) (int, error)
	UpdateOneP(id entity.ID, p *entity.UpdateProduct, v entity.Version) (int, error)
	DeleteOne(id entity.ID, v entity.Version) (int, error)
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
	Create(createProductDTO CreateProductDTO) (*entity.ID, *entity.Version, *entity.Error)
	UpdateOne(id entity.ID, v int32, updateProductDTO UpdateProductDTO) (*int32, *entity.Error)
	Delete(id entity.ID, version int32) *entity.Error
}

//UseCase use case interface
type UseCase interface {
	reader
	writer
}

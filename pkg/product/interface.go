package product

import "github.com/MarkusAzer/products-service/pkg/entity"

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
	FindOneByID(id entity.ID) (entity.Product, error)
}

//StoreWriter product writer interface
type storeWriter interface {
	StoreCommand(c *entity.Command)
	Create(p *entity.Product)
	UpdateOne(id entity.ID, p *entity.Product)
	UpdateOneP(id entity.ID, p *entity.UpdateProduct)
	DeleteOne(id entity.ID)
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
	Create(e *entity.Product) (entity.ID, error)
	UpdateOne(id entity.ID, version int32, p *entity.UpdateProduct) (int32, error)
	Publish(id entity.ID, version int32) (int32, error)
	Unpublish(id entity.ID, version int32) (int32, error)
	UpdatePrice(id entity.ID, version int32, price int) (int32, error)
	Delete(id entity.ID, version int32) error
}

//UseCase use case interface
type UseCase interface {
	reader
	writer
}

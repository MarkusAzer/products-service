package brand

import "github.com/markus-azer/products-service/pkg/entity"

//Service service interface
type Service struct {
	repo StoreRepository
}

//NewService create new service
func NewService(r StoreRepository) *Service {
	return &Service{
		repo: r,
	}
}

//FindOneByID new product
func (s *Service) FindOneByID(id entity.ID) (entity.Brand, error) {
	return s.repo.FindOneByID(id)
}

//FindOneByName new product
func (s *Service) FindOneByName(name string) (*entity.Brand, error) {
	return s.repo.FindOneByName(name)
}

//Create new product
func (s *Service) Create(b *entity.Brand) {
	s.repo.Create(b)
}

//UpdateOne brand
func (s *Service) UpdateOne(id entity.ID, b *entity.Brand) {
	s.repo.UpdateOne(id, b)
}

//DeleteOne brand
func (s *Service) DeleteOne(id entity.ID) {
	s.repo.DeleteOne(id)
}

package service

import (
	"gocouchbase/pkg/storage/couchbase"
)

type Service struct {
	couchbaseRepo couchbase.CouchbaseRepo
}

func NewDefault() Service {
	return Service{
		couchbaseRepo: couchbase.NewCouchbaseRepo()}
}

func New(cb couchbase.CouchbaseRepo) Service {
	return Service{
		couchbaseRepo: cb}
}

func (s *Service) GetAPIClient(clientID string) (client couchbase.APIClient, err error) {
	return s.couchbaseRepo.GetAPIClient(clientID)
}

func (s *Service) StoreAPIClient(client couchbase.APIClient) error {
	return s.couchbaseRepo.StoreAPIClient(client)
}

func (s *Service) GetServiceRole(roleName string) (role couchbase.ServiceRole, err error) {
	return s.couchbaseRepo.GetServiceRole(roleName)
}

func (s *Service) GetServiceRoles() (roles []couchbase.ServiceRole, err error) {
	return s.couchbaseRepo.GetServiceRoles()
}

func (s *Service) StoreServiceRole(role couchbase.ServiceRole) error {
	return s.couchbaseRepo.StoreServiceRole(role)
}

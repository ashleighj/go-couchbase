package couchbase

type CouchbaseRepo interface {
	GetAPIClient(clientID string) (APIClient, error)
	StoreAPIClient(client APIClient) (err error)

	GetServiceRoles() ([]ServiceRole, error)
	GetServiceRole(roleName string) (ServiceRole, error)
	StoreServiceRole(role ServiceRole) error
}

type couchbaseRepo struct {
	cluster CouchbaseCluster
}

func NewCouchbaseRepo() *couchbaseRepo {
	repo := &couchbaseRepo{}
	repo.Init()
	return repo
}

func (c *couchbaseRepo) GetAPIClient(clientID string) (client APIClient, err error) {
	err = c.cluster.GetDocument(clientID, &client)
	return
}

func (c *couchbaseRepo) StoreAPIClient(client APIClient) (err error) {
	err = c.cluster.Upsert(client.ClientID, client, 0)
	return
}

func (c *couchbaseRepo) StoreServiceRole(role ServiceRole) (err error) {
	err = c.cluster.Upsert(role.RoleName, role, 0)
	return
}

func (c *couchbaseRepo) GetServiceRole(roleName string) (role ServiceRole, err error) {
	err = c.cluster.GetDocument(roleName, &role)
	return
}

func (c *couchbaseRepo) GetServiceRoles() (roles []ServiceRole, err error) {
	interfaces, err := c.cluster.GetAll(ServiceRole{})
	for _, r := range interfaces {
		role := r.(*ServiceRole)
		roles = append(roles, *role)
	}
	return
}

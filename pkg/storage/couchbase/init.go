package couchbase

import (
	"gocouchbase/pkg/config"
	"gocouchbase/pkg/log"

	"github.com/couchbase/gocb"
)

var (
	cbConnect = "couchbase://"
	err       error
)

type CouchbaseCluster struct {
	Cluster       *gocb.Cluster
	ConnectString string
	Buckets       []couchbaseBucket
}

type couchbaseBucket struct {
	Name    string
	DocType string
	Bucket  *gocb.Bucket
}

func (c *couchbaseRepo) Init() {
	c.cluster = CouchbaseCluster{}
	c.cluster.Init()
}

func (cluster *CouchbaseCluster) Init() {
	conf := config.Get()
	cluster.ConnectString = cbConnect + conf.CouchbaseHost

	cluster.Cluster, err = gocb.Connect(cluster.ConnectString)
	if err != nil {
		log.Error(nil, err, log.ErrCouchbaseConnection)
	}

	err = cluster.Cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: conf.CouchbaseUsername,
		Password: conf.CouchbasePassword,
	})

	if err != nil {
		log.Error(nil, err, log.ErrCouchbaseAuth)
	}

	cluster.connectBuckets()
}

func (cluster *CouchbaseCluster) connectBuckets() {
	cluster.Buckets = []couchbaseBucket{
		{Name: "api_client", DocType: "couchbase.APIClient"},
		{Name: "service_roles", DocType: "couchbase.ServiceRole"}}

	for i, bucket := range cluster.Buckets {
		cluster.Buckets[i].Bucket, err = cluster.Cluster.OpenBucket(bucket.Name, "")
		if err != nil {
			log.Error(nil, err, log.ErrCouchbaseBucketOpen, bucket.Name)
		}
	}
}

package couchbase

import (
	"reflect"
	"strings"
)

func (cluster *CouchbaseCluster) getBucketByDocType(docType reflect.Type) couchbaseBucket {
	for _, bucket := range cluster.Buckets {
		if strings.Contains(docType.String(), bucket.DocType) {
			return bucket
		}
	}
	return couchbaseBucket{}
}

package couchbase

import (
	"fmt"
	"reflect"

	"gocouchbase/pkg/log"

	"github.com/couchbase/gocb"
)

const (
	// bucket name, bucket name
	getAllQuery    = "SELECT `%s`.* FROM `%s`"
	getAllIdsQuery = "SELECT meta().id FROM `%s`"
)

type BulkInsertItem struct {
	Key  string
	Item interface{}
}

func (c *CouchbaseCluster) Upsert(id string, document interface{}, expiry ...uint32) (err error) {
	bucket := c.getBucketByDocType(reflect.TypeOf(document))
	if bucket.Bucket == nil {
		err = log.GetError(log.ErrCouchbaseBucketRetrieval)
		return
	}

	expire := uint32(0)
	if expiry != nil {
		expire = expiry[0]
	}

	cas, err := bucket.Bucket.Insert(id, document, expire)
	if err != nil {
		_, err = bucket.Bucket.Replace(id, document, cas, expire)
		if err != nil {
			return
		}
	}
	return
}

func (c *CouchbaseCluster) GetDocument(id string, document interface{}) (err error) {
	bucket := c.getBucketByDocType(reflect.Indirect(reflect.ValueOf(document)).Type())
	if bucket.Bucket == nil {
		err = log.GetError(log.ErrCouchbaseBucketRetrieval)
		return
	}

	_, err = bucket.Bucket.Get(id, document)
	return
}

/*
*	Works by passing a reference to an empty struct matching
*	the document type to be fetched.
 */
func (c *CouchbaseCluster) GetAll(val interface{}) (all []interface{}, err error) {
	valType := reflect.Indirect(reflect.ValueOf(val)).Type()

	bucket := c.getBucketByDocType(valType)
	if bucket.Bucket == nil {
		err = log.GetError(log.ErrCouchbaseBucketRetrieval)
		return
	}

	query := gocb.NewN1qlQuery(fmt.Sprintf(getAllIdsQuery, bucket.Name))

	rows, err := bucket.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil || rows == nil {
		return
	}

	var ids []string
	type idResult struct {
		ID string `json:"id"`
	}
	var row idResult
	for rows.Next(&row) {
		ids = append(ids, row.ID)
	}

	var items []gocb.BulkOp
	for _, id := range ids {
		value := reflect.New(valType).Interface()
		item := &gocb.GetOp{Key: id, Value: value}
		items = append(items, item)
	}

	err = bucket.Bucket.Do(items)
	if err != nil {
		return
	}

	for _, item := range items {
		getItem := item.(*gocb.GetOp)
		val := reflect.ValueOf(getItem.Value)
		all = append(all, val.Interface())
	}

	return
}

func (c *CouchbaseCluster) BulkInsert(items []BulkInsertItem) error {
	bucket := c.getBucketByDocType(reflect.TypeOf(items[0].Item))
	if bucket.Bucket == nil {
		return log.GetError(log.ErrCouchbaseBucketRetrieval)
	}

	var insertItems []gocb.BulkOp
	for _, item := range items {
		key := item.Key
		item := item.Item
		insertItem := &gocb.InsertOp{Key: key, Value: item}
		insertItems = append(insertItems, insertItem)
	}

	return bucket.Bucket.Do(insertItems)
}

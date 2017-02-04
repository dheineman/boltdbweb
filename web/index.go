package boltbrowserweb

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/pkg/errors"
)

var Db *bolt.DB

func Index(c *gin.Context) {

	c.Redirect(301, "/web/html/layout.html")

}

func CreateBucket(c *gin.Context) {

	if c.PostForm("bucket") == "" {
		c.String(200, "no bucket name | n")
	}

	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		b = b
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	c.String(200, "ok")

}

func DeleteBucket(c *gin.Context) {

	if c.PostForm("bucket") == "" {
		c.String(200, "no bucket name | n")
	}

	Db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(c.PostForm("bucket")))

		if err != nil {

			c.String(200, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		return nil
	})

	c.String(200, "ok")

}

func DeleteKey(c *gin.Context) {

	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		c.String(200, "no bucket name or key | n")
	}

	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		b = b
		if err != nil {

			c.String(200, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		err = b.Delete([]byte(c.PostForm("key")))

		if err != nil {

			c.String(200, "error Deleting KV | n")
			return fmt.Errorf("delete kv: %s", err)
		}

		return nil
	})

	c.String(200, "ok")

}

func Put(c *gin.Context) {

	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		c.String(200, "no bucket name or key | n")
	}

	Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		b = b
		if err != nil {

			c.String(200, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(c.PostForm("key")), []byte(c.PostForm("value")))

		if err != nil {

			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})

	c.String(200, "ok")

}

func Get(c *gin.Context) {

	res := []string{"nok", ""}

	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		res[1] = "no bucket name or key | n"
		c.JSON(200, res)
	}

	// Get a list of buckets
	v := c.PostForm("bucket")
	bs := strings.Split(v, "/")

	Db.View(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, bs)
		if err != nil {
			res[1] = "Bucket not found"
			return fmt.Errorf(res[1])
		}

		if b != nil {

			v := b.Get([]byte(c.PostForm("key")))

			res[0] = "ok"
			res[1] = string(v)

			fmt.Printf("Key: %s\n", v)

		} else {

			res[1] = "error opening bucket| does it exist? | n"

		}
		return nil

	})

	c.JSON(200, res)

}

func Buckets(c *gin.Context) {

	res := []string{}

	Db.View(func(tx *bolt.Tx) error {

		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {

			b := []string{string(name)}
			res = append(res, b...)
			return nil
		})

	})

	c.JSON(200, res)

}

type Result struct {
	Result  string
	Buckets []string
	N       map[string]Node
}

type Node struct {
	Key    string
	Value  string
	Bucket bool
}

func Explore(c *gin.Context) {

	var prefix []byte

	res := Result{Result: "nok"}
	res.N = make(map[string]Node)

	if c.PostForm("bucket") == "" {
		res.Result = "No bucket provided"
		c.JSON(200, res)
		return
	}

	if c.PostForm("key") != "" {
		prefix = []byte(c.PostForm("key"))
	}

	// Get a list of buckets
	v := c.PostForm("bucket")
	bs := strings.Split(v, "/")
	res.Buckets = bs

	Db.View(func(tx *bolt.Tx) (err error) {
		b, err := getBucket(tx, bs)
		if err != nil {
			res.Result = "Bucket not found"
			return fmt.Errorf(res.Result)
		}

		c := b.Cursor()

		if len(prefix) > 0  {
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				bucket := false
				if v == nil {
					bucket = true
				}

				res.N[string(k)] = Node{
					Key: string(k),
					Value: string(v),
					Bucket: bucket,
				}
			}

		} else {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				bucket := false
				if v == nil {
					bucket = true
				}

				res.N[string(k)] = Node{
					Key: string(k),
					Value: string(v),
					Bucket: bucket,
				}
			}
		}

		res.Result = "ok"

		return
	})

	c.JSON(200, res)

}

// getBucket - Traverse trough the given buckets to deepest one
func getBucket(tx *bolt.Tx, bs []string) (b *bolt.Bucket, err error) {
	// Get the first bucket
	nb, bs := bs[0], bs[1:]
	b = tx.Bucket([]byte(nb))
	if b == nil {
		err = errors.New("Bucket not found")
		return
	}

	// Keep walking buckets
	for len(bs) > 0 {
		var nb string
		nb, bs = bs[0], bs[1:]

		cb := b.Bucket([]byte(nb))
		if cb == nil {
			err = errors.New("Bucket not found")
			return
		}

		b = cb
	}

	return
}

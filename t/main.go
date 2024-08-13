package main

import (
	"log"

	"github.com/averystampp/bl"
	"github.com/averystampp/sesame"
	bolt "go.etcd.io/bbolt"
)

func main() {

	db, err := bolt.Open("posts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("posts"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	db.Close()

	rtr := sesame.NewRouter()
	routes(rtr)

	rtr.StartServer(":5000")
}

func routes(rtr *sesame.Router) {
	rtr.Post("/post/create", bl.NewPost)

	rtr.Get("/post/all", bl.AllPosts)

}

package bl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/averystampp/sesame"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Post struct {
	ID      uuid.UUID
	Author  string
	Title   string
	Content template.HTML

	Metadata PostMetadata
}

type PostMetadata struct {
	DatePosted  time.Time
	DateUpdated time.Time
	DateDeleted time.Time
	IsLive      bool
}

func NewPost(ctx sesame.Context) error {
	if ctx.Request().Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("body must be a json object")
	}

	var post Post = Post{
		ID: uuid.New(),
		Metadata: PostMetadata{
			DateUpdated: time.Now(),
			DatePosted:  time.Now(),
			IsLive:      false,
		},
	}
	err := json.NewDecoder(ctx.Request().Body).Decode(&post)
	if err != nil {
		return err
	}

	db, err := bolt.Open("posts.db", 0600, nil)
	if err != nil {
		return err
	}
	a := template.HTML("")

	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))

		if b == nil {
			return fmt.Errorf("database could not open \"posts\" bucket")
		}

		body, err := json.Marshal(&post)
		if err != nil {
			return err
		}
		return b.Put([]byte(post.ID.String()), body)
	})
}

func AllPosts(ctx sesame.Context) error {
	db, err := bolt.Open("posts.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	var posts []Post
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		err = b.ForEach(func(k, v []byte) error {
			var post Post
			err := json.Unmarshal(v, &post)
			if err != nil {
				return err
			}
			posts = append(posts, post)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	ctx.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(ctx.Response()).Encode(&posts)
}

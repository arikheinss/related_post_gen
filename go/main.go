package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const topN = 5
const InitialTagMapSize = 100
const InitialPostsSliceCap = 0

type isize uint32

type Post struct {
	ID    string   `json:"_id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type RelatedPosts struct {
	ID      string      `json:"_id"`
	Tags    *[]string   `json:"tags"`
	Related [topN]*Post `json:"related"`
}

type PostWithSharedTags struct {
	Post       isize
	SharedTags isize
}

func main() {
	file, _ := os.Open("../posts.json")
	var posts = make([]Post, 0, InitialPostsSliceCap)
	err := json.NewDecoder(file).Decode(&posts)
	if err != nil {
		fmt.Println(err)
	}
	postsLen := len(posts)
	start := time.Now()
	// assumes that there are less than 100 tags
	tagMap := make(map[string][]isize, InitialTagMapSize)

	for i, post := range posts {
		for _, tag := range post.Tags {
			tagMap[tag] = append(tagMap[tag], isize(i))
		}
	}

	allRelatedPosts := make([]RelatedPosts, postsLen)
	taggedPostCount := make([]isize, postsLen)

	for i := range posts {
		for j := range taggedPostCount {
			taggedPostCount[j] = 0
		}
		// Count the number of tags shared between posts
		for _, tag := range posts[i].Tags {
			for _, otherPostIdx := range tagMap[tag] {
				taggedPostCount[otherPostIdx]++
			}
		}
		taggedPostCount[i] = 0 // Don't count self
		top5 := [topN]PostWithSharedTags{}

		for j, count := range taggedPostCount {
			if count > top5[0].SharedTags {
				top5[0] = PostWithSharedTags{Post: isize(j), SharedTags: count}
				for k := 0; k < topN-1; k++ {
					if top5[k].SharedTags > top5[k+1].SharedTags {
						top5[k], top5[k+1] = top5[k+1], top5[k]
					}
				}
			}
		}
		// Convert indexes back to Post pointers
		topPosts := [topN]*Post{}
		for idx, t := range top5 {
			topPosts[idx] = &posts[t.Post]
		}

		allRelatedPosts[i] = RelatedPosts{
			ID:      posts[i].ID,
			Tags:    &posts[i].Tags,
			Related: topPosts,
		}
	}

	fmt.Println("Processing time (w/o IO):", time.Since(start))
	file, _ = os.Create("../related_posts_go.json")
	_ = json.NewEncoder(file).Encode(allRelatedPosts)
}

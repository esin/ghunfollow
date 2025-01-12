package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	gitHubToken := os.Getenv("GITHUB_TOKEN")
	if gitHubToken == "" {
		log.Panicln("Env variable GITHUB_TOKEN not found")
	}

	gitHubUserName := os.Getenv("GITHUB_USERNAME")
	if gitHubUserName == "" {
		log.Panicln("Env variable GITHUB_USERNAME not found")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	currentUser, _, err := ghClient.Users.Get(ctx, gitHubUserName)
	if err != nil {
		log.Panicln(err)
	}

	allFollowing := []*github.User{}
	allFollowing_q := make(map[string]bool)

	followingsPerPage := 100

	log.Printf("Following count: %d", *currentUser.Following)
	log.Println("Finding all your followings...")

	for i := 0; i < (*currentUser.Following/followingsPerPage)+1; i++ {

		following, _, err := ghClient.Users.ListFollowing(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
		if err != nil {
			log.Println(err)
		}

		for _, fl := range following {
			allFollowing = append(allFollowing, fl)
			allFollowing_q[*fl.Login] = true
			//TODO: Need to rewrite
		}
		time.Sleep(300 * time.Millisecond)
	}

	log.Println("Finding all your followings...done")

	log.Println("Try to unfollow")
	unFollowCount := 0

	for _, f := range allFollowing {
		_, err := ghClient.Users.Unfollow(ctx, *f.Login)
		if err != nil {
			log.Println(err)
			continue
		}
		unFollowCount++
	}
	log.Println("Try to unfollow...done")
	log.Println("Unfollow count:", unFollowCount)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/redigo/redis"
	"github.com/lycoris0731/go-groovecoaster/groovecoaster"
)

var c redis.Conn

type item struct {
	name string
	diff []string
}

func (i item) String() string {
	text := i.name + "\n"

	for j, v := range i.diff {
		switch groovecoaster.Difficulty(j) {
		case groovecoaster.Simple:
			text += "  [S] "
		case groovecoaster.Normal:
			text += "  [N] "
		case groovecoaster.Hard:
			text += "  [H] "
		case groovecoaster.Extra:
			text += "  [E] "
		}

		text += v + "\n"
	}

	return text + "\n"
}

func lastDateFromRedis() (time.Time, error) {
	log.Println("Fetching last played date from Redis...")

	exists, err := redis.Bool(c.Do("EXISTS", "lastDate"))
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot read lastDate from Redis: %s", err)
	}

	if !exists {
		log.Println("Seted lastDate")
		c.Do("SET", "lastDate", time.Now().Format("2006-01-02 15:04:05"))
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")

	lastDateString, err := redis.String(c.Do("GET", "lastDate"))
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot get lastDate from Redis: %s", err)
	}

	lastDate, err := time.ParseInLocation("2006-01-02 15:04:05", lastDateString, loc)

	log.Printf("Last played date: %s", lastDate)

	return lastDate, nil
}

func splitMusicSummary(summary []*groovecoaster.MusicSummary, lastDate time.Time) []*groovecoaster.MusicSummary {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	for i, m := range summary {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", m.LastPlayTime, loc)
		if err != nil {
			log.Print(err)
			return nil
		}

		if date.Equal(lastDate) || lastDate.After(date) {
			return summary[:i]
		}
	}

	return []*groovecoaster.MusicSummary{}
}

func musicFromRedis(id string) (*groovecoaster.MusicDetail, error) {
	exists, err := redis.Bool(c.Do("EXISTS", id))
	if err != nil {
		return nil, fmt.Errorf("cannot confirm whether music is exists from Redis: %s", err)
	}
	if !exists {
		return nil, nil
	}

	musicString, err := redis.String(c.Do("GET", id))
	if err != nil {
		return nil, fmt.Errorf("cannot read music from Redis: %s", err)
	}

	var music groovecoaster.MusicDetail
	if err := json.Unmarshal([]byte(musicString), &music); err != nil {
		return nil, fmt.Errorf("cannot unmarshal music: %s", err)
	}

	return &music, nil
}

func tweet(items []item) error {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	tweet := ""
	for _, item := range items {
		if len(item.diff) != 0 {
			text := item.String()
			if len(tweet+text) > 140 {
				_, err := api.PostTweet(strings.TrimSpace(tweet), nil)
				if err != nil {
					return err
				}

				tweet = text
			}
		}
	}

	if tweet != "" {
		_, err := api.PostTweet(strings.TrimSpace(tweet), nil)
		if err != nil {
			return err
		}
	}

	log.Println("Tweeted")

	return nil
}
func main() {
	client := groovecoaster.New()

	var err error
	log.Println("Connecting to Redis server...")
	c, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Redis server")
	defer func() {
		c.Close()
		log.Println("Redis Connection closed")
	}()

	lastDate, err := lastDateFromRedis()
	if err != nil {
		log.Print(err)
		return
	}

	log.Println("Fetching musics summary from mypage...")
	summary, err := client.MusicSummary()
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Fetching musics summary was successful")

	splittedSummary := splitMusicSummary(summary, lastDate)
	items := make([]item, len(splittedSummary))
	for _, music := range splittedSummary {
		oldMusic, err := musicFromRedis(strconv.Itoa(music.ID))
		if err != nil {
			log.Print(err)
			return
		}

		newMusic, err := client.Music(music.ID)
		if err != nil {
			log.Print(err)
			return
		}

		var item item
		item.name = music.Title

		if oldMusic == nil {
			oldMusic = newMusic
		}

		results := []*groovecoaster.Result{newMusic.Simple, newMusic.Normal, newMusic.Hard}
		if newMusic.HasEx {
			results = append(results, newMusic.Extra)
		}

		for i, diff := range results {
			if diff == nil {
				continue
			}

			var oldMusicDiff *groovecoaster.Result
			archived := []string{}

			switch groovecoaster.Difficulty(i) {
			case groovecoaster.Simple:
				oldMusicDiff = oldMusic.Simple
			case groovecoaster.Normal:
				oldMusicDiff = oldMusic.Normal
			case groovecoaster.Hard:
				oldMusicDiff = oldMusic.Hard
			case groovecoaster.Extra:
				if !newMusic.HasEx {
					continue
				}
				oldMusicDiff = oldMusic.Extra
			}

			// どれかにあてはまる場合
			switch {
			case diff.Perfect == 1:
				archived = append(archived, "Perf")
			case diff.FullChain == 1:
				archived = append(archived, "FC")
			case diff.NoMiss == 1:
				archived = append(archived, "NM")
			}

			if diff.MaxChain > oldMusicDiff.MaxChain {
				archived = append(archived, fmt.Sprintf("⛓ +%d", diff.MaxChain-oldMusicDiff.MaxChain))
			}

			if diff.PlayCount == 100 {
				archived = append(archived, "100 Played!")
			}

			if diff.Score > oldMusicDiff.Score {
				archived = append(archived, fmt.Sprintf("💯 +%d", diff.MaxChain-oldMusicDiff.MaxChain))
			}

			if len(archived) == 0 {
				continue
			}

			item.diff = append(item.diff, strings.Join(archived, ", "))
		}

		bytes, err := json.Marshal(newMusic)
		if err != nil {
			log.Print(err)
			return
		}

		c.Do("SET", music.ID, string(bytes))
		log.Printf("Updated music cache: %s", music.Title)

		items = append(items, item)
	}

	tweet(items)

	c.Do("SET", "lastDate", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("Updated lastDate")

	log.Println("Completed")
}

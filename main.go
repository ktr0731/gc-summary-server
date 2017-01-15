package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/lycoris0731/go-groovecoaster/groovecoaster"
)

var c redis.Conn

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

func splitnewMusicSummary(summary []*groovecoaster.MusicSummary, lastDate time.Time) []*groovecoaster.MusicSummary {
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

	text := ""

	for _, music := range splitnewMusicSummary(summary, lastDate) {
		old, err := musicFromRedis(strconv.Itoa(music.ID))
		if err != nil {
			log.Print(err)
			return
		}

		newMusic, err := client.Music(music.ID)
		if err != nil {
			log.Print(err)
			return
		}

		if old == nil {
			bytes, err := json.Marshal(newMusic)
			if err != nil {
				log.Print(err)
				return
			}
			c.Do("SET", music.ID, string(bytes))

			old, err = musicFromRedis(strconv.Itoa(music.ID))
			if err != nil {
				log.Print(err)
				return
			}
		}

		results := []*groovecoaster.Result{newMusic.Simple, newMusic.Normal, newMusic.Hard}
		if newMusic.HasEx {
			results = append(results, newMusic.Extra)
		}

		var summary string
		for i, diff := range results {
			var oldDiff *groovecoaster.Result

			if diff == nil {
				continue
			}

			var diffName string

			switch groovecoaster.Difficulty(i) {
			case groovecoaster.Simple:
				diffName = "  [S] "
				oldDiff = old.Simple
			case groovecoaster.Normal:
				diffName = "  [N] "
				oldDiff = old.Normal
			case groovecoaster.Hard:
				diffName = "  [H] "
				oldDiff = old.Hard
			case groovecoaster.Extra:
				if !newMusic.HasEx {
					continue
				}
				diffName = "  [E] "
				oldDiff = old.Extra
			}

			archived := []string{}

			// どれかにあてはまる場合
			switch {
			case diff.Perfect == 1:
				archived = append(archived, "Perf")
			case diff.FullChain == 1:
				archived = append(archived, "FC")
			case diff.NoMiss == 1:
				archived = append(archived, "NM")
			}

			if diff.MaxChain > oldDiff.MaxChain {
				archived = append(archived, fmt.Sprintf("⛓ +%d", diff.MaxChain-oldDiff.MaxChain))
			}

			if diff.PlayCount == 100 {
				archived = append(archived, "100 Played!")
			}

			if diff.Score > oldDiff.Score {
				archived = append(archived, fmt.Sprintf("💯 +%d", diff.MaxChain-oldDiff.MaxChain))
			}

			if len(archived) == 0 {
				continue
			}

			summary += diffName + strings.Join(archived, ", ") + "\n"
		}

		if summary != "" {
			text += music.Title + "\n" + summary + "\n"
		}

		bytes, err := json.Marshal(newMusic)
		if err != nil {
			log.Print(err)
			return
		}

		c.Do("SET", music.ID, string(bytes))
		log.Printf("Updated music cache: %s", music.Title)
	}

	log.Printf("Tweet:\n%s", text)

	c.Do("SET", "lastDate", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("Updated lastDate")

	log.Println("Completed")
}

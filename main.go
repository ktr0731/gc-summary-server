package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
			text += "  [Simple] "
		case groovecoaster.Normal:
			text += "  [Normal] "
		case groovecoaster.Hard:
			text += "  [Hard]   "
		case groovecoaster.Extra:
			text += "  [Extra]  "
		}

		text += v + "\n"
	}

	return text + "\n"
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		summary := generateSummary()
		if summary == "" {
			summary = "No change...\nWhy don't play GrooveCoaster?\n"
		}
		io.WriteString(w, summary)
	})

	log.Println("Listen in", ":"+os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func generateSummary() string {
	client := groovecoaster.New()

	var err error
	log.Println("Connecting to Redis server...")
	c, err = redis.DialURL(os.Getenv("REDISTOGO_URL"))
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
		log.Println(err)
		return ""
	}
	log.Println("last updated:", lastDate)

	log.Println("Fetching musics summary from mypage...")
	summary, err := client.MusicSummary()
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println("Fetching musics summary was successful")

	splittedSummary := splitMusicSummary(summary, lastDate)
	items := make([]item, len(splittedSummary))
	for _, music := range splittedSummary {
		oldMusic, err := musicFromRedis(strconv.Itoa(music.ID))
		if err != nil {
			log.Println(err)
			return ""
		}

		newMusic, err := client.Music(music.ID)
		if err != nil {
			log.Println(err)
			return ""
		}

		var item item
		item.name = music.Title

		if oldMusic.ID == "" {
			oldMusic = newMusic
		}

		results := []groovecoaster.Result{newMusic.Simple, newMusic.Normal, newMusic.Hard}
		if newMusic.HasEx {
			results = append(results, newMusic.Extra)
		}

		for i, diff := range results {
			if diff.PlayCount == 0 {
				continue
			}

			var oldMusicDiff groovecoaster.Result
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
				archived = append(archived, "Perfect")
			case diff.FullChain == 1:
				archived = append(archived, "FullChain")
			case diff.NoMiss == 1:
				archived = append(archived, "NoMiss")
			}

			if diff.MaxChain > oldMusicDiff.MaxChain {
				archived = append(archived, fmt.Sprintf(":chains: Max Chain +%d", diff.MaxChain-oldMusicDiff.MaxChain))
			}

			if diff.PlayCount == 100 {
				archived = append(archived, ":tada: Total played time is over 100")
			}

			if diff.Score > oldMusicDiff.Score {
				archived = append(archived, fmt.Sprintf(":chart_with_upwards_trend: Total Score +%d", diff.MaxChain-oldMusicDiff.MaxChain))
			}

			if len(archived) == 0 {
				continue
			}

			item.diff = append(item.diff, strings.Join(archived, ", "))
		}

		bytes, err := json.Marshal(newMusic)
		if err != nil {
			log.Println(err)
			return ""
		}

		c.Do("SET", music.ID, string(bytes))
		log.Printf("Updated music cache: %s", music.Title)

		items = append(items, item)
	}

	// text := fmt.Sprintf("GrooveCoaster Summary of %s\n(From %s)\n\n", "ktr", lastDate.String())
	text := ""
	for _, item := range items {
		if len(item.diff) == 0 {
			continue
		}
		text += item.String() + "\n"
	}
	log.Println(text)

	c.Do("SET", "lastDate", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("Updated lastDate")

	log.Println("Completed")

	return strings.TrimSpace(text + "\n")
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

func splitMusicSummary(summary []groovecoaster.MusicSummary, lastDate time.Time) []groovecoaster.MusicSummary {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	for i, m := range summary {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", m.LastPlayTime, loc)
		if err != nil {
			log.Println(err)
			return nil
		}

		if date.Equal(lastDate) || lastDate.After(date) {
			return summary[:i]
		}
	}

	return []groovecoaster.MusicSummary{}
}

func musicFromRedis(id string) (groovecoaster.MusicDetail, error) {
	exists, err := redis.Bool(c.Do("EXISTS", id))
	if err != nil {
		return groovecoaster.MusicDetail{}, fmt.Errorf("cannot confirm whether music is exists from Redis: %s", err)
	}
	if !exists {
		return groovecoaster.MusicDetail{}, nil
	}

	musicString, err := redis.String(c.Do("GET", id))
	if err != nil {
		return groovecoaster.MusicDetail{}, fmt.Errorf("cannot read music from Redis: %s", err)
	}

	var music groovecoaster.MusicDetail
	if err := json.Unmarshal([]byte(musicString), &music); err != nil {
		return groovecoaster.MusicDetail{}, fmt.Errorf("cannot unmarshal music: %s", err)
	}

	return music, nil
}

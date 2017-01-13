package main

import (
	"log"
	"time"

	"github.com/lycoris0731/go-groovecoaster/groovecoaster"
)

func main() {
	client := groovecoaster.New()

	summary, err := client.MusicSummary()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Redisから取得
	lastDate, err := time.Parse("2006-01-02 15:04:05", "2017-01-01 10:21:30")
	if err != nil {
		log.Fatal(err)
	}

	for _, music := range summary {
		musicDate, err := time.Parse("2006-01-02 15:04:05", music.LastPlayTime)
		if err != nil {
			log.Fatal(err)
		}

		if musicDate.Equal(lastDate) || musicDate.After(lastDate) {
			log.Println("finish")
			break
		}

		// TODO: Redisから取得
		old := client.Music(music.ID)

		detail, err := client.Music(music.ID)
		if err != nil {
			log.Fatal(err)
		}

		/**
		 * さんぷる
		 *
		 * Got a pain cover?
		 *  [Simple] 初Perfect, +2011
		 *  [Hard] 初FullChain, +20300
		 *
		 * Axeria
		 *  [Normal] 初NoMiss
		 *
		 * #gcsummary bit.ly/hoge <- gc-summary-botのURL
		 */

		if old.Simple.Score < detail.Simple.Score {
			// TODO: 作成
		}
	}
}

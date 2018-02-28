package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type messageArray []*discordgo.Message

type messageSet struct {
	Items messageArray
	Set   map[string]bool
}

var s *discordgo.Session
var pinmap = map[string]*messageSet{} //Map of CHANNEL ID to STRUCT with MAP OF MESSAGE ID TO BOOLEAN and SORTED ARRAY OF MESSAGES
var guild *discordgo.Guild

func discordStart() {
	var err error
	s, err = discordgo.New("Bot " + os.Getenv("TOKEN"))
	p(err)

	p(s.Open())

	c := make(chan struct{})

	s.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		c <- struct{}{}
	})

	log.Println("Waiting for ready...")

	<-c
	log.Println("Ready recieved")
	guild, err = s.Guild(os.Getenv("SERVER"))
	p(err)

	file, err := ioutil.ReadFile(os.Getenv("HOME")+"/store.json")
	if err == nil {
		p(json.Unmarshal(file, &pinmap))
	}

	time.Sleep(4*time.Second)
	discordCheckAll(guild.Channels, 0)
	tick := time.NewTicker(time.Hour)
	go func() {
		for range tick.C {
			log.Println("Checking at ", time.Now().Format(time.Stamp))
			discordCheckAll(guild.Channels, 3*time.Second)
		}
	}()
}

func discordCheck(ids ...string) {
	for _, id := range ids {
		if _, ok := pinmap[id]; !ok {
			pinmap[id] = &messageSet{Items: make([]*discordgo.Message, 0), Set: make(map[string]bool)}
		}
		pins, err := s.ChannelMessagesPinned(id)
		check(err)
		for _, x := range pins {
			if _, ok := pinmap[id].Set[x.ID]; !ok {
				pinmap[id].Set[x.ID] = true
				pinmap[id].Items = append(pinmap[id].Items, x)
				sort.Sort(sort.Reverse(pinmap[id].Items))
			}
		}
	}
	data, err := json.Marshal(pinmap)
	check(err)
	check(ioutil.WriteFile(os.Getenv("HOME")+"/store.json", data, 0666))
}

func discordCheckAll(channels []*discordgo.Channel, d time.Duration) {
	for _, channel := range channels {
		discordCheck(channel.ID)
		time.Sleep(d)
	}
}

func (by messageArray) Len() int {
	return len(by)
}

func (by messageArray) Less(i, j int) bool {
	i2, err := strconv.Atoi(by[i].ID)
	p(err)
	j2, err := strconv.Atoi(by[j].ID)
	p(err)
	return i2 < j2
}

func (by messageArray) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}

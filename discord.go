package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"encoding/json"
	"time"
	"strconv"
	"sort"
	"log"
	"io/ioutil"
)

type messageArray []*discordgo.Message

type messageSet struct{
	Items messageArray
	Set map[string]bool
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
	guild, err = s.State.Guild(os.Getenv("SERVER"))
	p(err)

	file, err := ioutil.ReadFile("store.json")
	if err == nil {
		p(json.Unmarshal(file, &pinmap))
	}

	tick := time.NewTicker(time.Hour)
	defer func() {tick.Stop()}()
	go func() {
		for range tick.C {
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
		p(err)
		for _, x := range pins {
			if _, ok := pinmap[id].Set[x.ID]; !ok {
				pinmap[id].Set[x.ID] = true
				pinmap[id].Items = append(pinmap[id].Items, x)
				sort.Sort(sort.Reverse(pinmap[id].Items))
			}
		}
	}
	data, err := json.Marshal(pinmap)
	p(err)
	check(ioutil.WriteFile("store.json", data, 0666))
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

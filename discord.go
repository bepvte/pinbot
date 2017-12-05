package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"encoding/json"
	"time"
	"strconv"
	"sort"
	"log"
)

type messageArray []*discordgo.Message

type messageSet struct{
	Items messageArray
	Set map[string]bool
}

var s *discordgo.Session
var pinmap = map[string]*messageSet{} //Map of CHANNEL ID to STRUCT with MAP OF MESSAGE ID TO BOOLEAN and SORTED ARRAY OF MESSAGES
var guild *discordgo.Guild
var checkTimes = map[string]time.Time{}


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

	reply, err := db.Query("SELECT myblob FROM mydata WHERE myname = 'store'")
	p(err)
	if reply.Next() {
		var resp []byte
		err := reply.Scan(&resp)
		p(err)
		json.Unmarshal(resp, &pinmap)
	}
	time.Sleep(5*time.Second)
}

func discordCheck(ids ...string) {
	for _, id := range ids {
		checkTimes[id] = time.Now()
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
	_, err = db.Exec("INSERT INTO mydata VALUES ('store', $1) ON CONFLICT (myname) DO UPDATE SET myblob = $1", data)
	check(err)
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
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack"

	"github.com/bwmarrin/discordgo"
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
	s.LogLevel = discordgo.LogInformational

	p(err)

	c := make(chan struct{})

	s.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		c <- struct{}{}
	})
	p(s.Open())

	log.Println("Waiting for ready...")

	<-c
	log.Println("Ready recieved")
	guild, err = s.Guild(os.Getenv("SERVER"))
	p(err)

	file, err := os.OpenFile(os.Getenv("HOME")+"/pinbot.msgp.xz", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("failed to open pinbot.msgp")
	} else {
		cmd := exec.Command("xz", "-c")
		cmd.Stdin = file
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		if err := msgpack.NewDecoder(stdout).UseJSONTag(true).Decode(&pinmap); err != io.EOF && err != nil {
			panic(err)
		}
		stdout.Close()
		cmd.Wait()
	}

	time.Sleep(4 * time.Second)
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
	file, err := os.OpenFile(os.Getenv("HOME")+"/pinbot.msgp.xz", os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("xz", "-c")
	cmd.Stdout = file
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	if err := msgpack.NewEncoder(stdin).UseJSONTag(true).Encode(&pinmap); err != nil {
		panic(err)
	}
	stdin.Close()
	cmd.Wait()
	fmt.Println("DEBUG done")
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

package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"time"
)

type channelPage struct {
	Current  *discordgo.Channel
	Channels []*discordgo.Channel
	Pins     messageArray
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "channel.html", channelPage{Current: guild.Channels[0], Channels: guild.Channels, Pins: pinmap[guild.Channels[0].ID].Items})
	if check(err) {
		failed(w, err)
		return
	}
}

func channelHandler(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "ID")
	channel, err := s.State.Channel(channelID)
	if check(err) {
		failed(w, err)
	}
	err = t.ExecuteTemplate(w, "channel.html", channelPage{Current: channel, Channels: guild.Channels, Pins: pinmap[channel.ID].Items})
	if check(err) {
		failed(w, err)
	}
}

func reloadhandler(w http.ResponseWriter, r *http.Request) {
	discordCheckAll(guild.Channels, time.Second*0)
	io.WriteString(w, "Its refreshing! (as far as i can tell)")
}

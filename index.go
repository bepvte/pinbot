package main

import (
	"net/http"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"time"
)

type channelPage struct{
	Current *discordgo.Channel
	Channels []*discordgo.Channel
	Pins messageArray
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
	if checkTimes[channelID].IsZero() && time.Now().Sub(checkTimes[channelID]) >= time.Hour {
		discordCheck(channelID)
	}
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
	checkTimes = make(map[string]time.Time)
}
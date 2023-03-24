package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v6/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) FilterPost(post *model.Post) {
	conf := p.getConfiguration()
	_, fromBot := post.GetProps()["from_bot"]
	if fromBot || post.Message == conf.Response {
		return
	}

	contains := false
	for _, word := range conf.Keywords {
		if strings.Contains(post.Message, word) {
			contains = true
		}
	}

	if contains {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			ChannelId: post.ChannelId,
			RootId:    post.Id,
			Message:   conf.Response,
		})
	}
}

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	p.FilterPost(post)
}

// List of multistream service providers and their known servers by region

package config

import (
	"fmt"
)

// msServer holds information on each multi-streaming server
type msServer struct {
	name string
	url  string
}

// msService holds information on each multistreaming service and known servers
type msService struct {
	name    string
	servers map[string]msServer
}

// RTMPUrl builds URL for a service given a key and server. If `server` is not
// specified, `default` server is used.
func (m msService) RTMPUrl(key string, server string) string {
	var u string
	if server == "" {
		server = "default"
	}
	if s, ok := m.servers[server]; ok {
		u = fmt.Sprintf("%s/%s", s.url, key)
	}
	return u
}

// HasServer returns whether a server key is present for given msService
func (m msService) HasServer(key string) bool {
	_, ok := m.servers[key]
	return ok
}

// msServiceYouTube stores YouTube's server configurations for multistreaming
var msServiceYouTube = msService{
	name: "YouTube / YouTube Gaming",
	servers: map[string]msServer{
		"default": msServer{
			name: "YouTube default RTMP ingestion server",
			url:  "rtmp://a.rtmp.youtube.com/live2",
		},
		"backup": msServer{
			name: "YouTube backup RTMP ingestion server",
			url:  "rtmp://b.rtmp.youtube.com/live2?backup=1",
		},
	},
}

// msServiceTwitch stores Twitch's server configurations for multistreaming
var msServiceTwitch = msService{
	name: "Twitch.TV",
	servers: map[string]msServer{
		"default": msServer{
			name: "Twitch default RTMP ingestion server",
			url:  "rtmp://live.twitch.tv/app",
		},
		"sin": msServer{
			name: "Asia: Singapore",
			url:  "rtmp://live-sin.twitch.tv/app",
		},
		"ams": msServer{
			name: "EU: Amsterdam, NL",
			url:  "rtmp://live-ams.twitch.tv/app",
		},
		"lon": msServer{
			name: "EU: London, UK",
			url:  "rtmp://live-lhr.twitch.tv/app",
		},
		"nyc": msServer{
			name: "US East: New York, NY",
			url:  "rtmp://live-jfk.twitch.tv/app",
		},
		"sfo": msServer{
			name: "US West: San Francisco, CA",
			url:  "rtmp://live-sfo.twitch.tv/app",
		},
	},
}

// MSServices is a map of multistreaming services
var MSServices = map[string]*msService{
	"YouTube": &msServiceYouTube,
	"Twitch":  &msServiceTwitch,
}

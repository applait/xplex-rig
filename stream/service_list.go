package stream

// List of multistream service providers and their known servers by region
import (
	"fmt"
)

// IngestServer holds information on each multi-streaming server
type IngestServer struct {
	name string
	url  string
}

// Service holds information on each multistreaming service and known servers
type Service struct {
	Name    string
	servers map[string]IngestServer
}

// RTMPUrl builds URL for a service given a key and server. If `server` is not
// specified, `default` server is used.
func (m Service) RTMPUrl(key string, server string) string {
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
func (m Service) HasServer(key string) bool {
	_, ok := m.servers[key]
	return ok
}

// ServiceYouTube stores YouTube's server configurations for multistreaming
var ServiceYouTube = Service{
	Name: "YouTube / YouTube Gaming",
	servers: map[string]IngestServer{
		"default": IngestServer{
			name: "YouTube default RTMP ingestion server",
			url:  "rtmp://a.rtmp.youtube.com/live2",
		},
		"backup": IngestServer{
			name: "YouTube backup RTMP ingestion server",
			url:  "rtmp://b.rtmp.youtube.com/live2?backup=1",
		},
	},
}

// ServiceTwitch stores Twitch's server configurations for multistreaming
var ServiceTwitch = Service{
	Name: "Twitch.TV",
	servers: map[string]IngestServer{
		"default": IngestServer{
			name: "Twitch default RTMP ingestion server",
			url:  "rtmp://live.twitch.tv/app",
		},
		"sin": IngestServer{
			name: "Asia: Singapore",
			url:  "rtmp://live-sin.twitch.tv/app",
		},
		"ams": IngestServer{
			name: "EU: Amsterdam, NL",
			url:  "rtmp://live-ams.twitch.tv/app",
		},
		"lon": IngestServer{
			name: "EU: London, UK",
			url:  "rtmp://live-lhr.twitch.tv/app",
		},
		"nyc": IngestServer{
			name: "US East: New York, NY",
			url:  "rtmp://live-jfk.twitch.tv/app",
		},
		"sfo": IngestServer{
			name: "US West: San Francisco, CA",
			url:  "rtmp://live-sfo.twitch.tv/app",
		},
	},
}

// Services is a map of multistreaming services
var Services = map[string]*Service{
	"YouTube": &ServiceYouTube,
	"Twitch":  &ServiceTwitch,
}

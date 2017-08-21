/**
 * Definition of services available for multistream
 */
const MultiStreamServices = {

  // YouTube Live / Youtube Gaming configuration
  'YouTube': {
    name: 'YouTube / YouTube Gaming',
    servers: new Map([
      ['default', {
        name: 'YouTube default RTMP ingestion server',
        url: 'rtmp://a.rtmp.youtube.com/live2'
      }],
      ['backup', {
        name: 'YouTube backup RTMP ingestion server',
        url: 'rtmp://b.rtmp.youtube.com/live2?backup=1'
      }]
    ])
  },

  // Twitch configuration
  'Twitch': {
    name: 'Twitch.TV',
    servers: new Map([
      ['default', {
        name: 'Twitch default RTMP ingestion server',
        url: 'rtmp://live.twitch.tv/app'
      }],
      ['SIN', {
        name: 'Asia: Singapore',
        url: 'rtmp://live-sin.twitch.tv/app'
      }],
      ['AMS', {
        name: 'EU: Amsterdam, NL',
        url: 'rtmp://live-ams.twitch.tv/app'
      }],
      ['LON', {
        name: 'EU: London, UK',
        url: 'rtmp://live-lhr.twitch.tv/app'
      }],
      ['NYC', {
        name: 'US East: New York, NY',
        url: 'rtmp://live-jfk.twitch.tv/app'
      }],
      ['SFO', {
        name: 'US West: San Francisco, CA',
        url: 'rtmp://live-sfo.twitch.tv/app'
      }]
    ])
  }
}

module.exports = MultiStreamServices

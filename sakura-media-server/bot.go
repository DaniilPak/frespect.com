package sakura

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/oggreader"
)

type Bot struct {
	botID     string
	track     *webrtc.TrackLocalStaticSample
	rtpSender *webrtc.RTPSender
	room      *Room
}

var mediaManagerURL string = GetMediaManagerURL()
var audioURL = mediaManagerURL + "media/ef9DqYNBjUC"

func NewBot() (*Bot, error) {
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	return &Bot{botID: id}, nil
}

func (b *Bot) SetRoom(room *Room) {
	b.room = room
}

func (b *Bot) SetTrack(track *webrtc.TrackLocalStaticSample) {
	b.track = track
}

func (b *Bot) SetRTPSender(rtpSender *webrtc.RTPSender) {
	b.rtpSender = rtpSender
}

func (b *Bot) CreateAudioTrack() *webrtc.TrackLocalStaticSample {
	musicTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio-music", "pion-music")
	if err != nil {
		panic(err)
	}
	return musicTrack
}

func (b *Bot) WriteAudioToTrack(audioURL string, musicTrack *webrtc.TrackLocalStaticSample) {
	// Fetch the audio file from HTTP API
	resp, err := http.Get(audioURL)
	if err != nil {
		fmt.Println("Failed to fetch audio file.")
		return
	}
	defer resp.Body.Close()

	ogg, _, err := oggreader.NewWith(resp.Body)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(oggPageDuration)
	defer ticker.Stop()

	var lastGranule uint64
	for range ticker.C {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if errors.Is(err, io.EOF) {
			fmt.Println("[1]: All audio pages parsed.")
			b.BotStop()
			break
		}
		if err != nil {
			fmt.Println("[2]: All audio pages parsed.")
			b.BotStop()
		}

		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		if writeErr := musicTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); writeErr != nil {
			fmt.Println("[3]: All audio pages parsed.")
			b.BotStop()
		}
	}
}

func (b *Bot) BotStop() {
	roomManager := GetRoomManager()
	roomManager.RemoveTrackFromRoom(b.room, b.botID)
	roomManager.RenegotiateAllClients()
}

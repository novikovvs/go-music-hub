package main

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"os"
	"strconv"
	"time"
)

var Streamer beep.StreamSeekCloser
var Format beep.Format
var playingDone chan bool
var AP *audioPanel
var MusicQueue *TrackQueue
var TrackEnd chan bool
var CurrentTrack Track

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	seq        beep.Streamer
}

func (ap *audioPanel) play() {
	defer func() {
		TrackEnd <- true
		PlayerLogger.Info("Track end")
	}()
	PlayerLogger.Info("Playing audio")
	speaker.Play(ap.seq)
	<-playingDone
	PlayerLogger.Info("Stop playing audio")
}

func Play() {
	if AP != nil {
		return
	}

	TrackEnd = make(chan bool, 1)
	playingDone = make(chan bool, 1)
	go endTrackHandler()

	if MusicQueue == nil {
		MusicQueue = newMusicQueue()
	}

	if MusicQueue.size <= 0 {
		return
	}

	track, _ := MusicQueue.Dequeue()
	setCurrentTrack(track)

	PlayerLogger.Info(track.Path)

	var err error
	Streamer, Format, err = loadFile(track.Path)
	if err != nil {
		panic(err)
	}
	defer func(Streamer beep.StreamSeekCloser) {
		err := Streamer.Close()
		if err != nil {
			panic(err)
		}
	}(Streamer)

	err = speaker.Init(Format.SampleRate, Format.SampleRate.N(time.Second/30))
	if err != nil {
		panic(err)
	}

	AP = newAudioPanel(Format.SampleRate, Streamer)

	AP.play()
}

func Pause() {
	if AP == nil {
		return
	}

	speaker.Lock()
	AP.ctrl.Paused = !AP.ctrl.Paused
	PlayerLogger.Info("Paused: " + strconv.FormatBool(AP.ctrl.Paused))

	speaker.Unlock()
}

func Next() {
	if AP == nil {
		return
	}

	speaker.Lock()
	track, _ := MusicQueue.Dequeue()
	setCurrentTrack(track)
	stream, format, _ := loadFile(track.Path)
	AP.resampler.SetRatio(float64(format.SampleRate) / float64(AP.sampleRate))
	AP.ctrl.Streamer = stream
	PlayerLogger.Info("Next track: " + track.Path)
	speaker.Unlock()
}

func GetCurrentTrack() Track {
	return CurrentTrack
}

func loadFile(filePath string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(filePath)
	if err != nil {
		PlayerLogger.Fatal(err.Error())
		panic(err)
	}

	return mp3.Decode(f)
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resample := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resample, Base: 2}
	sequence := beep.Seq(volume, beep.Callback(func() { playingDone <- true }))
	return &audioPanel{sampleRate, streamer, ctrl, resample, volume, sequence}
}

func endTrackHandler() {
	for {
		select {
		case <-TrackEnd:
			{
				PlayerLogger.Info("Auto Next()")
				Next()
				go AP.play()
			}
		}
	}
}

func setCurrentTrack(track Track) {
	CurrentTrack = track
	PlayerLogger.Info("Current playing track: " + track.Label)
}

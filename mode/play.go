package mode

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audioutil"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/render"
)

// mode.ScenePlay is template struct. Fields can be set at outer function call.
// Update cannot be generalized; each scene use template fields in timely manner.
// Unit of time is a millisecond (1ms = 0.001s).
type ScenePlay struct {
	// General
	Tick    int
	MD5     [md5.Size]byte // MD5 for raw chart file.
	EndTime int64          // EndTime = Duration + WaitAfter
	// General: Graphics
	BackgroundDrawer BackgroundDrawer

	// Speed, BPM, Volume and Highlight
	MainBPM   float64
	SpeedBase float64
	*TransPoint

	// Audio
	LastVolume  float64
	Volume      float64
	MusicPlayer *audio.Player
	MusicCloser func() error
	Sounds      audioutil.SoundMap // A player for sample sound is generated at a place.

	// Input
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool

	// Note
	// Note: Graphics
	BarLineDrawer BarLineDrawer

	// Score
	Result
	NoteWeights float64
	Combo       int
	Flow        float64
	// Score: Graphics
	ScoreDrawer ScoreDrawer
	TimingMeter TimingMeter
}

// General
const (
	DefaultWaitBefore int64 = int64(-1.8 * 1000)
	DefaultWaitAfter  int64 = 3 * 1000
)

func (s ScenePlay) BeatRatio() float64 { return s.TransPoint.BPM / s.MainBPM }
func (s ScenePlay) Speed() float64     { return s.SpeedBase * s.BeatRatio() * s.BeatScale }

func MD5(path string) [md5.Size]byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return md5.Sum(b)
}

func TimeToTick(time int64) int { return int(float64(time) / 1000 * float64(MaxTPS)) }
func TickToTime(tick int) int64 { return int64(float64(tick) / float64(MaxTPS) * 1000) }
func (s ScenePlay) Time() int64 { return int64(float64(s.Tick) / float64(MaxTPS) * 1000) }

// SetTick returns the duration of waiting for the music / chart starts.
func (s *ScenePlay) SetTick(rf *osr.Format) int64 {
	waitBefore := DefaultWaitBefore
	if rf != nil && rf.BufferTime() < waitBefore {
		waitBefore = rf.BufferTime()
	}
	s.Tick = TimeToTick(waitBefore) - 1 // s.Update() starts with s.Tick++
	return waitBefore
}

// General: Graphics
func (s ScenePlay) SetWindowTitle(c Chart) {
	title := fmt.Sprintf("gosu - %s - [%s]", c.MusicName, c.ChartName)
	ebiten.SetWindowTitle(title)
}
func (s *ScenePlay) SetBackground(path string) {
	if img := render.NewImage(path); img != nil {
		sprite := render.Sprite{
			I:      img,
			Filter: ebiten.FilterLinear,
		}
		sprite.SetWidth(screenSizeX)
		sprite.SetCenterY(screenSizeY / 2)
		s.BackgroundDrawer.Sprite = sprite
	} else {
		s.BackgroundDrawer.Sprite = DefaultBackground
	}
}

// Speed, BPM, Volume and Highlight
func (s *ScenePlay) SetInitTransPoint(first *TransPoint) {
	s.TransPoint = first
	for s.TransPoint.Time == s.TransPoint.Next.Time {
		s.TransPoint = s.TransPoint.Next
	}
	s.Volume = Volume * s.TransPoint.Volume
}
func (s *ScenePlay) UpdateTransPoint() {
	for s.TransPoint.Next != nil && s.TransPoint.Next.Time <= s.Time() {
		s.TransPoint = s.TransPoint.Next
	}
	if s.LastVolume != s.TransPoint.Volume {
		s.LastVolume = s.Volume
		s.Volume = Volume * s.TransPoint.Volume
		// s.MusicPlayer.SetVolume(s.Volume)
	}
}

// Audio
func (s *ScenePlay) SetMusicPlayer(apath string) error { // apath stands for audio path.
	if apath == "virtual" || apath == "" {
		return nil
	}
	var err error
	s.MusicPlayer, s.MusicCloser, err = audioutil.NewPlayer(apath)
	if err != nil {
		return err
	}
	s.MusicPlayer.SetVolume(s.Volume)
	return nil
}
func (s *ScenePlay) SetSoundMap(cpath string, names []string) error {
	s.Sounds = audioutil.NewSoundMap(&Volume)
	for _, name := range names {
		path := filepath.Join(filepath.Dir(cpath), name)
		s.Sounds.Register(path)
	}
	return nil
}

// Input
func (s ScenePlay) KeyAction(k int) input.KeyAction {
	return input.CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}

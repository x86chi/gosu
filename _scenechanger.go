package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode"
)
// todo: game에 그대로 구현
// the word changer is more descriptive than manager
type SceneChanger struct {
	g              *Game
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
}

// The function is called every time when settings has been updated
func (g *Game) NewSceneChanger() *SceneChanger {
	sc := &SceneChanger{}
	sc.g = g
	sc.Scene = g.Scene
	p := mode.ScreenSize()
	sc.TransSceneFrom, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	sc.TransSceneTo, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	return sc
}

func (sc *SceneChanger) Update() error {
	if sc.done() {
		return nil
		// return sc.g.Scene.Update()
	}
	sc.TransCountdown--
	if sc.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	sc.g.Scene = sc.NextScene
	sc.g.Scene.Init()
	sc.NextScene = nil
	return nil
}

func (sc *SceneChanger) Draw(screen *ebiten.Image) {
	if sc.done() {
		return
	}
	var value float64
	{
		value = float64(sc.TransCountdown) / float64(sc.MaxTransCountDown())
		sc.TransSceneFrom.Clear()
		sc.Scene.Draw(sc.TransSceneFrom)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(sc.TransSceneFrom, &op)
	}
	{
		value = 1 - value
		sc.TransSceneTo.Clear()
		sc.NextScene.Draw(sc.TransSceneTo)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(sc.TransSceneTo, &op)
	}
}

func (sc *SceneChanger) changeScene(s Scene) {
	sc.NextScene = s
	sc.TransCountdown = sc.MaxTransCountDown()
}

// 모든 time 관련 단위는 ms
func (sc *SceneChanger) MaxTransCountDown() int { return mode.MaxTPS() * 4 / 5 }

func (sc *SceneChanger) done() bool { return sc.TransCountdown == 0 }
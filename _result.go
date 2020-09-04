package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
)

// 리절트창 따로 만들지 말고 mania에서 해결보기
// 버튼 만드는거 최소화
type SceneResult struct {
	ScreenResult *ebiten.Image
	Buttons      []ebitenui.Button
}

func (g *Game) NewSceneResult(sp ScenePlay) *SceneResult {
	sr := &SceneResult{}
	// 리절트 이미지 render
	// 리트라이 -> ScenePlay 에 있는 차트 with 모드 다시 실행
	// 리플레이 -> ScenePlay 에 있는 리플레이 with 모드 실행
	return sr
}

// select를 매번 새로 New 해야할까?
func (s *SceneResult) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.ChangeScene(NewSceneSelect())
	}
	for _, b := range s.Buttons {
		b.Update()
	}
	return nil
}

func (s *SceneResult) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.ScreenResult, &ebiten.DrawImageOptions{})
	for _, b := range s.Buttons {
		b.Draw(screen)
	}
}
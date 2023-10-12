package coldBrew

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SceneFace interface {
	New(*Manager) *Scene
	Index() string
}

type Scene struct {
	Manager *Manager
	World   donburi.World
	Systems []System
}

func NewScene(m *Manager) *Scene {

	return &Scene{
		Manager: m,
		World:   donburi.NewWorld(),
	}
}

func (s *Scene) AddEntity(components ...donburi.IComponentType) *donburi.Entry {
	entity := s.World.Entry(s.World.Create(components...))

	return entity
}

func (s *Scene) AddSystem(newSystem System) {
	s.Systems = append(s.Systems, newSystem)
}

func (s *Scene) Sync() {
	s.processSystems(ClientType, nil, nil)
}

func (s *Scene) Sim(dt float64) {
	s.processSystems(SimType, dt, nil)
}

func (s *Scene) Draw(screen *ebiten.Image) {

	s.processSystems(RenderType, 0, screen)
}

func (s *Scene) processSystems(sysType SystemType, args ...interface{}) {

	for _, system := range s.Systems {

		query := &donburi.Query{}
		customIteration := false

		if q, ok := system.(Query); ok {
			query = q.Query()

		} else {
			query = nil
		}

		if sysWithCustomIterMethod, ok := system.(CustomIteration); ok {
			customIteration = sysWithCustomIterMethod.CustomIteration()
		}

		runOnce := query == nil || customIteration

		if sysType == ClientType {

			if clientSys, ok := system.(Client); ok {

				if runOnce {
					clientSys.Sync(nil)
					continue
				}

				query.Each(s.World, func(entry *donburi.Entry) {
					clientSys.Sync(entry)
				})
				continue
			}
		}

		if sysType == SimType {

			dt := args[0].(float64)

			if simSys, ok := system.(Sim); ok {

				if runOnce {
					simSys.Run(dt, nil)
					continue
				}

				query.Each(s.World, func(entry *donburi.Entry) {
					simSys.Run(dt, entry)
				})

				continue
			}
		}

		if sysType == RenderType {

			screen := args[1].(*ebiten.Image)

			if renderSys, ok := system.(Render); ok {

				if runOnce {
					renderSys.Draw(screen, nil)
					continue
				}

				query.Each(s.World, func(entry *donburi.Entry) {
					renderSys.Draw(screen, entry)
				})
				continue
			}
		}
	}
}

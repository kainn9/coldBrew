package coldBrew

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SceneFace interface {
	New(*Manager) *Scene
	Index() string
}

type Scene struct {
	Height, Width int
	Manager       *Manager
	World         donburi.World
	Systems       []System

	mu        *sync.Mutex
	isLoading bool
	isLoaded  bool
}

func NewScene(m *Manager, width, height int) *Scene {

	return &Scene{
		Manager: m,
		World:   donburi.NewWorld(),
		Height:  height,
		Width:   width,
		mu:      &sync.Mutex{},
	}
}

func (s *Scene) SetIsLoading(loading bool) {
	s.mu.Lock()
	s.isLoading = loading
	s.mu.Unlock()
}

func (s *Scene) SetIsLoaded(loaded bool) {
	s.mu.Lock()
	s.isLoaded = loaded
	s.mu.Unlock()
}

func (s *Scene) IsCurrentlyLoading() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isLoading
}

func (s *Scene) IsCurrentlyLoaded() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isLoaded
}

func (s *Scene) AddEntity(components ...donburi.IComponentType) *donburi.Entry {
	entity := s.World.Entry(s.World.Create(components...))

	return entity
}

func (s *Scene) AddSystem(newSystem System) {
	s.Systems = append(s.Systems, newSystem)
}

func (s *Scene) Sync() {

	if !s.IsCurrentlyLoaded() {
		return
	}

	s.processSystems(ClientType, nil, nil)
}

func (s *Scene) Sim(dt float64) {

	if !s.IsCurrentlyLoaded() {
		return
	}

	s.processSystems(SimType, dt, nil)
}

func (s *Scene) Draw(screen *ebiten.Image) {

	if !s.IsCurrentlyLoaded() {

		loaderImage := s.Manager.LoaderImage
		screen.DrawImage(loaderImage, &ebiten.DrawImageOptions{})

		return
	}

	s.processSystems(RenderType, nil, screen)
}

func (s *Scene) Load() {

	if s.IsCurrentlyLoaded() || s.IsCurrentlyLoading() {
		return
	}

	s.SetIsLoading(true)
	s.processSystems(LoaderType, nil, nil)
	s.SetIsLoading(false)
	s.SetIsLoaded(true)
}

func (s *Scene) processSystems(sysType SystemType, args ...interface{}) {

	for _, system := range s.Systems {

		query := &donburi.Query{}

		if q, ok := system.(Query); ok {
			query = q.Query()

		} else {
			query = nil
		}

		if sysType == LoaderType {

			if loaderSys, ok := system.(Load); ok {

				if query == nil {
					loaderSys.Load(nil)
					continue
				}

				query.Each(s.World, func(entity *donburi.Entry) {
					loaderSys.Load(entity)
				})

				continue
			}
		}

		if sysType == ClientType {

			if clientSys, ok := system.(Client); ok {

				if query == nil {
					clientSys.Sync(nil)
					continue
				}

				query.Each(s.World, func(entity *donburi.Entry) {
					clientSys.Sync(entity)
				})
				continue
			}
		}

		if sysType == SimType {

			dt := args[0].(float64)

			if simSys, ok := system.(Sim); ok {

				if query == nil {
					simSys.Run(dt, nil)
					continue
				}

				query.Each(s.World, func(entity *donburi.Entry) {
					simSys.Run(dt, entity)
				})

				continue
			}
		}

		if sysType == RenderType {

			screen := args[1].(*ebiten.Image)

			if renderSys, ok := system.(Render); ok {

				if query == nil {
					renderSys.Draw(screen, nil)
					continue
				}

				query.Each(s.World, func(entity *donburi.Entry) {
					renderSys.Draw(screen, entity)
				})
				continue
			}
		}
	}
}

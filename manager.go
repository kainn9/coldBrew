package coldBrew

import "github.com/hajimehoshi/ebiten/v2"

type Manager struct {
	activeScene *Scene
	sceneCache  *sceneCache
	TickHandler *TickHandler
	LoaderImage *ebiten.Image
}

type ManagerError struct {
	msg string
}

func (e ManagerError) Error() string {
	return e.msg
}

func NewManager(cacheLimit, maxTick int, loaderImage *ebiten.Image) *Manager {

	return &Manager{
		sceneCache:  newSceneCache(cacheLimit),
		TickHandler: NewTickHandler(maxTick),
		LoaderImage: loaderImage,
	}
}

func (m *Manager) LoadScene(s SceneFace) error {
	key := s.Index()

	if ok, scene := m.sceneCache.check(key); ok {

		if scene == m.activeScene {
			return &ManagerError{msg: "attempting to load scene that is already active."}
		}

		// TODO:
		// Re-Index the cache buffer, so that the current scene is at the end of the buffer.
		// This is because we want to keep N most recently accessed scenes, not N most recently
		// loaded scenes. The distinction being that we may only load a scene once, but we may
		// access it many times.
		m.activeScene = scene

	} else {
		m.activeScene = s.New(m)
		m.sceneCache.add(key, m.activeScene)
	}

	m.activeScene.LastActiveTick = m.TickHandler.CurrentTick()

	return nil
}

func (m *Manager) ActiveScene() *Scene {
	return m.activeScene
}

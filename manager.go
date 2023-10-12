package coldBrew

type Manager struct {
	activeScene *Scene
	sceneCache  *sceneCache
}

type ManagerError struct {
	msg string
}

func (e ManagerError) Error() string {
	return e.msg
}

func NewManager(cacheLimit int) *Manager {

	return &Manager{
		sceneCache: newSceneCache(cacheLimit),
	}
}

func (m *Manager) LoadScene(s SceneFace) error {
	key := s.Index()

	if ok, sceneAdmin := m.sceneCache.check(key); ok {

		if sceneAdmin == m.activeScene {
			return &ManagerError{msg: "attempting to load scene that is already active."}
		}
		m.sceneCache.add(key, m.activeScene)
		m.activeScene = sceneAdmin

	} else {

		m.activeScene = s.New(m)
	}

	return nil
}

func (m *Manager) GetActiveScene() *Scene {
	return m.activeScene
}

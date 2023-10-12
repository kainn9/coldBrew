package coldBrew

// Ring buffer like cache for scenes.
type sceneCache struct {
	cachedScenes    map[string]*Scene
	keys            []string
	currentKeyIndex int
	limit           int
}

func newSceneCache(limit int) *sceneCache {
	return &sceneCache{
		cachedScenes: make(map[string]*Scene),
		keys:         make([]string, 0),
		limit:        limit,
	}
}

func (sc *sceneCache) add(key string, scene *Scene) {

	if len(sc.keys) < sc.limit {
		sc.keys = append(sc.keys, key)
	} else {
		oldKey := sc.keys[sc.currentKeyIndex]
		delete(sc.cachedScenes, oldKey)
		sc.keys[sc.currentKeyIndex] = key
	}

	sc.cachedScenes[key] = scene
	sc.currentKeyIndex = (sc.currentKeyIndex + 1) % sc.limit
}

func (sc *sceneCache) check(key string) (bool, *Scene) {
	scene, ok := sc.cachedScenes[key]
	return ok, scene
}

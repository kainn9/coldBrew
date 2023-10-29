package coldBrew

type tickHandler struct {
	currentTick int
	maxTick     int
}

func NewTickHandler(maxTick int) *tickHandler {
	return &tickHandler{
		currentTick: 0,
		maxTick:     maxTick,
	}
}

func (t *tickHandler) CurrentTick() int {
	return t.currentTick
}

func (t *tickHandler) IncrementTick() {
	t.currentTick++

	if t.currentTick >= t.maxTick {
		t.currentTick = 0
	}

}

func (t *tickHandler) TicksSinceNTicks(n int) int {

	if n < 0 || n >= t.maxTick {
		return 0
	}

	if n > t.currentTick {
		return t.maxTick - n + t.currentTick
	}

	return t.currentTick - n
}

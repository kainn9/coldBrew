package coldBrew

type TickHandler struct {
	currentTick int
	maxTick     int
}

func NewTickHandler(maxTick int) *TickHandler {
	return &TickHandler{
		currentTick: 0,
		maxTick:     maxTick,
	}
}

func (t *TickHandler) CurrentTick() int {
	return t.currentTick
}

func (t *TickHandler) IncrementTick() {
	t.currentTick++

	if t.currentTick >= t.maxTick {
		t.currentTick = 0
	}

}

func (t *TickHandler) TicksSinceNTicks(n int) int {

	if n < 0 || n >= t.maxTick {
		return 0
	}

	if n > t.currentTick {
		return t.maxTick - n + t.currentTick
	}

	return t.currentTick - n
}

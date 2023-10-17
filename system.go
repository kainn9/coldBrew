package coldBrew

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type System interface{}

type SystemType int

const (
	ClientType SystemType = iota
	RenderType
	SimType
	LoaderType
)

// Query interface for systems that need to query a world for a single type of entity.
// By default, the system will iterate through all entities returned by the query and
// run the main func(sim, render, etc) on each entry/entity.
// If more robust querying or combinations of entities are required, do not 
// implement this.
type Query interface {
	Query() *donburi.Query
}

// Sim interface for systems that simulate the game logic.
type Sim interface {
	Run(dt float64, entry *donburi.Entry)
}

// Render interface for systems that render the game, based on the game state.
type Render interface {
	Draw(screen *ebiten.Image, entry *donburi.Entry)
}

// Client interface for systems that handle
// "client logic" (i.e. input, sounds, etc).
type Client interface {
	Sync(entry *donburi.Entry)
}

// For loading assets.
type Load interface {
	Load(entry *donburi.Entry)
}

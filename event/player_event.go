package event

//import "github.com/cr0sh/encore"

/**
 * Player Event
*/
type PlayerEvent struct {
	*Event
	Player string //*encore.Player
}

/**
 * Player Join Event
 */
type PlayerJoinEvent struct {
	*PlayerEvent
	Handlers HandlerList
	JoinMessage string
}

func (e *PlayerJoinEvent) New(player /*encore.Player*/string, joinMessage string) *PlayerJoinEvent{
	e.Name = "PlayerJoinEvent"
	e.Handlers = HandlerList{}
	e.IsCancelled = false
	e.JoinMessage = joinMessage
	e.Player = player
	return e
}



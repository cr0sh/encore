package event

//Event
type Event struct {
	Name string
	IsCancelled bool
}

//get name of event
func (e Event) GetEventName() string {
	if e.Name != "" {
		return "Event"
	} else {
		return e.Name
	}
}

type Listener interface {}

type HandlerList struct {}
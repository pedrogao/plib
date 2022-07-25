package actor

type Actor[Msg any] interface {
	Send(msg Msg)
}

type reducerActor[Msg, St any] struct {
	state   St
	queue   chan Msg
	reducer func(Msg, St) St
}

func NewFromReducer[Msg, St any](initial St,
	reducer func(Msg, St) St) Actor[Msg] {
	ga := &reducerActor[Msg, St]{
		state:   initial,
		queue:   make(chan Msg, 1),
		reducer: reducer,
	}
	ga.start()
	return ga
}

func (r *reducerActor[Msg, _]) Send(msg Msg) {
	r.queue <- msg
}

func (r *reducerActor[_, _]) start() {
	go r.receiveLoop()
}

func (r *reducerActor[_, _]) receiveLoop() {
	for msg := range r.queue {
		r.state = r.reducer(msg, r.state)
	}
}

type Message[St any] interface {
	Apply(St) St
}

func NewTyped[St any](initial St) Actor[Message[St]] {
	return NewFromReducer(initial,
		func(msg Message[St], state St) St {
			return msg.Apply(state)
		},
	)
}

// getProj Message implement for add & get
type getProj[St, Ret any] struct {
	reply      chan Ret
	projection func(St) Ret
}

func (g *getProj[St, _]) Apply(state St) St {
	g.reply <- g.projection(state) // 映射
	return state
}

// GetAsync applies the projection func to the actor's state and returns the
// result asynchronously
func GetAsync[St, Ret any](a Actor[Message[St]],
	projection func(St) Ret) <-chan Ret {
	reply := make(chan Ret, 1)
	gp := &getProj[St, Ret]{reply: reply, projection: projection}
	a.Send(gp)
	return reply
}

// Get applies the projection func to the actor's state and returns the result
func Get[St, Ret any](a Actor[Message[St]], projection func(St) Ret) Ret {
	reply := GetAsync(a, projection)
	return <-reply
}

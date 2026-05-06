// Package fsm provides a lightweight, unidirectional Finite State Machine engine.
// Import path: example.com/fsm
// Go version: 1.21+ (uses generics)
// Dependencies: standard library only.
//
// The FSM is NOT safe for concurrent use; callers must serialize access.
package fsm

import (
	"errors"
	"fmt"
)

// State represents a distinct phase in the workflow.
type State string

// Event represents a trigger that causes a transition.
type Event string

// GuardFunc inspects the context and determines whether a transition is allowed.
// If error != nil, the transition fails immediately with that error.
// If allowed == false with nil error, the transition is treated as invalid.
type GuardFunc[C any] func(ctx *C) (bool, error)

// ActionFunc is a side-effect callback executed on state entry or exit.
// Actions must not fail (no error return) and should not panic.
type ActionFunc[C any] func(ctx *C)

// TransitionActionFunc is executed during a transition, after OnExit and before OnEnter.
type TransitionActionFunc[C any] func(ctx *C, event Event)

// Transition defines a directed edge from one state to another.
type Transition[C any] struct {
	From   State
	Event  Event
	To     State
	Guard  GuardFunc[C]            // optional; nil = always allowed
	Action TransitionActionFunc[C] // optional
}

// InvalidReason indicates why a transition was rejected.
type InvalidReason int

const (
	// ReasonNoTransition: no transition defined for (from, event).
	ReasonNoTransition InvalidReason = iota
	// ReasonGuardDenied: guard returned (false, nil).
	ReasonGuardDenied
)

// Observer allows external code to observe FSM transitions without modifying core logic.
type Observer[C any] interface {
	// OnTransition is called after a successful transition completes
	// (state updated, OnEnter action finished).
	OnTransition(from, to State, event Event, ctx *C)
	// OnInvalidTransition is called when Send fails due to no transition or guard denial.
	OnInvalidTransition(from State, event Event, ctx *C, reason InvalidReason)
}

// InvalidTransitionHandler is a standalone callback for invalid transitions.
// Called before any observers' OnInvalidTransition.
type InvalidTransitionHandler[C any] func(from State, event Event, ctx *C, reason InvalidReason)

// Predefined errors.
var (
	ErrInvalidTransition = errors.New("fsm: invalid transition")
	ErrNilContext        = errors.New("fsm: nil context")
)

// config holds optional FSM configuration provided at Build time.
type config[C any] struct {
	observers                []Observer[C]
	invalidTransitionHandler InvalidTransitionHandler[C]
}

// Option configures optional FSM behavior at Build time.
type Option[C any] func(*config[C])

// WithObserver attaches an observer to the FSM.
// Observers are called in order after successful transitions.
func WithObserver[C any](o Observer[C]) Option[C] {
	return func(cfg *config[C]) {
		cfg.observers = append(cfg.observers, o)
	}
}

// WithInvalidTransitionHandler sets the handler for invalid transitions.
// This handler is called before any observers' OnInvalidTransition.
func WithInvalidTransitionHandler[C any](h InvalidTransitionHandler[C]) Option[C] {
	return func(cfg *config[C]) {
		cfg.invalidTransitionHandler = h
	}
}

// Builder constructs an FSM instance imperatively.
// It is generic over the user-defined context type C.
type Builder[C any] struct {
	states      map[State]struct{}
	transitions map[State]map[Event]Transition[C] // from -> event -> transition
	onEnter     map[State]ActionFunc[C]
	onExit      map[State]ActionFunc[C]
}

// NewBuilder creates a new FSM builder generic over context type C.
func NewBuilder[C any]() *Builder[C] {
	return &Builder[C]{
		states:      make(map[State]struct{}),
		transitions: make(map[State]map[Event]Transition[C]),
		onEnter:     make(map[State]ActionFunc[C]),
		onExit:      make(map[State]ActionFunc[C]),
	}
}

// AddState registers a state. Adding the same state multiple times is a no-op.
func (b *Builder[C]) AddState(state State) *Builder[C] {
	if state == "" {
		// Empty state names are invalid; defer error to Build() for consistent validation.
		return b
	}
	b.states[state] = struct{}{}
	return b
}

// AddTransition registers a transition.
// Duplicate (from, event) pairs silently overwrite previous definitions (last write wins).
// Validation is deferred to Build().
func (b *Builder[C]) AddTransition(t Transition[C]) *Builder[C] {
	if t.From == "" || t.Event == "" || t.To == "" {
		// Empty identifiers are invalid; defer to Build().
		return b
	}
	if _, ok := b.transitions[t.From]; !ok {
		b.transitions[t.From] = make(map[Event]Transition[C])
	}
	b.transitions[t.From][t.Event] = t
	return b
}

// OnEnter registers an action to execute every time the machine enters the given state.
// Overwrites any previously registered action for that state.
func (b *Builder[C]) OnEnter(state State, action ActionFunc[C]) *Builder[C] {
	if state == "" {
		return b
	}
	b.onEnter[state] = action
	return b
}

// OnExit registers an action to execute every time the machine leaves the given state.
// Overwrites any previously registered action for that state.
func (b *Builder[C]) OnExit(state State, action ActionFunc[C]) *Builder[C] {
	if state == "" {
		return b
	}
	b.onExit[state] = action
	return b
}

// Build validates the state machine definition and returns a ready-to-use FSM instance.
// Validation errors:
//   - initial state must be registered via AddState
//   - all states referenced in transitions must be registered
//   - no empty State or Event identifiers allowed
func (b *Builder[C]) Build(initial State, opts ...Option[C]) (*FSM[C], error) {
	// Validate initial state
	if initial == "" {
		return nil, errors.New("fsm: initial state cannot be empty")
	}
	if _, ok := b.states[initial]; !ok {
		return nil, fmt.Errorf("fsm: initial state %q not registered", initial)
	}

	// Validate all referenced states in transitions
	for from, events := range b.transitions {
		if _, ok := b.states[from]; !ok {
			return nil, fmt.Errorf("fsm: transition source state %q not registered", from)
		}
		for _, t := range events {
			if _, ok := b.states[t.To]; !ok {
				return nil, fmt.Errorf("fsm: transition target state %q not registered", t.To)
			}
			// Validate event is non-empty (already checked in AddTransition, but be defensive)
			if t.Event == "" {
				return nil, errors.New("fsm: transition event cannot be empty")
			}
		}
	}

	// Apply options
	cfg := &config[C]{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &FSM[C]{
		builder:        b,
		current:        initial,
		observers:      cfg.observers,
		invalidHandler: cfg.invalidTransitionHandler,
	}, nil
}

// FSM is a unidirectional finite state machine instance.
// It is NOT safe for concurrent use; callers must serialize access.
type FSM[C any] struct {
	builder        *Builder[C]
	current        State
	observers      []Observer[C]
	invalidHandler InvalidTransitionHandler[C]
}

// CurrentState returns the current state of the machine.
func (m *FSM[C]) CurrentState() State {
	return m.current
}

// Send fires an event, attempting to perform the corresponding transition.
// The context pointer must not be nil (returns ErrNilContext if it is).
// Returns nil on success, or an error if the transition is impossible or a guard failed.
//
// Execution order on success:
// 1. Guard evaluation (if present)
// 2. OnExit action (if registered for current state)
// 3. Transition action (if present)
// 4. State update to target
// 5. OnEnter action (if registered for new state)
// 6. Observer notifications
//
// Panic semantics:
// - If OnExit or transition action panics: state remains unchanged (update not yet performed)
// - If OnEnter panics: state has already been updated to target, but entry action incomplete
// - No rollback of side effects is attempted
func (m *FSM[C]) Send(event Event, ctx *C) error {
	if ctx == nil {
		return ErrNilContext
	}

	from := m.current

	// Look up transition
	events, hasFrom := m.builder.transitions[from]
	if !hasFrom {
		// No transitions defined from this state
		return m.handleInvalidTransition(from, event, ctx, ReasonNoTransition)
	}

	t, exists := events[event]
	if !exists {
		// No transition for this (from, event) pair
		return m.handleInvalidTransition(from, event, ctx, ReasonNoTransition)
	}

	// Evaluate guard if present
	if t.Guard != nil {
		allowed, err := t.Guard(ctx)
		if err != nil {
			// Guard error: return immediately, state unchanged, no hooks
			return err
		}
		if !allowed {
			// Guard denied: invalid transition
			return m.handleInvalidTransition(from, event, ctx, ReasonGuardDenied)
		}
	}

	// Execute OnExit action if registered
	if onExit, ok := m.builder.onExit[from]; ok {
		onExit(ctx)
	}

	// Execute transition action if registered
	if t.Action != nil {
		t.Action(ctx, event)
	}

	// Update state (point of no return for state change)
	m.current = t.To

	// Execute OnEnter action if registered
	if onEnter, ok := m.builder.onEnter[m.current]; ok {
		onEnter(ctx)
	}

	// Notify observers of successful transition (after OnEnter completes)
	for _, obs := range m.observers {
		obs.OnTransition(from, m.current, event, ctx)
	}

	return nil
}

// handleInvalidTransition handles the case where a transition cannot be taken.
// Returns ErrInvalidTransition after invoking configured hooks.
func (m *FSM[C]) handleInvalidTransition(from State, event Event, ctx *C, reason InvalidReason) error {
	// Call standalone handler first (if configured)
	if m.invalidHandler != nil {
		m.invalidHandler(from, event, ctx, reason)
	}
	// Then notify observers
	for _, obs := range m.observers {
		obs.OnInvalidTransition(from, event, ctx, reason)
	}
	return ErrInvalidTransition
}

// Reset immediately sets the current state to the given state.
// The state must be a registered state (previously added via AddState).
// No guards, actions, or observers are invoked.
// Panics if the state is not registered (caller must ensure validity).
func (m *FSM[C]) Reset(state State) {
	// Per spec: state must be registered. Panic on misuse to fail fast.
	if _, ok := m.builder.states[state]; !ok {
		panic(fmt.Sprintf("fsm: Reset to unregistered state %q", state))
	}
	m.current = state
}

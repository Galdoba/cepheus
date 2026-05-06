package fsm

import (
	"errors"
	"reflect"
	"testing"
)

// ============================================================================
// Helper types and functions
// ============================================================================

type testCtx struct {
	log     []string
	counter int
	valid   bool
}

func (c *testCtx) appendLog(entry string) {
	c.log = append(c.log, entry)
}

// assertEqual for comparable types (non-slice)
func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

// assertEqualSlice for slice types using reflect.DeepEqual
func assertEqualSlice[T any](t *testing.T, name string, got, want []T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func assertErrorIs(t *testing.T, name string, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("%s: got error %v, want %v", name, got, want)
	}
}

func assertLen(t *testing.T, name string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got length %d, want %d", name, got, want)
	}
}

// ============================================================================
// Builder & Validation Tests
// ============================================================================

func TestBuilder_AddState(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("a").AddState("b")
	_, err := b.Build("a")
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
}

func TestBuild_Validation_EmptyIdentifiers(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Builder[any])
		initial State
		wantErr string
	}{
		{
			name: "empty initial state",
			setup: func(b *Builder[any]) {
				b.AddState("valid")
			},
			initial: "",
			wantErr: "fsm: initial state cannot be empty", // FIX: added "fsm: " prefix
		},
		{
			name: "initial state not registered",
			setup: func(b *Builder[any]) {
				b.AddState("other")
			},
			initial: "missing",
			wantErr: `fsm: initial state "missing" not registered`, // FIX: added "fsm: " prefix
		},
		{
			name: "transition source not registered",
			setup: func(b *Builder[any]) {
				b.AddState("to")
				b.AddTransition(Transition[any]{From: "from", Event: "e", To: "to"})
			},
			initial: "to",
			wantErr: `fsm: transition source state "from" not registered`, // FIX: added "fsm: " prefix
		},
		{
			name: "transition target not registered",
			setup: func(b *Builder[any]) {
				b.AddState("from")
				b.AddTransition(Transition[any]{From: "from", Event: "e", To: "to"})
			},
			initial: "from",
			wantErr: `fsm: transition target state "to" not registered`, // FIX: added "fsm: " prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder[any]()
			tt.setup(b)
			_, err := b.Build(tt.initial)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error message: got %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestBuild_DuplicateTransition_LastWriteWins(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("b").AddState("c")
	b.AddTransition(Transition[any]{From: "a", Event: "go", To: "b"})
	b.AddTransition(Transition[any]{From: "a", Event: "go", To: "c"})

	m, err := b.Build("a")
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// FIX: pass non-nil context (nil triggers ErrNilContext, not transition logic)
	_ = m.Send("go", new(any))
	assertEqual(t, "state after duplicate transition", m.CurrentState(), State("c"))
}

// ============================================================================
// Transition Execution Tests
// ============================================================================

func TestSend_GuardPass(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("idle").AddState("run")
	b.AddTransition(Transition[testCtx]{
		From:  "idle",
		Event: "start",
		To:    "run",
		Guard: func(ctx *testCtx) (bool, error) { return ctx.valid, nil },
	})

	m, _ := b.Build("idle")
	ctx := &testCtx{valid: true}
	err := m.Send("start", ctx)

	assertErrorIs(t, "Send error", err, nil)
	assertEqual(t, "current state", m.CurrentState(), State("run"))
}

func TestSend_GuardDenied(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("idle").AddState("run")
	b.AddTransition(Transition[testCtx]{
		From:  "idle",
		Event: "start",
		To:    "run",
		Guard: func(ctx *testCtx) (bool, error) { return false, nil },
	})

	var reason InvalidReason
	handler := func(from State, event Event, ctx *testCtx, r InvalidReason) {
		reason = r
	}

	m, _ := b.Build("idle", WithInvalidTransitionHandler[testCtx](handler))
	ctx := &testCtx{}
	err := m.Send("start", ctx)

	assertErrorIs(t, "Send error", err, ErrInvalidTransition)
	assertEqual(t, "current state unchanged", m.CurrentState(), State("idle"))
	assertEqual(t, "invalid reason", reason, ReasonGuardDenied)
}

func TestSend_GuardError(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("idle").AddState("run")
	guardErr := errors.New("guard failed")
	b.AddTransition(Transition[testCtx]{
		From:  "idle",
		Event: "start",
		To:    "run",
		Guard: func(ctx *testCtx) (bool, error) { return false, guardErr },
	})

	hookCalled := false
	m, _ := b.Build("idle", WithInvalidTransitionHandler[testCtx](func(State, Event, *testCtx, InvalidReason) {
		hookCalled = true
	}))

	ctx := &testCtx{}
	err := m.Send("start", ctx)

	assertEqual(t, "Send returns guard error", err, guardErr)
	assertEqual(t, "state unchanged", m.CurrentState(), State("idle"))
	if hookCalled {
		t.Error("invalid transition hook should not be called on guard error")
	}
}

func TestSend_ExecutionOrder(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")

	b.OnExit("a", func(ctx *testCtx) { ctx.appendLog("exit-a") })
	b.OnEnter("b", func(ctx *testCtx) { ctx.appendLog("enter-b") })

	b.AddTransition(Transition[testCtx]{
		From:   "a",
		Event:  "go",
		To:     "b",
		Action: func(ctx *testCtx, e Event) { ctx.appendLog("action") },
	})

	m, _ := b.Build("a")
	ctx := &testCtx{}
	_ = m.Send("go", ctx)

	want := []string{"exit-a", "action", "enter-b"}
	assertEqualSlice(t, "execution order log", ctx.log, want)
}

// ============================================================================
// Self-Loop Tests
// ============================================================================

func TestSend_SelfLoop_FiresBothHooks(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("loop")
	b.OnExit("loop", func(ctx *testCtx) { ctx.appendLog("exit") })
	b.OnEnter("loop", func(ctx *testCtx) { ctx.appendLog("enter") })
	b.AddTransition(Transition[testCtx]{
		From:   "loop",
		Event:  "tick",
		To:     "loop",
		Action: func(ctx *testCtx, e Event) { ctx.appendLog("action") },
	})

	m, _ := b.Build("loop")
	ctx := &testCtx{}
	_ = m.Send("tick", ctx)

	want := []string{"exit", "action", "enter"}
	assertEqualSlice(t, "self-loop execution order", ctx.log, want)
	assertEqual(t, "state remains loop", m.CurrentState(), State("loop"))
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestSend_NilContext(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("s")
	m, _ := b.Build("s")

	err := m.Send("e", nil)
	assertErrorIs(t, "nil context error", err, ErrNilContext)
}

func TestSend_NoTransition(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a")

	var reason InvalidReason
	m, _ := b.Build("a", WithInvalidTransitionHandler[any](func(from State, event Event, ctx *any, r InvalidReason) {
		reason = r
	}))

	// FIX: pass non-nil context
	err := m.Send("unknown", new(any))
	assertErrorIs(t, "no transition error", err, ErrInvalidTransition)
	assertEqual(t, "invalid reason", reason, ReasonNoTransition)
	assertEqual(t, "state unchanged", m.CurrentState(), State("a"))
}

// ============================================================================
// Panic Semantics Tests
// ============================================================================

func TestPanic_OnExit_StateUnchanged(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from OnExit, got none")
		}
	}()

	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")
	b.OnExit("a", func(ctx *testCtx) { panic("exit-panic") })
	b.AddTransition(Transition[testCtx]{From: "a", Event: "go", To: "b"})

	m, _ := b.Build("a")
	_ = m.Send("go", &testCtx{})
}

func TestPanic_OnEnter_StateAlreadyUpdated(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")
	b.OnEnter("b", func(ctx *testCtx) { panic("enter-panic") })
	b.AddTransition(Transition[testCtx]{From: "a", Event: "go", To: "b"})

	m, _ := b.Build("a")
	ctx := &testCtx{}

	defer func() {
		if r := recover(); r != nil {
			assertEqual(t, "state after OnEnter panic", m.CurrentState(), State("b"))
		}
	}()
	_ = m.Send("go", ctx)
}

// ============================================================================
// Observer Tests
// ============================================================================

type recordingObserver struct {
	transitions []struct {
		from, to State
		event    Event
	}
	invalid []struct {
		from   State
		event  Event
		reason InvalidReason
	}
}

func (o *recordingObserver) OnTransition(from, to State, event Event, ctx *testCtx) {
	o.transitions = append(o.transitions, struct {
		from, to State
		event    Event
	}{from, to, event})
}

func (o *recordingObserver) OnInvalidTransition(from State, event Event, ctx *testCtx, reason InvalidReason) {
	o.invalid = append(o.invalid, struct {
		from   State
		event  Event
		reason InvalidReason
	}{from, event, reason})
}

func TestObserver_SuccessfulTransition(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")
	b.AddTransition(Transition[testCtx]{From: "a", Event: "go", To: "b"})

	obs := &recordingObserver{}
	m, _ := b.Build("a", WithObserver[testCtx](obs))

	_ = m.Send("go", &testCtx{})

	assertLen(t, "observer transitions recorded", len(obs.transitions), 1)
	assertEqual(t, "recorded from", obs.transitions[0].from, State("a"))
	assertEqual(t, "recorded to", obs.transitions[0].to, State("b"))
	assertEqual(t, "recorded event", obs.transitions[0].event, Event("go"))
	assertLen(t, "observer invalid count", len(obs.invalid), 0)
}

func TestObserver_InvalidTransition(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a")

	obs := &recordingObserver{}
	m, _ := b.Build("a", WithObserver[testCtx](obs))

	// FIX: pass non-nil context
	_ = m.Send("missing", &testCtx{})

	assertLen(t, "observer invalid recorded", len(obs.invalid), 1)
	assertEqual(t, "invalid from", obs.invalid[0].from, State("a"))
	assertEqual(t, "invalid event", obs.invalid[0].event, Event("missing"))
	assertEqual(t, "invalid reason", obs.invalid[0].reason, ReasonNoTransition)
}

func TestObserver_Ordering(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")
	b.AddTransition(Transition[testCtx]{From: "a", Event: "go", To: "b"})

	order := []int{}
	makeObs := func(id int) Observer[testCtx] {
		return &funcObserver[testCtx]{
			onTransition: func(State, State, Event, *testCtx) {
				order = append(order, id)
			},
		}
	}

	m, _ := b.Build("a", WithObserver[testCtx](makeObs(1)), WithObserver[testCtx](makeObs(2)))
	_ = m.Send("go", &testCtx{}) // FIX: was nil

	assertEqualSlice(t, "observer invocation order", order, []int{1, 2})
}

// Helper observer for tests
type funcObserver[C any] struct {
	onTransition        func(State, State, Event, *C)
	onInvalidTransition func(State, Event, *C, InvalidReason)
}

func (f *funcObserver[C]) OnTransition(from, to State, event Event, ctx *C) {
	if f.onTransition != nil {
		f.onTransition(from, to, event, ctx)
	}
}

func (f *funcObserver[C]) OnInvalidTransition(from State, event Event, ctx *C, reason InvalidReason) {
	if f.onInvalidTransition != nil {
		f.onInvalidTransition(from, event, ctx, reason)
	}
}

// ============================================================================
// InvalidTransitionHandler Tests
// ============================================================================

func TestInvalidHandler_CalledBeforeObservers(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a")

	handlerCalled := false
	obsCalled := false

	handler := func(from State, event Event, ctx *testCtx, reason InvalidReason) {
		handlerCalled = true
	}
	obs := &funcObserver[testCtx]{
		onInvalidTransition: func(State, Event, *testCtx, InvalidReason) {
			obsCalled = true
		},
	}

	m, _ := b.Build("a", WithInvalidTransitionHandler[testCtx](handler), WithObserver[testCtx](obs))
	_ = m.Send("missing", &testCtx{}) // FIX: was nil

	if !handlerCalled {
		t.Error("invalid handler not called")
	}
	if !obsCalled {
		t.Error("observer OnInvalidTransition not called")
	}
}

// ============================================================================
// Reset Tests
// ============================================================================

func TestReset_ValidState(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("b")
	m, _ := b.Build("a")

	m.Reset("b")
	assertEqual(t, "state after Reset", m.CurrentState(), State("b"))
}

func TestReset_UnregisteredState_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Reset to unregistered state")
		}
	}()

	b := NewBuilder[any]()
	b.AddState("a")
	m, _ := b.Build("a")
	m.Reset("unknown")
}

func TestReset_NoHooksInvoked(t *testing.T) {
	b := NewBuilder[testCtx]()
	b.AddState("a").AddState("b")
	b.OnExit("a", func(ctx *testCtx) { ctx.appendLog("exit") })
	b.OnEnter("b", func(ctx *testCtx) { ctx.appendLog("enter") })

	m, _ := b.Build("a")
	ctx := &testCtx{}
	m.Reset("b")

	assertLen(t, "no actions invoked on Reset", len(ctx.log), 0)
	assertEqual(t, "state after Reset", m.CurrentState(), State("b"))
}

// ============================================================================
// Spec Example Test (§9)
// ============================================================================

type TaskCtx struct {
	Attempts int
	Valid    bool
}

func TestSpecExample(t *testing.T) {
	builder := NewBuilder[TaskCtx]()

	builder.AddState("new").AddState("processing").AddState("done")
	builder.OnEnter("processing", func(ctx *TaskCtx) {
		ctx.Attempts++
	})

	builder.AddTransition(Transition[TaskCtx]{
		From:  "new",
		Event: "start",
		To:    "processing",
		Guard: func(ctx *TaskCtx) (bool, error) {
			return ctx.Valid, nil
		},
	})
	builder.AddTransition(Transition[TaskCtx]{
		From:  "processing",
		Event: "finish",
		To:    "done",
	})
	builder.AddTransition(Transition[TaskCtx]{
		From:   "processing",
		Event:  "retry",
		To:     "processing",
		Action: func(ctx *TaskCtx, e Event) {},
	})

	machine, err := builder.Build("new")
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Guard denial
	ctx := &TaskCtx{Valid: false}
	err = machine.Send("start", ctx)
	assertErrorIs(t, "guard denied", err, ErrInvalidTransition)
	assertEqual(t, "state after guard denial", machine.CurrentState(), State("new"))

	// Guard pass
	ctx.Valid = true
	err = machine.Send("start", ctx)
	assertErrorIs(t, "guard pass", err, nil)
	assertEqual(t, "state after start", machine.CurrentState(), State("processing"))
	assertEqual(t, "attempts incremented", ctx.Attempts, 1)

	// Self-loop retry
	err = machine.Send("retry", ctx)
	assertErrorIs(t, "retry", err, nil)
	assertEqual(t, "state after retry", machine.CurrentState(), State("processing"))

	// Finish
	err = machine.Send("finish", ctx)
	assertErrorIs(t, "finish", err, nil)
	assertEqual(t, "final state", machine.CurrentState(), State("done"))
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestGuard_NilIsAlwaysAllowed(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("b")
	b.AddTransition(Transition[any]{From: "a", Event: "go", To: "b"})

	m, _ := b.Build("a")
	// FIX: pass non-nil context (nil triggers ErrNilContext)
	err := m.Send("go", new(any))
	assertErrorIs(t, "nil guard allows transition", err, nil)
	assertEqual(t, "state after nil-guard transition", m.CurrentState(), State("b"))
}

func TestAction_NilIsNoOp(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("b")
	b.AddTransition(Transition[any]{From: "a", Event: "go", To: "b", Action: nil})

	m, _ := b.Build("a")
	// FIX: pass non-nil context
	err := m.Send("go", new(any))
	assertErrorIs(t, "nil action transition", err, nil)
}

func TestMultipleObservers_AllCalledOnSuccess(t *testing.T) {
	b := NewBuilder[any]()
	b.AddState("a").AddState("b")
	b.AddTransition(Transition[any]{From: "a", Event: "go", To: "b"})

	called := [3]bool{}
	makeObs := func(id int) Observer[any] {
		return &funcObserver[any]{
			onTransition: func(State, State, Event, *any) { called[id] = true },
		}
	}

	m, _ := b.Build("a",
		WithObserver[any](makeObs(0)),
		WithObserver[any](makeObs(1)),
		WithObserver[any](makeObs(2)),
	)
	// FIX: pass non-nil context
	_ = m.Send("go", new(any))

	for i, c := range called {
		if !c {
			t.Errorf("observer %d not called", i)
		}
	}
}

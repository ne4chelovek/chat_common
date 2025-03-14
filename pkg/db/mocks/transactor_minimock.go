package mocks

//go:generate minimock -i chat-server/internal/client/db.Transactor -o transactor_minimock.go -n TransactorMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
)

// TransactorMock implements mm_db.Transactor
type TransactorMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcBeginTx          func(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error)
	funcBeginTxOrigin    string
	inspectFuncBeginTx   func(ctx context.Context, txOptions pgx.TxOptions)
	afterBeginTxCounter  uint64
	beforeBeginTxCounter uint64
	BeginTxMock          mTransactorMockBeginTx
}

// NewTransactorMock returns a mock for mm_db.Transactor
func NewTransactorMock(t minimock.Tester) *TransactorMock {
	m := &TransactorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginTxMock = mTransactorMockBeginTx{mock: m}
	m.BeginTxMock.callArgs = []*TransactorMockBeginTxParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mTransactorMockBeginTx struct {
	optional           bool
	mock               *TransactorMock
	defaultExpectation *TransactorMockBeginTxExpectation
	expectations       []*TransactorMockBeginTxExpectation

	callArgs []*TransactorMockBeginTxParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// TransactorMockBeginTxExpectation specifies expectation struct of the Transactor.BeginTx
type TransactorMockBeginTxExpectation struct {
	mock               *TransactorMock
	params             *TransactorMockBeginTxParams
	paramPtrs          *TransactorMockBeginTxParamPtrs
	expectationOrigins TransactorMockBeginTxExpectationOrigins
	results            *TransactorMockBeginTxResults
	returnOrigin       string
	Counter            uint64
}

// TransactorMockBeginTxParams contains parameters of the Transactor.BeginTx
type TransactorMockBeginTxParams struct {
	ctx       context.Context
	txOptions pgx.TxOptions
}

// TransactorMockBeginTxParamPtrs contains pointers to parameters of the Transactor.BeginTx
type TransactorMockBeginTxParamPtrs struct {
	ctx       *context.Context
	txOptions *pgx.TxOptions
}

// TransactorMockBeginTxResults contains results of the Transactor.BeginTx
type TransactorMockBeginTxResults struct {
	t1  pgx.Tx
	err error
}

// TransactorMockBeginTxOrigins contains origins of expectations of the Transactor.BeginTx
type TransactorMockBeginTxExpectationOrigins struct {
	origin          string
	originCtx       string
	originTxOptions string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmBeginTx *mTransactorMockBeginTx) Optional() *mTransactorMockBeginTx {
	mmBeginTx.optional = true
	return mmBeginTx
}

// Expect sets up expected params for Transactor.BeginTx
func (mmBeginTx *mTransactorMockBeginTx) Expect(ctx context.Context, txOptions pgx.TxOptions) *mTransactorMockBeginTx {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &TransactorMockBeginTxExpectation{}
	}

	if mmBeginTx.defaultExpectation.paramPtrs != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by ExpectParams functions")
	}

	mmBeginTx.defaultExpectation.params = &TransactorMockBeginTxParams{ctx, txOptions}
	mmBeginTx.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmBeginTx.expectations {
		if minimock.Equal(e.params, mmBeginTx.defaultExpectation.params) {
			mmBeginTx.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmBeginTx.defaultExpectation.params)
		}
	}

	return mmBeginTx
}

// ExpectCtxParam1 sets up expected param ctx for Transactor.BeginTx
func (mmBeginTx *mTransactorMockBeginTx) ExpectCtxParam1(ctx context.Context) *mTransactorMockBeginTx {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &TransactorMockBeginTxExpectation{}
	}

	if mmBeginTx.defaultExpectation.params != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Expect")
	}

	if mmBeginTx.defaultExpectation.paramPtrs == nil {
		mmBeginTx.defaultExpectation.paramPtrs = &TransactorMockBeginTxParamPtrs{}
	}
	mmBeginTx.defaultExpectation.paramPtrs.ctx = &ctx
	mmBeginTx.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmBeginTx
}

// ExpectTxOptionsParam2 sets up expected param txOptions for Transactor.BeginTx
func (mmBeginTx *mTransactorMockBeginTx) ExpectTxOptionsParam2(txOptions pgx.TxOptions) *mTransactorMockBeginTx {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &TransactorMockBeginTxExpectation{}
	}

	if mmBeginTx.defaultExpectation.params != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Expect")
	}

	if mmBeginTx.defaultExpectation.paramPtrs == nil {
		mmBeginTx.defaultExpectation.paramPtrs = &TransactorMockBeginTxParamPtrs{}
	}
	mmBeginTx.defaultExpectation.paramPtrs.txOptions = &txOptions
	mmBeginTx.defaultExpectation.expectationOrigins.originTxOptions = minimock.CallerInfo(1)

	return mmBeginTx
}

// Inspect accepts an inspector function that has same arguments as the Transactor.BeginTx
func (mmBeginTx *mTransactorMockBeginTx) Inspect(f func(ctx context.Context, txOptions pgx.TxOptions)) *mTransactorMockBeginTx {
	if mmBeginTx.mock.inspectFuncBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("Inspect function is already set for TransactorMock.BeginTx")
	}

	mmBeginTx.mock.inspectFuncBeginTx = f

	return mmBeginTx
}

// Return sets up results that will be returned by Transactor.BeginTx
func (mmBeginTx *mTransactorMockBeginTx) Return(t1 pgx.Tx, err error) *TransactorMock {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &TransactorMockBeginTxExpectation{mock: mmBeginTx.mock}
	}
	mmBeginTx.defaultExpectation.results = &TransactorMockBeginTxResults{t1, err}
	mmBeginTx.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmBeginTx.mock
}

// Set uses given function f to mock the Transactor.BeginTx method
func (mmBeginTx *mTransactorMockBeginTx) Set(f func(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error)) *TransactorMock {
	if mmBeginTx.defaultExpectation != nil {
		mmBeginTx.mock.t.Fatalf("Default expectation is already set for the Transactor.BeginTx method")
	}

	if len(mmBeginTx.expectations) > 0 {
		mmBeginTx.mock.t.Fatalf("Some expectations are already set for the Transactor.BeginTx method")
	}

	mmBeginTx.mock.funcBeginTx = f
	mmBeginTx.mock.funcBeginTxOrigin = minimock.CallerInfo(1)
	return mmBeginTx.mock
}

// When sets expectation for the Transactor.BeginTx which will trigger the result defined by the following
// Then helper
func (mmBeginTx *mTransactorMockBeginTx) When(ctx context.Context, txOptions pgx.TxOptions) *TransactorMockBeginTxExpectation {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("TransactorMock.BeginTx mock is already set by Set")
	}

	expectation := &TransactorMockBeginTxExpectation{
		mock:               mmBeginTx.mock,
		params:             &TransactorMockBeginTxParams{ctx, txOptions},
		expectationOrigins: TransactorMockBeginTxExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmBeginTx.expectations = append(mmBeginTx.expectations, expectation)
	return expectation
}

// Then sets up Transactor.BeginTx return parameters for the expectation previously defined by the When method
func (e *TransactorMockBeginTxExpectation) Then(t1 pgx.Tx, err error) *TransactorMock {
	e.results = &TransactorMockBeginTxResults{t1, err}
	return e.mock
}

// Times sets number of times Transactor.BeginTx should be invoked
func (mmBeginTx *mTransactorMockBeginTx) Times(n uint64) *mTransactorMockBeginTx {
	if n == 0 {
		mmBeginTx.mock.t.Fatalf("Times of TransactorMock.BeginTx mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmBeginTx.expectedInvocations, n)
	mmBeginTx.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmBeginTx
}

func (mmBeginTx *mTransactorMockBeginTx) invocationsDone() bool {
	if len(mmBeginTx.expectations) == 0 && mmBeginTx.defaultExpectation == nil && mmBeginTx.mock.funcBeginTx == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmBeginTx.mock.afterBeginTxCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmBeginTx.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// BeginTx implements mm_db.Transactor
func (mmBeginTx *TransactorMock) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (t1 pgx.Tx, err error) {
	mm_atomic.AddUint64(&mmBeginTx.beforeBeginTxCounter, 1)
	defer mm_atomic.AddUint64(&mmBeginTx.afterBeginTxCounter, 1)

	mmBeginTx.t.Helper()

	if mmBeginTx.inspectFuncBeginTx != nil {
		mmBeginTx.inspectFuncBeginTx(ctx, txOptions)
	}

	mm_params := TransactorMockBeginTxParams{ctx, txOptions}

	// Record call args
	mmBeginTx.BeginTxMock.mutex.Lock()
	mmBeginTx.BeginTxMock.callArgs = append(mmBeginTx.BeginTxMock.callArgs, &mm_params)
	mmBeginTx.BeginTxMock.mutex.Unlock()

	for _, e := range mmBeginTx.BeginTxMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.t1, e.results.err
		}
	}

	if mmBeginTx.BeginTxMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmBeginTx.BeginTxMock.defaultExpectation.Counter, 1)
		mm_want := mmBeginTx.BeginTxMock.defaultExpectation.params
		mm_want_ptrs := mmBeginTx.BeginTxMock.defaultExpectation.paramPtrs

		mm_got := TransactorMockBeginTxParams{ctx, txOptions}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmBeginTx.t.Errorf("TransactorMock.BeginTx got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmBeginTx.BeginTxMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.txOptions != nil && !minimock.Equal(*mm_want_ptrs.txOptions, mm_got.txOptions) {
				mmBeginTx.t.Errorf("TransactorMock.BeginTx got unexpected parameter txOptions, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmBeginTx.BeginTxMock.defaultExpectation.expectationOrigins.originTxOptions, *mm_want_ptrs.txOptions, mm_got.txOptions, minimock.Diff(*mm_want_ptrs.txOptions, mm_got.txOptions))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmBeginTx.t.Errorf("TransactorMock.BeginTx got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmBeginTx.BeginTxMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmBeginTx.BeginTxMock.defaultExpectation.results
		if mm_results == nil {
			mmBeginTx.t.Fatal("No results are set for the TransactorMock.BeginTx")
		}
		return (*mm_results).t1, (*mm_results).err
	}
	if mmBeginTx.funcBeginTx != nil {
		return mmBeginTx.funcBeginTx(ctx, txOptions)
	}
	mmBeginTx.t.Fatalf("Unexpected call to TransactorMock.BeginTx. %v %v", ctx, txOptions)
	return
}

// BeginTxAfterCounter returns a count of finished TransactorMock.BeginTx invocations
func (mmBeginTx *TransactorMock) BeginTxAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.afterBeginTxCounter)
}

// BeginTxBeforeCounter returns a count of TransactorMock.BeginTx invocations
func (mmBeginTx *TransactorMock) BeginTxBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.beforeBeginTxCounter)
}

// Calls returns a list of arguments used in each call to TransactorMock.BeginTx.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmBeginTx *mTransactorMockBeginTx) Calls() []*TransactorMockBeginTxParams {
	mmBeginTx.mutex.RLock()

	argCopy := make([]*TransactorMockBeginTxParams, len(mmBeginTx.callArgs))
	copy(argCopy, mmBeginTx.callArgs)

	mmBeginTx.mutex.RUnlock()

	return argCopy
}

// MinimockBeginTxDone returns true if the count of the BeginTx invocations corresponds
// the number of defined expectations
func (m *TransactorMock) MinimockBeginTxDone() bool {
	if m.BeginTxMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.BeginTxMock.invocationsDone()
}

// MinimockBeginTxInspect logs each unmet expectation
func (m *TransactorMock) MinimockBeginTxInspect() {
	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactorMock.BeginTx at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterBeginTxCounter := mm_atomic.LoadUint64(&m.afterBeginTxCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.BeginTxMock.defaultExpectation != nil && afterBeginTxCounter < 1 {
		if m.BeginTxMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to TransactorMock.BeginTx at\n%s", m.BeginTxMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to TransactorMock.BeginTx at\n%s with params: %#v", m.BeginTxMock.defaultExpectation.expectationOrigins.origin, *m.BeginTxMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBeginTx != nil && afterBeginTxCounter < 1 {
		m.t.Errorf("Expected call to TransactorMock.BeginTx at\n%s", m.funcBeginTxOrigin)
	}

	if !m.BeginTxMock.invocationsDone() && afterBeginTxCounter > 0 {
		m.t.Errorf("Expected %d calls to TransactorMock.BeginTx at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.BeginTxMock.expectedInvocations), m.BeginTxMock.expectedInvocationsOrigin, afterBeginTxCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TransactorMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockBeginTxInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TransactorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *TransactorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockBeginTxDone()
}

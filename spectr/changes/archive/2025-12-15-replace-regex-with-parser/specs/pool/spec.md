## ADDED Requirements

### Requirement: Token Pool
The system SHALL provide an object pool for Token instances to reduce allocation pressure during lexing.

#### Scenario: Token pool structure
- **WHEN** a token pool is created
- **THEN** it SHALL use `sync.Pool` for thread-safe pooling
- **AND** each pool instance SHALL be local to a parse operation

#### Scenario: Token acquisition from pool
- **WHEN** `pool.GetToken()` is called
- **THEN** it SHALL return a Token from the pool if available
- **AND** it SHALL allocate a new Token if pool is empty
- **AND** the returned Token SHALL have zero values (ready for initialization)

#### Scenario: Token return to pool
- **WHEN** `pool.PutToken(t *Token)` is called
- **THEN** the Token SHALL be returned to the pool for reuse
- **AND** Token fields SHALL be cleared before return
- **AND** Source slice reference SHALL be nilled to avoid memory leaks

#### Scenario: Token pool per-parse isolation
- **WHEN** multiple Parse calls run concurrently
- **THEN** each SHALL use its own pool or share safely via sync.Pool
- **AND** tokens SHALL NOT leak between parse operations

### Requirement: Node Pool
The system SHALL provide object pools for AST Node instances to reduce GC pressure.

#### Scenario: Typed node pools
- **WHEN** node pools are created
- **THEN** there SHALL be a separate pool for each node type
- **AND** `pool.GetSection()`, `pool.GetRequirement()`, etc. SHALL return typed nodes

#### Scenario: Node acquisition from pool
- **WHEN** `pool.GetNode[T]()` is called
- **THEN** it SHALL return a *T from the appropriate pool
- **AND** it SHALL allocate a new *T if pool is empty
- **AND** the returned node SHALL be in a clean state

#### Scenario: Node return to pool
- **WHEN** `pool.PutNode(n Node)` is called
- **THEN** the node SHALL be returned to the appropriate typed pool
- **AND** node fields SHALL be cleared (children slice, source, etc.)
- **AND** children SHALL NOT be recursively returned (caller's responsibility)

### Requirement: Children Slice Pool
The system SHALL provide a pool for children slices to reduce slice allocation.

#### Scenario: Children slice acquisition
- **WHEN** `pool.GetChildren(capacity int)` is called
- **THEN** it SHALL return a `[]*Node` slice with at least the requested capacity
- **AND** slice length SHALL be 0 (capacity available)

#### Scenario: Children slice sizing
- **WHEN** pooling children slices
- **THEN** pools SHALL be bucketed by size: small (<=4), medium (<=16), large (<=64)
- **AND** requests larger than 64 SHALL allocate fresh slices (not pooled)

#### Scenario: Children slice return
- **WHEN** `pool.PutChildren(s []*Node)` is called
- **THEN** the slice SHALL be returned to the appropriate bucket
- **AND** slice SHALL be cleared (elements set to nil) before return

### Requirement: Pool Lifecycle Management
The system SHALL manage pool lifecycle to prevent memory leaks.

#### Scenario: Pool scope
- **WHEN** a parse operation begins
- **THEN** pools MAY be created fresh or reused from sync.Pool
- **AND** the parser SHALL track which objects came from pools

#### Scenario: Pool cleanup
- **WHEN** a parse operation completes
- **THEN** temporary objects (tokens) SHALL be returned to pools
- **AND** AST nodes SHALL NOT be returned (they are the result)
- **AND** pool references SHALL be cleared to enable GC

#### Scenario: No pool leaks
- **WHEN** the AST is no longer referenced
- **THEN** pooled allocations behind it SHALL be eligible for GC
- **AND** pools SHALL NOT hold strong references to AST nodes

### Requirement: Pool Statistics
The system SHALL optionally track pool statistics for performance analysis.

#### Scenario: Pool hit/miss tracking
- **WHEN** pool statistics are enabled
- **THEN** `pool.Stats()` SHALL return hit count, miss count, return count
- **AND** statistics SHALL be thread-safe

#### Scenario: Pool statistics disable
- **WHEN** statistics are not needed
- **THEN** a build tag or compile-time flag SHALL disable tracking
- **AND** disabled tracking SHALL have zero runtime overhead

### Requirement: Stateless Parse with Internal Pooling
The system SHALL use pools internally in the stateless Parse function.

#### Scenario: Parse function uses pools
- **WHEN** `Parse(source []byte)` is called
- **THEN** it SHALL internally use pools for tokens and temporary nodes
- **AND** the caller SHALL NOT need to manage pools
- **AND** pools SHALL be accessed from sync.Pool (no explicit creation)

#### Scenario: Pool reuse across calls
- **WHEN** multiple Parse calls are made
- **THEN** sync.Pool MAY reuse allocations across calls
- **AND** this is transparent to the caller
- **AND** no explicit pool passing is required

### Requirement: Builder with Pool Integration
The system SHALL integrate pools with the node builder API.

#### Scenario: Builder acquires from pool
- **WHEN** `NewBuilder[T](pool *Pool)` is called
- **THEN** the builder SHALL acquire its node from the pool
- **AND** if pool is nil, it SHALL allocate normally

#### Scenario: Builder without pool
- **WHEN** building nodes outside parse context
- **THEN** builders SHALL work without pools (pool parameter optional)
- **AND** allocation SHALL fall back to normal `new()`

### Requirement: Pool Thread Safety
The system SHALL ensure pools are safe for concurrent use.

#### Scenario: Concurrent pool access
- **WHEN** multiple goroutines access pools
- **THEN** Get and Put operations SHALL be thread-safe
- **AND** sync.Pool provides this guarantee automatically

#### Scenario: No pool contention
- **WHEN** pools are used during parsing
- **THEN** parsers SHALL use local references to avoid repeated sync.Pool access
- **AND** batch acquisition MAY be used for known allocation patterns

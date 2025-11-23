# Large Test Corpus

## Purpose
This file is designed to stress-test the parser with many requirements and scenarios.

## Requirements

### Requirement: API Gateway Routing
The system SHALL route incoming requests to appropriate microservices.

#### Scenario: Route to user service
- **WHEN** request path starts with /users
- **THEN** route to user service
- **AND** preserve headers and query parameters

#### Scenario: Route to payment service
- **WHEN** request path starts with /payments
- **THEN** route to payment service
- **AND** add service-specific headers

#### Scenario: Invalid route
- **WHEN** request path does not match any service
- **THEN** return 404 Not Found
- **AND** log routing failure

### Requirement: Load Balancing
The system SHALL distribute requests across multiple service instances.

#### Scenario: Round-robin distribution
- **WHEN** multiple instances are available
- **THEN** distribute requests evenly
- **AND** track instance health

#### Scenario: Health check failure
- **WHEN** instance fails health check
- **THEN** remove from rotation
- **AND** retry health check periodically

### Requirement: Request Validation
All incoming requests SHALL be validated before processing.

#### Scenario: Valid request
- **WHEN** request meets all validation rules
- **THEN** allow request to proceed
- **AND** log validation success

#### Scenario: Missing required fields
- **WHEN** request missing required fields
- **THEN** return 400 Bad Request
- **AND** list missing fields in response

#### Scenario: Invalid field types
- **WHEN** request contains invalid field types
- **THEN** return 400 Bad Request
- **AND** specify expected types

### Requirement: Rate Limiting
The system SHALL enforce rate limits per API key.

#### Scenario: Within rate limit
- **WHEN** requests within allowed rate
- **THEN** process all requests normally
- **AND** include rate limit headers

#### Scenario: Exceeded rate limit
- **WHEN** requests exceed rate limit
- **THEN** return 429 Too Many Requests
- **AND** include retry-after header

### Requirement: Authentication
The system SHALL authenticate all requests using API keys or JWT tokens.

#### Scenario: Valid API key
- **WHEN** request includes valid API key
- **THEN** authenticate request
- **AND** attach key metadata to context

#### Scenario: Valid JWT token
- **WHEN** request includes valid JWT
- **THEN** authenticate request
- **AND** extract user claims

#### Scenario: Missing credentials
- **WHEN** request lacks authentication
- **THEN** return 401 Unauthorized
- **AND** include WWW-Authenticate header

### Requirement: Response Caching
The system SHALL cache responses to reduce backend load.

#### Scenario: Cache hit
- **WHEN** cached response exists
- **THEN** return cached data
- **AND** include cache headers

#### Scenario: Cache miss
- **WHEN** no cached response exists
- **THEN** forward to backend
- **AND** cache response if cacheable

#### Scenario: Cache invalidation
- **WHEN** resource is modified
- **THEN** invalidate related cache entries
- **AND** notify cache cluster

### Requirement: Request Tracing
The system SHALL trace requests across all services.

#### Scenario: New trace
- **WHEN** request arrives without trace ID
- **THEN** generate new trace ID
- **AND** attach to all downstream requests

#### Scenario: Continue trace
- **WHEN** request includes trace ID
- **THEN** preserve trace ID
- **AND** add span for gateway processing

### Requirement: Error Handling
The system SHALL handle errors gracefully and provide meaningful responses.

#### Scenario: Backend timeout
- **WHEN** backend service times out
- **THEN** return 504 Gateway Timeout
- **AND** log timeout details

#### Scenario: Backend error
- **WHEN** backend returns 5xx error
- **THEN** return 502 Bad Gateway
- **AND** mask internal error details

### Requirement: Metrics Collection
The system SHALL collect and expose metrics for monitoring.

#### Scenario: Request metrics
- **WHEN** request is processed
- **THEN** record latency and status code
- **AND** increment request counter

#### Scenario: Error metrics
- **WHEN** error occurs
- **THEN** record error type and count
- **AND** update error rate gauge

### Requirement: Circuit Breaking
The system SHALL implement circuit breaker pattern for failing services.

#### Scenario: Circuit closed
- **WHEN** service is healthy
- **THEN** allow all requests through
- **AND** monitor failure rate

#### Scenario: Circuit open
- **WHEN** failure threshold exceeded
- **THEN** return cached responses or errors
- **AND** stop forwarding requests

#### Scenario: Circuit half-open
- **WHEN** recovery timeout expires
- **THEN** allow test requests
- **AND** close circuit if successful

### Requirement: Request Transformation
The system SHALL transform requests between client and backend formats.

#### Scenario: Header transformation
- **WHEN** forwarding to backend
- **THEN** add required headers
- **AND** remove client-specific headers

#### Scenario: Body transformation
- **WHEN** format conversion needed
- **THEN** transform request body
- **AND** update content-type header

### Requirement: WebSocket Support
The system SHALL proxy WebSocket connections.

#### Scenario: WebSocket upgrade
- **WHEN** client requests WebSocket upgrade
- **THEN** forward upgrade to backend
- **AND** maintain connection

#### Scenario: WebSocket message
- **WHEN** message received on WebSocket
- **THEN** forward to backend
- **AND** route response back to client

### Requirement: CORS Handling
The system SHALL handle Cross-Origin Resource Sharing appropriately.

#### Scenario: Preflight request
- **WHEN** OPTIONS request received
- **THEN** return allowed origins and methods
- **AND** include appropriate CORS headers

#### Scenario: Cross-origin request
- **WHEN** request from allowed origin
- **THEN** add CORS headers to response
- **AND** allow request to proceed

### Requirement: Request Logging
The system SHALL log all requests for audit and debugging.

#### Scenario: Successful request logged
- **WHEN** request completes successfully
- **THEN** log request details
- **AND** include response time

#### Scenario: Failed request logged
- **WHEN** request fails
- **THEN** log failure details
- **AND** include error information

### Requirement: SSL Termination
The system SHALL terminate SSL connections and forward to backends.

#### Scenario: HTTPS request
- **WHEN** HTTPS request received
- **THEN** decrypt SSL
- **AND** forward HTTP to backend

#### Scenario: Certificate validation
- **WHEN** checking client certificate
- **THEN** validate certificate chain
- **AND** extract certificate metadata

### Requirement: Request Retry
The system SHALL retry failed requests with exponential backoff.

#### Scenario: Transient failure
- **WHEN** request fails with retriable error
- **THEN** retry with backoff
- **AND** limit retry attempts

#### Scenario: Non-retriable error
- **WHEN** request fails permanently
- **THEN** return error immediately
- **AND** do not retry

### Requirement: Request Queuing
The system SHALL queue requests during peak load.

#### Scenario: Queue request
- **WHEN** all backend workers busy
- **THEN** add request to queue
- **AND** process when worker available

#### Scenario: Queue overflow
- **WHEN** queue is full
- **THEN** return 503 Service Unavailable
- **AND** include retry-after header

### Requirement: Response Compression
The system SHALL compress responses to reduce bandwidth.

#### Scenario: Compress response
- **WHEN** client supports compression
- **THEN** compress response body
- **AND** add content-encoding header

#### Scenario: Skip compression
- **WHEN** response already compressed
- **THEN** forward without compression
- **AND** preserve original headers

### Requirement: IP Whitelisting
The system SHALL restrict access based on IP addresses.

#### Scenario: Whitelisted IP
- **WHEN** request from whitelisted IP
- **THEN** allow request to proceed
- **AND** log IP address

#### Scenario: Non-whitelisted IP
- **WHEN** request from non-whitelisted IP
- **THEN** return 403 Forbidden
- **AND** log blocked attempt

### Requirement: Request Deduplication
The system SHALL detect and handle duplicate requests.

#### Scenario: Duplicate request
- **WHEN** identical request received within window
- **THEN** return cached response
- **AND** do not forward to backend

#### Scenario: Unique request
- **WHEN** request is unique
- **THEN** process normally
- **AND** cache for deduplication

### Requirement: API Versioning
The system SHALL support multiple API versions simultaneously.

#### Scenario: Version in header
- **WHEN** version specified in header
- **THEN** route to appropriate version
- **AND** validate version compatibility

#### Scenario: Version in path
- **WHEN** version in URL path
- **THEN** extract and route to version
- **AND** rewrite path for backend

### Requirement: Request Sanitization
The system SHALL sanitize inputs to prevent injection attacks.

#### Scenario: SQL injection attempt
- **WHEN** request contains SQL patterns
- **THEN** sanitize or reject request
- **AND** log security event

#### Scenario: XSS attempt
- **WHEN** request contains script tags
- **THEN** escape or reject content
- **AND** log security event

### Requirement: Response Filtering
The system SHALL filter sensitive data from responses.

#### Scenario: Filter PII
- **WHEN** response contains personal data
- **THEN** redact based on permissions
- **AND** log data access

#### Scenario: Filter credentials
- **WHEN** response contains credentials
- **THEN** remove from response
- **AND** log security warning

### Requirement: Request Batching
The system SHALL batch multiple requests for efficiency.

#### Scenario: Batch requests
- **WHEN** multiple requests for same service
- **THEN** combine into single batch
- **AND** distribute responses

#### Scenario: Batch timeout
- **WHEN** batch wait time expires
- **THEN** process accumulated requests
- **AND** reset batch window

### Requirement: GraphQL Support
The system SHALL proxy GraphQL queries to backend services.

#### Scenario: GraphQL query
- **WHEN** GraphQL query received
- **THEN** parse and validate query
- **AND** forward to GraphQL backend

#### Scenario: GraphQL mutation
- **WHEN** mutation received
- **THEN** validate mutation schema
- **AND** forward to appropriate service

### Requirement: Request Prioritization
The system SHALL prioritize requests based on business rules.

#### Scenario: High priority request
- **WHEN** request marked as high priority
- **THEN** process before normal requests
- **AND** allocate more resources

#### Scenario: Low priority request
- **WHEN** system under load
- **THEN** delay low priority requests
- **AND** preserve ordering within priority

### Requirement: Service Discovery
The system SHALL discover backend services dynamically.

#### Scenario: Service registration
- **WHEN** new service instance starts
- **THEN** register with discovery
- **AND** add to routing table

#### Scenario: Service deregistration
- **WHEN** service instance stops
- **THEN** remove from discovery
- **AND** drain existing connections

### Requirement: Traffic Splitting
The system SHALL split traffic for A/B testing.

#### Scenario: Canary deployment
- **WHEN** new version deployed
- **THEN** route percentage to new version
- **AND** monitor metrics

#### Scenario: Feature flag routing
- **WHEN** feature flag enabled
- **THEN** route to feature backend
- **AND** track feature usage

### Requirement: Request Mocking
The system SHALL mock responses for testing.

#### Scenario: Mock enabled
- **WHEN** mock mode is active
- **THEN** return mock response
- **AND** do not call backend

#### Scenario: Selective mocking
- **WHEN** specific endpoints mocked
- **THEN** mock those endpoints only
- **AND** forward others to backend

### Requirement: Response Aggregation
The system SHALL aggregate responses from multiple services.

#### Scenario: Parallel requests
- **WHEN** aggregating multiple sources
- **THEN** request in parallel
- **AND** combine responses

#### Scenario: Sequential dependencies
- **WHEN** responses have dependencies
- **THEN** request in sequence
- **AND** pass data between calls

### Requirement: Request Validation Schema
The system SHALL validate requests against OpenAPI schemas.

#### Scenario: Schema validation success
- **WHEN** request matches schema
- **THEN** allow request to proceed
- **AND** log validation success

#### Scenario: Schema validation failure
- **WHEN** request violates schema
- **THEN** return detailed error
- **AND** reference schema violation

### Requirement: Connection Pooling
The system SHALL pool connections to backend services.

#### Scenario: Reuse connection
- **WHEN** connection available in pool
- **THEN** reuse existing connection
- **AND** avoid connection overhead

#### Scenario: Create connection
- **WHEN** pool is empty
- **THEN** create new connection
- **AND** add to pool after use

### Requirement: Request Timeout
The system SHALL enforce timeouts on all requests.

#### Scenario: Request within timeout
- **WHEN** response received before timeout
- **THEN** return response normally
- **AND** update latency metrics

#### Scenario: Request exceeds timeout
- **WHEN** timeout expires
- **THEN** cancel request
- **AND** return timeout error

### Requirement: Health Check Endpoint
The system SHALL provide health check endpoint.

#### Scenario: Healthy system
- **WHEN** all services operational
- **THEN** return 200 OK
- **AND** include service status

#### Scenario: Degraded system
- **WHEN** some services down
- **THEN** return 503
- **AND** list unavailable services

### Requirement: Admin API
The system SHALL provide admin endpoints for management.

#### Scenario: Admin authentication
- **WHEN** accessing admin endpoint
- **THEN** require admin credentials
- **AND** log admin access

#### Scenario: Configuration update
- **WHEN** admin updates configuration
- **THEN** validate changes
- **AND** apply without restart

### Requirement: Request Debugging
The system SHALL support request debugging for troubleshooting.

#### Scenario: Debug mode enabled
- **WHEN** debug header present
- **THEN** include debug information
- **AND** log detailed trace

#### Scenario: Debug mode disabled
- **WHEN** production mode
- **THEN** omit debug data
- **AND** use standard logging

### Requirement: Response Headers
The system SHALL add standard response headers.

#### Scenario: Security headers
- **WHEN** responding to request
- **THEN** add security headers
- **AND** include HSTS and CSP

#### Scenario: Custom headers
- **WHEN** service specifies headers
- **THEN** merge with defaults
- **AND** avoid header conflicts

### Requirement: Request Size Limits
The system SHALL enforce request size limits.

#### Scenario: Within size limit
- **WHEN** request size acceptable
- **THEN** process request normally
- **AND** track request size

#### Scenario: Exceeds size limit
- **WHEN** request too large
- **THEN** return 413 Payload Too Large
- **AND** log oversized request

### Requirement: Graceful Shutdown
The system SHALL shutdown gracefully without dropping requests.

#### Scenario: Shutdown initiated
- **WHEN** shutdown signal received
- **THEN** stop accepting new requests
- **AND** complete in-flight requests

#### Scenario: Shutdown timeout
- **WHEN** shutdown timeout expires
- **THEN** force close connections
- **AND** log incomplete requests

### Requirement: Multi-Region Support
The system SHALL route requests to nearest region.

#### Scenario: Regional routing
- **WHEN** request from specific region
- **THEN** route to nearest datacenter
- **AND** fallback to other regions

#### Scenario: Region failover
- **WHEN** regional service unavailable
- **THEN** failover to backup region
- **AND** maintain session state

### Requirement: Request Analytics
The system SHALL collect analytics on request patterns.

#### Scenario: Usage analytics
- **WHEN** processing requests
- **THEN** track endpoint usage
- **AND** identify popular routes

#### Scenario: Performance analytics
- **WHEN** measuring performance
- **THEN** track p50, p95, p99 latencies
- **AND** identify slow endpoints

### Requirement: Webhook Support
The system SHALL forward webhook events to registered handlers.

#### Scenario: Webhook delivery
- **WHEN** webhook event received
- **THEN** forward to registered URL
- **AND** retry on failure

#### Scenario: Webhook validation
- **WHEN** webhook includes signature
- **THEN** validate signature
- **AND** reject if invalid

### Requirement: Request Throttling
The system SHALL throttle requests from abusive clients.

#### Scenario: Normal usage
- **WHEN** client behaves normally
- **THEN** process all requests
- **AND** track usage patterns

#### Scenario: Abusive pattern detected
- **WHEN** abuse threshold exceeded
- **THEN** throttle client requests
- **AND** notify administrators

### Requirement: Service Mesh Integration
The system SHALL integrate with service mesh for advanced routing.

#### Scenario: Mesh routing
- **WHEN** service mesh enabled
- **THEN** use mesh routing rules
- **AND** delegate traffic management

#### Scenario: Mesh observability
- **WHEN** mesh active
- **THEN** export traces to mesh
- **AND** use mesh metrics

### Requirement: Custom Middleware
The system SHALL support custom middleware plugins.

#### Scenario: Middleware execution
- **WHEN** request processed
- **THEN** execute middleware chain
- **AND** preserve execution order

#### Scenario: Middleware error
- **WHEN** middleware fails
- **THEN** handle error appropriately
- **AND** continue or abort based on config

### Requirement: Request Rewriting
The system SHALL rewrite requests for backend compatibility.

#### Scenario: Path rewriting
- **WHEN** backend expects different path
- **THEN** rewrite path before forwarding
- **AND** preserve query parameters

#### Scenario: Header rewriting
- **WHEN** backend needs specific headers
- **THEN** add or modify headers
- **AND** document transformations

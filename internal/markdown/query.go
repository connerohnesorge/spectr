package markdown

// Predicate is a function type for node matching predicates.
type Predicate func(Node) bool

// Find returns all nodes in the AST where pred returns true.
// Results are in pre-order traversal order (depth-first).
// Returns an empty slice (not nil) if no nodes match.
func Find(root Node, pred Predicate) []Node {
	var results []Node
	if root == nil {
		return make([]Node, 0)
	}
	findRecursive(root, pred, &results)

	return results
}

func findRecursive(
	node Node,
	pred Predicate,
	results *[]Node,
) {
	if node == nil {
		return
	}
	if pred(node) {
		*results = append(*results, node)
	}
	for _, child := range node.Children() {
		findRecursive(child, pred, results)
	}
}

// FindFirst returns the first node where pred returns true, or nil if none.
// Traversal stops immediately upon finding a match (short-circuit).
func FindFirst(root Node, pred Predicate) Node {
	if root == nil {
		return nil
	}

	return findFirstRecursive(root, pred)
}

func findFirstRecursive(
	node Node,
	pred Predicate,
) Node {
	if node == nil {
		return nil
	}
	if pred(node) {
		return node
	}
	for _, child := range node.Children() {
		if found := findFirstRecursive(child, pred); found != nil {
			return found
		}
	}

	return nil
}

// FindByType returns all nodes of type T in the AST.
// Results are properly typed and in traversal order.
// The type parameter T must implement the Node interface.
func FindByType[T Node](root Node) []T {
	var results []T
	if root == nil {
		return results
	}
	findByTypeRecursive[T](root, &results)

	return results
}

func findByTypeRecursive[T Node](
	node Node,
	results *[]T,
) {
	if node == nil {
		return
	}
	if typed, ok := node.(T); ok {
		*results = append(*results, typed)
	}
	for _, child := range node.Children() {
		findByTypeRecursive[T](child, results)
	}
}

// FindFirstByType returns the first node of type T, or nil if none exists.
// The type parameter T must implement the Node interface.
func FindFirstByType[T Node](root Node) T {
	var zero T
	if root == nil {
		return zero
	}
	result, _ := findFirstByTypeRecursive[T](root)

	return result
}

func findFirstByTypeRecursive[T Node](
	node Node,
) (result T, found bool) {
	if node == nil {
		return result, false
	}
	if typed, ok := node.(T); ok {
		return typed, true
	}
	for _, child := range node.Children() {
		if result, found = findFirstByTypeRecursive[T](child); found {
			return result, true
		}
	}

	return result, false
}

// And returns a predicate that is true when both p1 AND p2 are true.
// Uses short-circuit evaluation: p2 is not called if p1 is false.
func And(p1, p2 Predicate) Predicate {
	return func(n Node) bool {
		return p1(n) && p2(n)
	}
}

// Or returns a predicate that is true when p1 OR p2 is true.
// Uses short-circuit evaluation: p2 is not called if p1 is true.
func Or(p1, p2 Predicate) Predicate {
	return func(n Node) bool {
		return p1(n) || p2(n)
	}
}

// Not returns a predicate that negates p.
// Not(p)(node) equals !p(node).
func Not(p Predicate) Predicate {
	return func(n Node) bool {
		return !p(n)
	}
}

// All returns a predicate that is true when all preds are true.
// Uses short-circuit evaluation: stops on first false.
// Returns true for empty predicate list.
func All(preds ...Predicate) Predicate {
	return func(n Node) bool {
		for _, p := range preds {
			if !p(n) {
				return false
			}
		}

		return true
	}
}

// Any returns a predicate that is true when any pred is true.
// Uses short-circuit evaluation: stops on first true.
// Returns false for empty predicate list.
func Any(preds ...Predicate) Predicate {
	return func(n Node) bool {
		for _, p := range preds {
			if p(n) {
				return true
			}
		}

		return false
	}
}

// IsType returns a predicate that matches nodes of type T.
// Type checking uses Go type assertion.
func IsType[T Node]() Predicate {
	return func(n Node) bool {
		_, ok := n.(T)

		return ok
	}
}

// Named is an interface for nodes that have a Name() method.
// This includes NodeRequirement and NodeScenario.
type Named interface {
	// Name returns the name of the node.
	Name() string
}

// HasName returns a predicate matching nodes with Name() == name.
// Works for NodeRequirement, NodeScenario, and any type implementing Named.
func HasName(name string) Predicate {
	return func(n Node) bool {
		if named, ok := n.(Named); ok {
			return named.Name() == name
		}

		return false
	}
}

// InRange returns a predicate matching nodes within [start, end).
// A node is in range if its span overlaps the given range.
// Overlap occurs when node.start < end AND node.end > start.
func InRange(start, end int) Predicate {
	return func(n Node) bool {
		nodeStart, nodeEnd := n.Span()
		// Check for overlap: node.start < end AND node.end > start
		return nodeStart < end && nodeEnd > start
	}
}

// HasChild returns a predicate that is true if any direct child matches pred.
// Only immediate children are checked (not descendants).
func HasChild(pred Predicate) Predicate {
	return func(n Node) bool {
		for _, child := range n.Children() {
			if pred(child) {
				return true
			}
		}

		return false
	}
}

// HasDescendant returns a predicate true if any descendant matches pred.
// All descendants are checked recursively.
func HasDescendant(pred Predicate) Predicate {
	return func(n Node) bool {
		return hasDescendantRecursive(n, pred)
	}
}

func hasDescendantRecursive(
	node Node,
	pred Predicate,
) bool {
	for _, child := range node.Children() {
		if pred(child) {
			return true
		}
		if hasDescendantRecursive(child, pred) {
			return true
		}
	}

	return false
}

// Count returns the number of nodes matching pred without allocating a slice.
// Traverses the entire tree and increments a counter for each match.
func Count(root Node, pred Predicate) int {
	if root == nil {
		return 0
	}

	return countRecursive(root, pred)
}

func countRecursive(
	node Node,
	pred Predicate,
) int {
	count := 0
	if node == nil {
		return 0
	}
	if pred(node) {
		count = 1
	}
	for _, child := range node.Children() {
		count += countRecursive(child, pred)
	}

	return count
}

// Exists returns true if any node matches pred.
// Uses short-circuit evaluation: stops on first match.
func Exists(root Node, pred Predicate) bool {
	if root == nil {
		return false
	}

	return existsRecursive(root, pred)
}

func existsRecursive(
	node Node,
	pred Predicate,
) bool {
	if node == nil {
		return false
	}
	if pred(node) {
		return true
	}
	for _, child := range node.Children() {
		if existsRecursive(child, pred) {
			return true
		}
	}

	return false
}

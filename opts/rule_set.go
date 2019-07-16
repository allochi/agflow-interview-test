package opts

// RuleSet a data structure that hold the rules between options
type RuleSet struct {
	edges      map[string][]string
	dependents map[string][]string
	conflicts  map[string][]string
}

// NewRuleSet create a bew RuleSet
func NewRuleSet() *RuleSet {
	return &RuleSet{
		edges:      make(map[string][]string),
		dependents: make(map[string][]string),
		conflicts:  make(map[string][]string),
	}
}

// AddDep adds a dependency between options
func (rs *RuleSet) AddDep(a, b string) {
	rs.dependents[a] = append(rs.dependents[a], b)
	rs.edges[a] = append(rs.edges[a], b)
	rs.edges[b] = append(rs.edges[b], a)
}

// AddConflict adds a conflict between options
func (rs *RuleSet) AddConflict(a, b string) {
	rs.conflicts[a] = append(rs.conflicts[a], b)
	rs.conflicts[b] = append(rs.conflicts[b], a)
}

// IsCoherent checks if the rule set is coherent
func (rs *RuleSet) IsCoherent() bool {
	// walk the graph
	for start, ends := range rs.conflicts {
		for _, end := range ends {
			// check if start and end of a conflict rule are dependents
			canWalk := rs.canWalkBetween(start, end)
			if canWalk {
				return false
			}
		}
	}

	return true
}

// canWalkBetween check if there is a connection between two options in the graph
func (rs *RuleSet) canWalkBetween(a, b string) bool {
	visited := make(map[string]bool)
	queue := Queue{}
	visited[a] = true
	queue.Push(rs.edges[a]...)

	for queue.Len() > 0 {
		node := queue.Pop()
		if node == b {
			return true
		}

		if !visited[node] {
			queue.Push(rs.edges[node]...)
			visited[node] = true
		}
	}

	return false
}

// Opts option type
type Opts struct {
	rs      *RuleSet
	options map[string]bool
}

// New returns a new (empty) collection of selected options (Opts) for the rule set rs
func New(rs *RuleSet) *Opts {
	opts := Opts{
		rs:      rs,
		options: make(map[string]bool),
	}
	opts.updateOptions()
	return &opts
}

// updateOptions update the set of options from RuleSet graph
// this method is used to update the options set before operations on Opts
func (opts *Opts) updateOptions() {
	for option := range opts.rs.edges {
		if _, ok := opts.options[option]; !ok {
			// by default the option state is false
			opts.options[option] = false
			// if a parent exist the new option takes its state
			if len(opts.rs.edges[option]) > 0 {
				parent := opts.rs.edges[option][0] // TODO: what if there are multiple parents?
				opts.options[option] = opts.options[parent]
			}
		}
	}

	for option := range opts.rs.conflicts {
		if _, ok := opts.options[option]; !ok {
			opts.options[option] = false
		}
	}
}

// Toggle method to set or unset an option
func (opts *Opts) Toggle(opt string) {
	opts.updateOptions()
	visited := make(map[string]bool)
	opts.toggle(opt, visited)
}

func (opts *Opts) toggle(opt string, visited map[string]bool) {
	// skip toggle if visited
	if visited[opt] {
		return
	}

	// toggle the option and mark it as visited
	opts.options[opt] = !opts.options[opt]
	visited[opt] = true

	// propagate through conflicts
	for _, option := range opts.rs.conflicts[opt] {
		// options are mutually exclusive in conflict rule
		if opts.options[option] && opts.options[opt] {
			opts.Toggle(option)
		}
	}

	// propagate through all dependents
	for _, option := range opts.rs.dependents[opt] {
		opts.toggle(option, visited)
	}
}

// StringSlice returns a slice of string with the current list of selected options
func (opts *Opts) StringSlice() (sl []string) {
	opts.updateOptions()
	for option, ok := range opts.options {
		if ok {
			sl = append(sl, option)
		}
	}
	return
}

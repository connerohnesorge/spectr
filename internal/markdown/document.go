package markdown

// Document query methods

// GetSection returns a section by name.
func (d *Document) GetSection(name string) *Section {
	return d.Sections[name]
}

// GetSectionContent returns the content of a section by name.
func (d *Document) GetSectionContent(name string) string {
	if s := d.Sections[name]; s != nil {
		return s.Content
	}

	return ""
}

// GetDeltaSection returns a delta section by type (ADDED, MODIFIED, etc.).
func (d *Document) GetDeltaSection(deltaType string) *Section {
	return d.Sections[deltaType+" Requirements"]
}

// GetRequirement returns a requirement by name.
func (d *Document) GetRequirement(name string) *Requirement {
	return d.Requirements[name]
}

// GetRequirementNames returns all requirement names in document order.
func (d *Document) GetRequirementNames() []string {
	names := make([]string, 0, len(d.H3Headers))

	for _, h := range d.H3Headers {
		if name := parseRequirementName(h.Text); name != "" {
			names = append(names, name)
		}
	}

	return names
}

// GetScenario returns a scenario by name.
func (d *Document) GetScenario(name string) *Scenario {
	return d.Scenarios[name]
}

// GetHeadersByLevel returns all headers at a specific level.
func (d *Document) GetHeadersByLevel(level int) []Header {
	switch level {
	case headerLevelH2:
		return d.H2Headers
	case headerLevelH3:
		return d.H3Headers
	case headerLevelH4:
		return d.H4Headers
	default:
		var headers []Header
		for _, h := range d.Headers {
			if h.Level == level {
				headers = append(headers, h)
			}
		}

		return headers
	}
}

// GetAllTasks returns a flat list of all tasks including nested ones.
func (d *Document) GetAllTasks() []Task {
	var result []Task
	var flatten func(tasks []Task)

	flatten = func(tasks []Task) {
		for _, t := range tasks {
			result = append(result, t)
			if len(t.Children) > 0 {
				flatten(t.Children)
			}
		}
	}

	flatten(d.Tasks)

	return result
}

// CountTasks returns total and completed task counts.
func (d *Document) CountTasks() (total, completed int) {
	all := d.GetAllTasks()
	total = len(all)

	for _, t := range all {
		if t.Checked {
			completed++
		}
	}

	return total, completed
}

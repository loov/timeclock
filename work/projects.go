package work

type Projects struct {
	Active map[ProjectID]*Project
}

func (projects *Projects) Lookup(pid ProjectID) (*Project, bool) {
	p, ok := projects.Active[pid]
	return p, ok
}

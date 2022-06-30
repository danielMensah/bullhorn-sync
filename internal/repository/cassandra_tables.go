package repository

import (
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	candidateMetadata = table.Metadata{
		Name:    "candidates",
		Columns: []string{"first_name", "last_name"},
		PartKey: []string{"first_name"},
	}

	CandidateTable = table.New(candidateMetadata)
)

type Candidate struct {
	FirstName string
	LastName  string
}

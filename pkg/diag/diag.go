// Copyright 2016 Marapongo, Inc. All rights reserved.

package diag

// ID is a unique diagnostics identifier.
type ID int

// Diag is an instance of an error or warning generated by the compiler.
type Diag struct {
	ID      ID        // a unique identifier for this diagnostic.
	Message string    // a human-friendly message for this diagnostic.
	Doc     *Document // the document in which this diagnostic occurred.
	Loc     *Location // the document location at which this diagnostic occurred.
}

// WithFile adds a file to an existing diagnostic, retaining its ID and message.
func (diag *Diag) WithFile(file string) *Diag {
	return &Diag{
		ID:      diag.ID,
		Message: diag.Message,
		Doc:     NewDocument(file),
		Loc:     &EmptyLocation,
	}
}

// WithDocument adds a file to an existing diagnostic, retaining its ID and message.
func (diag *Diag) WithDocument(doc *Document) *Diag {
	return &Diag{
		ID:      diag.ID,
		Message: diag.Message,
		Doc:     doc,
		Loc:     &EmptyLocation,
	}
}

// WithLocation adds a file and position to an existing diagnostic, retaining its ID and message.
func (diag *Diag) WithLocation(doc *Document, loc *Location) *Diag {
	return &Diag{
		ID:      diag.ID,
		Message: diag.Message,
		Doc:     doc,
		Loc:     loc,
	}
}

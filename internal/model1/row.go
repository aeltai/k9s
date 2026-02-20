// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package model1

// Row represents a collection of columns.
type Row struct {
	ID     string
	Fields Fields
}

// NewRow returns a new row with initialized fields.
func NewRow(size int) Row {
	return Row{Fields: make([]string, size)}
}

// Labelize returns a new row based on labels.
func (r Row) Labelize(cols []int, labelCol int, labels []string) Row {
	out := NewRow(len(cols) + len(labels))
	for _, col := range cols {
		out.Fields = append(out.Fields, r.Fields[col])
	}
	m := labelize(r.Fields[labelCol])
	for _, label := range labels {
		out.Fields = append(out.Fields, m[label])
	}

	return out
}

// Customize returns a row subset based on given col indices.
func (r Row) Customize(cols []int) Row {
	out := NewRow(len(cols))
	r.Fields.Customize(cols, out.Fields)
	out.ID = r.ID

	return out
}

// Diff returns true if row differ or false otherwise.
func (r Row) Diff(ro Row, ageCol int) bool {
	if r.ID != ro.ID {
		return true
	}
	return r.Fields.Diff(ro.Fields, ageCol)
}

// Clone copies a row.
func (r Row) Clone() Row {
	return Row{
		ID:     r.ID,
		Fields: r.Fields.Clone(),
	}
}

// Len returns the length of the row.
func (r Row) Len() int {
	return len(r.Fields)
}

// ----------------------------------------------------------------------------

// RowSorter sorts rows.
type RowSorter struct {
	Rows       Rows
	Index      int
	IsNumber   bool
	IsDuration bool
	IsCapacity bool
	Asc        bool
}

func (s RowSorter) Len() int {
	return len(s.Rows)
}

func (s RowSorter) Swap(i, j int) {
	s.Rows[i], s.Rows[j] = s.Rows[j], s.Rows[i]
}

func (s RowSorter) Less(i, j int) bool {
	v1, v2 := s.Rows[i].Fields[s.Index], s.Rows[j].Fields[s.Index]
	id1, id2 := s.Rows[i].ID, s.Rows[j].ID
	less := Less(s.IsNumber, s.IsDuration, s.IsCapacity, id1, id2, v1, v2)
	if s.Asc {
		return less
	}
	return !less
}

// ----------------------------------------------------------------------------

// MultiContextSep separates the context name from the resource path in row IDs.
const MultiContextSep = "@@"

// SplitMultiContextID splits a multi-context row ID into context and path.
// Returns ("", id) if not a multi-context ID.
func SplitMultiContextID(id string) (ctx, path string) {
	if i := indexOf(id, MultiContextSep); i >= 0 {
		return id[:i], id[i+len(MultiContextSep):]
	}
	return "", id
}

// JoinMultiContextID creates a multi-context row ID.
func JoinMultiContextID(ctx, path string) string {
	return ctx + MultiContextSep + path
}

func indexOf(s, sep string) int {
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}

// ----------------------------------------------------------------------------
// Helpers...

// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The MIT License (MIT)

// Copyright (c) 2016 Jerry Bai

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sqlparser

// Insert represents an INSERT statement.
type Insert struct {
	Comments Comments
	Ignore   string
	Table    *TableName
	Columns  Columns
	Rows     InsertRows
	OnDup    OnDup
}

func (node *Insert) Format(buf *TrackedBuffer) {
	buf.Fprintf("insert %v%s into %v%v %v%v",
		node.Comments, node.Ignore,
		node.Table, node.Columns, node.Rows, node.OnDup)
}
func (node *Insert) IStatement() {}

// Update represents an UPDATE statement.
type Update struct {
	Comments Comments
	Table    *TableName
	Exprs    UpdateExprs
	Where    *Where
	OrderBy  OrderBy
	Limit    *Limit
}

func (node *Update) Format(buf *TrackedBuffer) {
	buf.Fprintf("update %v%v set %v%v%v%v",
		node.Comments, node.Table,
		node.Exprs, node.Where, node.OrderBy, node.Limit)
}

func (node *Update) IStatement() {}

// Delete represents a DELETE statement.
type Delete struct {
	Comments Comments
	Table    *TableName
	Where    *Where
	OrderBy  OrderBy
	Limit    *Limit
}

func (node *Delete) Format(buf *TrackedBuffer) {
	buf.Fprintf("delete %vfrom %v%v%v%v",
		node.Comments,
		node.Table, node.Where, node.OrderBy, node.Limit)
}

func (node *Delete) IStatement() {}

// Replace represents an REPLACE statement.
type Replace struct {
	Comments Comments
	Table    *TableName
	Columns  Columns
	Rows     InsertRows
}

func (node *Replace) Format(buf *TrackedBuffer) {
	buf.Fprintf("replace %vinto %v%v %v",
		node.Comments,
		node.Table, node.Columns, node.Rows)
}

func (node *Replace) IStatement() {}

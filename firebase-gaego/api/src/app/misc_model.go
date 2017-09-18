package app

import (
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

// Noop is no-op.
type Noop struct {
}

// IntIDReq is int id request for in query value.
type IntIDReq struct {
	ID int64 `json:"id,string" swagger:",in=query"`
}

// IntIDInPathReq is int id request for in path value.
type IntIDInPathReq struct {
	ID int64 `json:"id,string" swagger:",in=path"`
}

// StringIDReq is string id request for in query value.
type StringIDReq struct {
	ID string `json:"id" swagger:",in=query"`
}

// ReqListOptions is request that use for list query options.
type ReqListOptions struct {
	Limit  int    `json:"limit" swagger:",in=query,d=10"`
	Offset int    `json:"offset" swagger:",in=query"`
	Cursor string `json:"cursor" swagger:",in=query"`
}

// RespListOptions is response of list query.
type RespListOptions struct {
	Cursor string `json:"cursor,omitempty" swagger:",in=query"`
}

// QueryListLoader interface is automagically query loader. see ExecQuery.
type QueryListLoader interface {
	LoadEntity(g *goon.Goon, key *datastore.Key) (interface{}, error)
	Append(v interface{}) error
	PostProcess(g *goon.Goon) error
	ReqListOptions() *ReqListOptions
	RespListOptions() *RespListOptions
}

// ExecQuery implements general query method for Datastore.
func ExecQuery(g *goon.Goon, q *datastore.Query, ldr QueryListLoader) error {
	req := ldr.ReqListOptions()

	if req.Limit == 0 {
		req.Limit = 10
	} else if req.Limit != -1 {
		q = q.Limit(req.Limit + 1) // get one more! it uses for `has next`.
	}

	if req.Cursor != "" {
		cursor, err := datastore.DecodeCursor(req.Cursor)
		if err != nil {
			return err
		}
		q = q.Start(cursor)
	}

	q = q.KeysOnly() // use PostProcess

	t := g.Run(q)

	count := 0
	hasNext := false
	var cursor datastore.Cursor
	for {
		key, err := t.Next(nil)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return err
		}
		count++
		if req.Limit != -1 && req.Limit < count {
			hasNext = true
			break
		}
		inst, err := ldr.LoadEntity(g, key)
		if err != nil {
			return err
		}
		err = ldr.Append(inst)
		if err != nil {
			return err
		}
		if req.Limit == count {
			cursor, err = t.Cursor()
			if err != nil {
				return err
			}
		}
	}

	err := ldr.PostProcess(g)
	if err != nil {
		return err
	}

	resp := ldr.RespListOptions()

	if hasNext {
		resp.Cursor = cursor.String()
	}

	return nil
}

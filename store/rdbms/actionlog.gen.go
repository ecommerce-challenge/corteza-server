package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/actionlog.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
)

// SearchActionlogs returns all matching rows
//
// This function calls convertActionlogFilter with the given
// actionlog.Filter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error) {
	q, err := s.convertActionlogFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap := DefaultSliceCapacity

	var (
		set = make([]*actionlog.Action, 0, scap)
		res *actionlog.Action
	)

	return set, f, func() error {
		rows, err := s.Query(ctx, q)
		if err != nil {
			return err
		}

		for rows.Next() {
			if res, err = s.internalActionlogRowScanner(rows, rows.Err()); err != nil {
				if cerr := rows.Close(); cerr != nil {
					return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
				}

				return err
			}

			set = append(set, res)
		}

		return rows.Close()
	}()
}

// CreateActionlog creates one or more rows in actionlog table
func (s Store) CreateActionlog(ctx context.Context, rr ...*actionlog.Action) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ActionlogTable()).SetMap(s.internalActionlogEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateActionlog updates one or more existing rows in actionlog
func (s Store) UpdateActionlog(ctx context.Context, rr ...*actionlog.Action) error {
	return s.PartialUpdateActionlog(ctx, nil, rr...)
}

// PartialUpdateActionlog updates one or more existing rows in actionlog
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateActionlog(ctx context.Context, onlyColumns []string, rr ...*actionlog.Action) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateActionlogs(
				ctx,
				squirrel.Eq{s.preprocessColumn("alg.id", ""): s.preprocessValue(res.ID, "")},
				s.internalActionlogEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveActionlog removes one or more rows from actionlog table
func (s Store) RemoveActionlog(ctx context.Context, rr ...*actionlog.Action) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ActionlogTable("alg")).Where(squirrel.Eq{s.preprocessColumn("alg.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveActionlogByID removes row from the actionlog table
func (s Store) RemoveActionlogByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ActionlogTable("alg")).Where(squirrel.Eq{s.preprocessColumn("alg.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateActionlogs removes all rows from the actionlog table
func (s Store) TruncateActionlogs(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ActionlogTable())
}

// ExecUpdateActionlogs updates all matched (by cnd) rows in actionlog with given data
func (s Store) ExecUpdateActionlogs(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ActionlogTable("alg")).Where(cnd).SetMap(set))
}

// ActionlogLookup prepares Actionlog query and executes it,
// returning actionlog.Action (or error)
func (s Store) ActionlogLookup(ctx context.Context, cnd squirrel.Sqlizer) (*actionlog.Action, error) {
	return s.internalActionlogRowScanner(s.QueryRow(ctx, s.QueryActionlogs().Where(cnd)))
}

func (s Store) internalActionlogRowScanner(row rowScanner, err error) (*actionlog.Action, error) {
	if err != nil {
		return nil, err
	}

	var res = &actionlog.Action{}
	if _, has := s.config.RowScanners["actionlog"]; has {
		scanner := s.config.RowScanners["actionlog"].(func(rowScanner, *actionlog.Action) error)
		err = scanner(row, res)
	} else {
		err = s.scanActionlogRow(row, res)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Actionlog: %w", err)
	} else {
		return res, nil
	}
}

// QueryActionlogs returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryActionlogs() squirrel.SelectBuilder {
	return s.Select(s.ActionlogTable("alg"), s.ActionlogColumns("alg")...)
}

// ActionlogTable name of the db table
func (Store) ActionlogTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "actionlog" + alias
}

// ActionlogColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ActionlogColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "ts",
		alias + "request_origin",
		alias + "request_id",
		alias + "actor_ip_addr",
		alias + "actor_id",
		alias + "resource",
		alias + "action",
		alias + "error",
		alias + "severity",
		alias + "description",
		alias + "meta",
	}
}

// internalActionlogEncoder encodes fields from actionlog.Action to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeActionlog
// func when rdbms.customEncoder=true
func (s Store) internalActionlogEncoder(res *actionlog.Action) store.Payload {
	return s.encodeActionlog(res)
}
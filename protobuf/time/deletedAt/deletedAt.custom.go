package deletedAt

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	timei "github.com/hopeio/cherry/utils/time"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"io"
	"reflect"
	"strings"
	"time"
)

func (t *DeletedAt) Time() time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func (x *DeletedAt) IsValid() bool {
	return x != nil && x.check() == 0
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
)

func (x *DeletedAt) check() uint {
	const minTimestamp = -62135596800  // Seconds between 1970-01-01T00:00:00Z and 0001-01-01T00:00:00Z, inclusive
	const maxTimestamp = +253402300799 // Seconds between 1970-01-01T00:00:00Z and 9999-12-31T23:59:59Z, inclusive
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < minTimestamp:
		return invalidUnderflow
	case secs > maxTimestamp:
		return invalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return invalidNanos
	default:
		return 0
	}
}

// Scan implements the Scanner interface.
func (ts *DeletedAt) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	if err != nil {
		return err
	}
	if nullTime.Valid {
		*ts = DeletedAt{Seconds: nullTime.Time.Unix(), Nanos: int32(nullTime.Time.Nanosecond())}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t *DeletedAt) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return time.Unix(t.Seconds, int64(t.Nanos)), nil
}

func (ts *DeletedAt) GormDataType() string {
	return "time"
}

var (
	FlagDeleted = 1
	FlagActived = 0
)

func (*DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteQueryClause{Field: f}}
}

type SoftDeleteQueryClause struct {
	Field *schema.Field
}

func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

func (sd SoftDeleteQueryClause) Build(clause.Builder) {
}

func (sd SoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok && !stmt.Statement.Unscoped {
		if c, ok := stmt.Clauses["WHERE"]; ok {
			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) >= 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
						c.Expression = where
						stmt.Clauses["WHERE"] = c
						break
					}
				}
			}
		}

		if sd.Field.DefaultValue == "null" {
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
			}})
		} else {
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: FlagActived},
			}})
		}
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

func (*DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	settings := schema.ParseTagSetting(f.TagSettings["SOFTDELETE"], ",")
	softDeleteClause := SoftDeleteDeleteClause{
		Field:    f,
		Flag:     settings["FLAG"] != "",
		TimeType: getTimeType(settings),
	}
	if v := settings["DELETEDATFIELD"]; v != "" { // DeletedAtField
		softDeleteClause.DeleteAtField = f.Schema.LookUpField(v)
	}
	return []clause.Interface{softDeleteClause}
}

func (*DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteUpdateClause{Field: f}}
}

type SoftDeleteUpdateClause struct {
	Field *schema.Field
}

func (sd SoftDeleteUpdateClause) Name() string {
	return ""
}

func (sd SoftDeleteUpdateClause) Build(clause.Builder) {
}

func (sd SoftDeleteUpdateClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		SoftDeleteQueryClause(sd).ModifyStatement(stmt)
	}
}

type SoftDeleteDeleteClause struct {
	Field         *schema.Field
	Flag          bool
	TimeType      schema.TimeType
	DeleteAtField *schema.Field
}

func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

func (sd SoftDeleteDeleteClause) Build(clause.Builder) {
}

func (sd SoftDeleteDeleteClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		var (
			curTime = stmt.DB.NowFunc()
			set     clause.Set
		)

		if deleteAtField := sd.DeleteAtField; deleteAtField != nil {
			var value interface{}
			if deleteAtField.GORMDataType == "time" {
				value = curTime
			} else {
				value = sd.timeToUnix(curTime)
			}
			set = append(set, clause.Assignment{Column: clause.Column{Name: deleteAtField.DBName}, Value: value})
			stmt.SetColumn(deleteAtField.DBName, value, true)
		}

		if sd.Flag {
			set = append(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: FlagDeleted}}, set...)
			stmt.SetColumn(sd.Field.DBName, FlagDeleted, true)
			stmt.AddClause(set)
		} else {
			var curUnix = sd.timeToUnix(curTime)
			set = append(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: curUnix}}, set...)
			stmt.AddClause(set)
			stmt.SetColumn(sd.Field.DBName, curUnix, true)
		}

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(stmt.Context, reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		SoftDeleteQueryClause{Field: sd.Field}.ModifyStatement(stmt)
		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}

func (sd SoftDeleteDeleteClause) timeToUnix(curTime time.Time) int64 {
	switch sd.TimeType {
	case schema.UnixNanosecond:
		return curTime.UnixNano()
	case schema.UnixMillisecond:
		return curTime.UnixMilli()
	default:
		return curTime.Unix()
	}
}

func getTimeType(settings map[string]string) schema.TimeType {
	if settings["NANO"] != "" {
		return schema.UnixNanosecond
	}

	if settings["MILLI"] != "" {
		return schema.UnixMillisecond
	}

	fieldUnit := strings.ToUpper(settings["DELETEDATFIELDUNIT"])

	if fieldUnit == "NANO" {
		return schema.UnixNanosecond
	}

	if fieldUnit == "MILLI" {
		return schema.UnixMillisecond
	}

	return schema.UnixSecond
}

func (x *DeletedAt) MarshalGQL(w io.Writer) {
	data, _ := timei.MarshalText(x.Time())
	w.Write(data)
}

func (x *DeletedAt) UnmarshalGQL(v interface{}) error {
	var t time.Time
	if i, ok := v.(string); ok {
		err := timei.UnmarshalText(&t, []byte(i))
		if err != nil {
			return err
		}
		*x = DeletedAt{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
		return nil
	}
	return errors.New("enum need integer type")
}

type DeletedAtInput = DeletedAt

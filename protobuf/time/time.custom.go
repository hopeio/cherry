package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

func (ts *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Time{T: nullTime.Time.UnixNano()}
	return
}

func (ts Time) Value() (driver.Value, error) {
	return time.Unix(0, ts.T), nil
}

func (ts Time) Format(foramt string) string {
	return time.Unix(0, ts.T).Format(foramt)
}

// GormDataType gorm common data type
func (ts Time) GormDataType() string {
	return "datetime"
}

func (ts Time) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(ts.T >> 56), // bytes 1-8: seconds
		byte(ts.T >> 48),
		byte(ts.T >> 40),
		byte(ts.T >> 32),
		byte(ts.T >> 24),
		byte(ts.T >> 16),
		byte(ts.T >> 8),
		byte(ts.T),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Time) UnmarshalBinary(data []byte) error {
	ts.T = int64(data[7]) | int64(data[6])<<8 | int64(data[5])<<16 | int64(data[4])<<24 |
		int64(data[3])<<32 | int64(data[2])<<40 | int64(data[1])<<48 | int64(data[0])<<56
	return nil
}

func (ts Time) GobEncode() ([]byte, error) {
	return ts.MarshalBinary()
}

func (ts *Time) GobDecode(data []byte) error {
	return ts.UnmarshalBinary(data)
}

func (ts Time) MarshalJSON() ([]byte, error) {
	t := time.Unix(0, ts.T)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.DateTime)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.DateTime)
	b = append(b, '"')
	return b, nil
}

func (ts *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*ts = Time{T: t.UnixNano()}
	return err
}

// Scan implements the Scanner interface.
func (ts *DeletedAt) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	if err != nil {
		return err
	}
	if !nullTime.Valid {
		*ts = DeletedAt{Seconds: 0, Nanos: -1}
	}
	*ts = DeletedAt{Seconds: int64(nullTime.Time.Second()), Nanos: int32(nullTime.Time.Nanosecond())}
	return err
}

// Value implements the driver Valuer interface.
func (t DeletedAt) Value() (driver.Value, error) {
	if t.Nanos < 0 {
		return nil, nil
	}
	return time.Unix(t.Seconds, int64(t.Nanos)), nil
}

func (t *DeletedAt) InValid() {
	t.Nanos = -1
}

var (
	FlagDeleted = 1
	FlagActived = 0
)

func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
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

func (DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
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

func (DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
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
		return curTime.UnixNano() / 1e6
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

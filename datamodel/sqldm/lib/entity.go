package lib

import (
	"fmt"
	"github.com/csby/database/sqldb"
	"os"
	"path/filepath"
	"strings"
)

type entity struct {
	pkg *Package
}

func (s *entity) openFile(name string, overwrite bool) (*os.File, error) {
	path := filepath.Join(s.pkg.Folder, fmt.Sprintf("%s.go", name))

	info := ", new"
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		if !overwrite {
			return nil, os.ErrExist
		}
		info = ", overwrite"
	}

	os.MkdirAll(s.pkg.Folder, 0777)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	apsPath, _ := filepath.Abs(path)
	fmt.Print("file: ", apsPath)
	fmt.Println(info)

	return file, nil
}

func (s *entity) toFirstUpper(v string) string {
	count := len(v)
	if count == 1 {
		return strings.ToUpper(v)
	} else if count > 1 {
		rs := []rune(v)
		return strings.ToUpper(string(rs[0:1])) + string(rs[1:])
	}

	return v
}

func (s *entity) toFirstLower(v string) string {
	count := len(v)
	if count == 1 {
		return strings.ToUpper(v)
	} else if count > 1 {
		rs := []rune(v)
		return strings.ToLower(string(rs[0:1])) + string(rs[1:])
	}

	return v
}

func (s *entity) toFileName(v string) string {
	vt := strings.ReplaceAll(v, ".", "_")
	vs := strings.Split(vt, "_")
	count := len(vs)
	if count < 1 {
		return ""
	}

	sb := strings.Builder{}
	name := strings.ToLower(vs[0])
	c := 0
	if len(name) > 0 {
		sb.WriteString(name)
		c++
	}
	for i := 1; i < count; i++ {
		name = strings.ToLower(vs[i])
		if len(name) > 0 {
			if c > 0 {
				sb.WriteString(".")
			}
			sb.WriteString(name)
			c++
		}
	}

	return sb.String()
}

func (s *entity) toEntityName(v string) string {
	vt := strings.ReplaceAll(v, ".", "_")
	vs := strings.Split(vt, "_")
	count := len(vs)
	if count < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString(s.toFirstUpper(strings.ToLower(vs[0])))
	for i := 1; i < count; i++ {
		sb.WriteString(s.toFirstUpper(strings.ToLower(vs[i])))
	}

	return sb.String()
}

func (s *entity) toFieldName(columnName string) string {
	if strings.ToUpper(columnName) == "ID" {
		return "ID"
	}

	return s.toEntityName(columnName)
}

func (s *entity) toJsonName(columnName string) string {
	if strings.ToUpper(columnName) == "ID" {
		return "id"
	}

	v := s.toEntityName(columnName)

	return s.toFirstLower(v)
}

func (s *entity) toRuntimeType(dataType string, nullable bool, scale *int) string {
	v := strings.ToLower(dataType)
	scaleValue := 0
	if scale != nil {
		scaleValue = *scale
	}

	if v == "varchar" || v == "varchar2" || v == "nvarchar" || v == "nvarchar2" || v == "text" || v == "ntext" || v == "json" {
		if nullable {
			return "*string"
		} else {
			return "string"
		}
	} else if v == "longtext" || v == "blob" || v == "time" || v == "binary" || v == "varbinary" || v == "longblob" {
		return "[]byte"
	} else if v == "int" || v == "int64" || v == "bigint" || v == "tinyint" || v == "bit" {
		if nullable {
			return "*uint64"
		} else {
			return "uint64"
		}
	} else if v == "uint" || v == "uint64" {
		if nullable {
			return "*uint64"
		} else {
			return "uint64"
		}
	} else if v == "float" || v == "decimal" || v == "numeric" {
		if nullable {
			return "*float64"
		} else {
			return "float64"
		}
	} else if v == "datetime" || v == "date" || v == "datetime2" || v == "timestamp" {
		return "*time.Time"
	} else if v == "uniqueidentifier" {
		if nullable {
			return "*types.UniqueIdentifier"
		} else {
			return "types.UniqueIdentifier"
		}
	} else if v == "number" {
		if scaleValue > 0 {
			if nullable {
				return "*float64"
			} else {
				return "float64"
			}
		} else {
			if nullable {
				return "*uint64"
			} else {
				return "uint64"
			}
		}
	}

	return v
}

func (s *entity) toNoPointerType(v string) string {
	rs := []rune(v)
	if len(rs) < 2 {
		return v
	}

	return string(rs[1:])
}

func (s *entity) importPackages(columns []*sqldb.SqlColumn, model bool) []string {
	packages := make([]string, 0)
	packages = append(packages, "\"bytes\"")
	packages = append(packages, "\"encoding/gob\"")

	temps := make(map[string]string, 0)
	for _, column := range columns {
		columnType := strings.ToLower(column.DataType)
		if columnType == "datetime" || columnType == "date" || columnType == "datetime2" || columnType == "timestamp" {
			if _, ok := temps[columnType]; !ok {
				temps[columnType] = columnType
				packages = append(packages, "\"time\"")
			}
		} else if columnType == "json" {
			if model {
				if _, ok := temps[columnType]; !ok {
					temps[columnType] = columnType
					packages = append(packages, "\"encoding/json\"")
				}
			}
		} else if columnType == "uniqueidentifier" {
			if _, ok := temps[columnType]; !ok {
				temps[columnType] = columnType
				packages = append(packages, "\"github.com/csby/wsf/types\"")
			}
		}
	}

	return packages
}

func (s *entity) getComments(comment string) []string {
	comments := make([]string, 0)
	if len(comment) > 0 {
		lines := strings.Split(comment, "\n")
		for _, line := range lines {
			comments = append(comments, strings.TrimRight(line, "\r"))
		}
	}

	return comments
}

func (s *entity) getNote(comment string) string {
	note := &strings.Builder{}
	cmd := strings.Replace(comment, "\r", "", -1)
	cs := strings.Split(cmd, "\n")
	if len(cs) > 0 {
		note.WriteString(cs[0])
	}
	if len(cs) > 3 {
		note.WriteString(" ")
		note.WriteString(cs[3])
	}
	if len(cs) > 1 {
		note.WriteString(" 例如:")
		note.WriteString(cs[1])
	}

	return note.String()
}

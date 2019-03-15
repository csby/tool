package lib

import (
	"fmt"
	"github.com/csby/database/sqldb"
	"os"
	"sort"
	"strings"
)

type entityModel struct {
	entity
}

func (s *entityModel) toRuntimeType(v string, nullAble bool) string {
	if v == "datetime" || v == "date" {
		return "*types.DateTime"
	} else if v == "varchar" || v == "nvarchar" || v == "text" || v == "longtext" || v == "time" {
		return "string"
	} else if v == "json" {
		return "interface{}"
	}

	return s.entity.toRuntimeType(v, nullAble)
}

func (s *entityModel) importPackages(columns []*sqldb.SqlColumn) []string {
	packages := make([]string, 0)

	temps := make(map[string]string, 0)
	for _, column := range columns {
		if column.DataType == "datetime" || column.DataType == "date" {
			if _, ok := temps[column.DataType]; !ok {
				temps[column.DataType] = column.DataType
				packages = append(packages, "\"github.com/csby/wsf/types\"")
			}
		}
	}

	return packages
}

func (s *entityModel) create(table *sqldb.SqlTable, columns []*sqldb.SqlColumn) error {
	entityFileName := s.toFileName(table.Name)
	entityFile, err := s.openFile(entityFileName, true)
	if err != nil {
		return err
	}
	defer entityFile.Close()
	fmt.Fprintln(entityFile, "package", s.pkg.Name)
	fmt.Fprintln(entityFile, "")

	importPackages := s.importPackages(columns)
	importPackagesCount := len(importPackages)
	if importPackagesCount == 1 {
		fmt.Fprintln(entityFile, "import", importPackages[0])
		fmt.Fprintln(entityFile, "")
	} else if importPackagesCount > 1 {
		sort.Slice(importPackages, func(i, j int) bool {
			return strings.Compare(importPackages[i], importPackages[j]) < 0
		})
		fmt.Fprintln(entityFile, "import", "(")
		importedPackages := make(map[string]string)
		for _, importPkg := range importPackages {
			if _, ok := importedPackages[importPkg]; ok {
				continue
			}
			importedPackages[importPkg] = ""
			fmt.Fprint(entityFile, "	", importPkg)
			fmt.Fprintln(entityFile)
		}
		fmt.Fprintln(entityFile, ")")
		fmt.Fprintln(entityFile, "")
	}

	entityName := s.toEntityName(table.Name)
	fmt.Fprintln(entityFile, "type", entityName, "struct", "{")
	columnNameMaxLength := 0
	columnTypeMaxLength := 0
	for _, column := range columns {
		n := len(s.toFieldName(column.Name))
		if columnNameMaxLength < n {
			columnNameMaxLength = n
		}

		n = len(s.toRuntimeType(column.DataType, column.Nullable))
		if columnTypeMaxLength < n {
			columnTypeMaxLength = n
		}
	}

	for _, column := range columns {
		fieldName := s.toFieldName(column.Name)
		fmt.Fprint(entityFile, "	", fieldName)
		n := len(fieldName)
		for i := n; i <= columnNameMaxLength; i++ {
			fmt.Fprint(entityFile, " ")
		}

		columnType := s.toRuntimeType(column.DataType, column.Nullable)
		fmt.Fprint(entityFile, columnType)
		n = len(columnType)
		for i := n; i <= columnTypeMaxLength; i++ {
			fmt.Fprint(entityFile, " ")
		}

		fmt.Fprint(entityFile, "`json:\"", s.toFirstLower(column.Name), "\"")
		fmt.Fprint(entityFile, " note:\"", s.getNote(column.Comment), "\"")
		fmt.Fprintln(entityFile, "`")
	}
	fmt.Fprintln(entityFile, "}")

	// ext file
	extFile, err := s.openFile(fmt.Sprintf("%s.ext", entityFileName), false)
	if err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}
	defer extFile.Close()
	fmt.Fprintln(extFile, "package", s.pkg.Name)

	return nil
}

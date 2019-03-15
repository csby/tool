package lib

import (
	"fmt"
	"github.com/csby/database/sqldb"
	"os"
	"sort"
	"strings"
)

type entityDatabase struct {
	entity
}

func (s *entityDatabase) create(table *sqldb.SqlTable, columns []*sqldb.SqlColumn, model *entityModel) error {
	entityFileName := s.toFileName(table.Name)
	entityFile, err := s.openFile(entityFileName, true)
	if err != nil {
		return err
	}
	defer entityFile.Close()
	fmt.Fprintln(entityFile, "package", s.pkg.Name)
	fmt.Fprintln(entityFile, "")

	createModel := false
	if model != nil {
		if model.pkg != nil {
			createModel = true
		}
	}
	importPackages := s.importPackages(columns, createModel)
	if createModel {
		importPackages = append(importPackages, fmt.Sprintf("\"%s\"", model.pkg.Path))
		modelImportPackages := model.importPackages(columns)
		if len(modelImportPackages) > 0 {
			importPackages = append(importPackages, modelImportPackages...)
		}
	}

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
	// base for table name
	entityBaseName := fmt.Sprintf("%sBase", entityName)
	s.createBase(entityFile, entityBaseName, table, columns)

	entityComments := s.getComments(table.Description)
	if len(entityComments) > 0 {
		for _, entityComment := range entityComments {
			fmt.Fprint(entityFile, "//")
			if len(entityComment) > 1 {
				fmt.Fprint(entityFile, " ", entityComment)
			}
			fmt.Fprintln(entityFile)
		}
	}

	fmt.Fprintln(entityFile, "type", entityName, "struct", "{")
	fmt.Fprint(entityFile, "	", entityBaseName)
	fmt.Fprintln(entityFile, "")
	for _, column := range columns {
		note := strings.TrimSpace(s.getNote(column.Comment))
		fmt.Fprint(entityFile, "	//")
		if len(note) > 0 {
			fmt.Fprint(entityFile, " ", note)
		}
		fmt.Fprintln(entityFile)

		fmt.Fprint(entityFile, "	", s.toFirstUpper(column.Name))
		fmt.Fprint(entityFile, " ", s.toRuntimeType(column.DataType, column.Nullable))
		fmt.Fprint(entityFile, " `sql:\"", column.Name, "\"")
		if column.AutoIncrement {
			fmt.Fprint(entityFile, " auto:\"true\"")
		}
		if column.PrimaryKey {
			fmt.Fprint(entityFile, " primary:\"true\"")
		}
		fmt.Fprintln(entityFile, "`")
	}
	fmt.Fprintln(entityFile, "}")

	fmt.Fprintln(entityFile, "")
	s.createCloneFunction(entityFile, entityName)

	if createModel {
		fmt.Fprintln(entityFile, "")
		s.createCopyToFunction(entityFile, entityName, table, columns, model)

		fmt.Fprintln(entityFile, "")
		s.createCopyFromFunction(entityFile, entityName, table, columns, model)
	}

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

func (s *entityDatabase) createBase(entityFile *os.File, entityBaseName string, table *sqldb.SqlTable, columns []*sqldb.SqlColumn) {
	fmt.Fprintln(entityFile, "type", entityBaseName, "struct", "{")
	fmt.Fprintln(entityFile, "}")
	fmt.Fprintln(entityFile, "")
	fmt.Fprint(entityFile, "func (s ", entityBaseName, ") TableName()", " string")
	fmt.Fprintln(entityFile, " {")
	fmt.Fprint(entityFile, "	return \"", table.Name, "\"")
	fmt.Fprintln(entityFile, "")
	fmt.Fprintln(entityFile, "}")
	fmt.Fprintln(entityFile, "")

	fmt.Fprint(entityFile, "func (s ", entityBaseName, ") SetFilter(v interface{})")
	fmt.Fprintln(entityFile, " {")
	fmt.Fprintln(entityFile, "")
	fmt.Fprintln(entityFile, "}")
	fmt.Fprintln(entityFile, "")

	fmt.Fprint(entityFile, "func (s *", entityBaseName, ") Clone() (interface{}, error)")
	fmt.Fprintln(entityFile, " {")
	fmt.Fprint(entityFile, "	return &", entityBaseName, "{}, nil")
	fmt.Fprintln(entityFile, "")
	fmt.Fprintln(entityFile, "}")
	fmt.Fprintln(entityFile, "")
}

func (s *entityDatabase) createCloneFunction(entityFile *os.File, entityName string) {
	fmt.Fprint(entityFile, "func (s *", entityName, ") Clone() (interface{}, error)")
	fmt.Fprintln(entityFile, " {")
	fmt.Fprint(entityFile, "	t := &", entityName, "{}")
	fmt.Fprintln(entityFile, "")
	fmt.Fprintln(entityFile, "")

	fmt.Fprintln(entityFile, "	buf := new(bytes.Buffer)")
	fmt.Fprintln(entityFile, "	enc := gob.NewEncoder(buf)")
	fmt.Fprintln(entityFile, "	dec := gob.NewDecoder(buf)")
	fmt.Fprintln(entityFile, "	err := enc.Encode(s)")
	fmt.Fprintln(entityFile, "	if err != nil {")
	fmt.Fprintln(entityFile, "		return nil, err")
	fmt.Fprintln(entityFile, "	}")
	fmt.Fprintln(entityFile, "	err = dec.Decode(t)")
	fmt.Fprintln(entityFile, "	if err != nil {")
	fmt.Fprintln(entityFile, "		return nil, err")
	fmt.Fprintln(entityFile, "	}")

	fmt.Fprintln(entityFile, "")
	fmt.Fprintln(entityFile, "	return t, nil")
	fmt.Fprintln(entityFile, "}")
}

func (s *entityDatabase) createCopyToFunction(entityFile *os.File, entityName string, table *sqldb.SqlTable, columns []*sqldb.SqlColumn, model *entityModel) {
	fmt.Fprint(entityFile, "func (s *", entityName, ") CopyTo(target *", model.pkg.Name, ".", entityName)
	fmt.Fprintln(entityFile, ") {")
	fmt.Fprintln(entityFile, "	if target == nil {")
	fmt.Fprintln(entityFile, "		return")
	fmt.Fprintln(entityFile, "	}")

	for _, column := range columns {
		fieldName := s.toFirstUpper(column.Name)
		sourceType := s.toRuntimeType(column.DataType, column.Nullable)
		targetType := model.toRuntimeType(column.DataType, column.Nullable)
		if sourceType == targetType {
			fmt.Fprintf(entityFile, "	target.%s = s.%s", fieldName, fieldName)
		} else {
			tmpValName := s.toJsonName(column.Name)
			if column.Type == "json" {
				if column.Nullable {
					fmt.Fprint(entityFile, "	", tmpValName, " := \"\"")
					fmt.Fprintln(entityFile)
					fmt.Fprintf(entityFile, "	if s.%s != nil {", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "		", tmpValName, " = *s.", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprintln(entityFile, "	}")
				} else {
					fmt.Fprint(entityFile, "	", tmpValName, " := s.", fieldName)
					fmt.Fprintln(entityFile)
				}

				fmt.Fprintf(entityFile, "	if len(%s) > 0 {", tmpValName)
				fmt.Fprintln(entityFile)
				fmt.Fprintf(entityFile, "		json.Unmarshal([]byte(%s), &target.%s)", tmpValName, fieldName)
				fmt.Fprintln(entityFile)
				fmt.Fprintln(entityFile, "	}")

			} else if strings.HasPrefix(targetType, "*") {
				if strings.HasPrefix(sourceType, "*") {
					fmt.Fprintf(entityFile, "	if s.%s == nil {", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprintf(entityFile, "		target.%s = nil", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprintln(entityFile, "	}", "else", "{")
					fmt.Fprint(entityFile, "		", tmpValName, " := ", s.toNoPointerType(targetType), "(*s.", fieldName, ")")
					fmt.Fprintln(entityFile, "")
					fmt.Fprint(entityFile, "		target.", fieldName, " = &", tmpValName)
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "	}")
				} else {
					fmt.Fprint(entityFile, "	", tmpValName, " := ", s.toNoPointerType(targetType), "(s.", fieldName, ")")
					fmt.Fprintln(entityFile, "")
					fmt.Fprint(entityFile, "	target.", fieldName, " = &", tmpValName)
				}
			} else {
				if strings.HasPrefix(sourceType, "*") {
					fmt.Fprintf(entityFile, "	if s.%s != nil {", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "		target.", fieldName, " = ", targetType, "(*s.", fieldName, ")")
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "	}")
				} else {
					fmt.Fprint(entityFile, "	target.", fieldName, " = ", targetType, "(s.", fieldName, ")")
				}
			}
		}

		fmt.Fprintln(entityFile, "")
	}

	fmt.Fprintln(entityFile, "}")
}

func (s *entityDatabase) createCopyFromFunction(entityFile *os.File, entityName string, table *sqldb.SqlTable, columns []*sqldb.SqlColumn, model *entityModel) {
	fmt.Fprint(entityFile, "func (s *", entityName, ") CopyFrom(source *", model.pkg.Name, ".", entityName)
	fmt.Fprintln(entityFile, ") {")
	fmt.Fprintln(entityFile, "	if source == nil {")
	fmt.Fprintln(entityFile, "		return")
	fmt.Fprintln(entityFile, "	}")

	for _, column := range columns {
		fieldName := s.toFieldName(column.Name)
		sourceType := model.toRuntimeType(column.DataType, column.Nullable)
		targetType := s.toRuntimeType(column.DataType, column.Nullable)
		if sourceType == targetType {
			fmt.Fprintf(entityFile, "	s.%s = source.%s", fieldName, fieldName)
		} else {
			tmpValName := s.toJsonName(fieldName)
			if strings.HasPrefix(targetType, "*") {
				if strings.HasPrefix(sourceType, "*") {
					fmt.Fprintf(entityFile, "	if source.%s == nil {", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprintf(entityFile, "		s.%s = nil", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprintln(entityFile, "	}", "else", "{")
					fmt.Fprint(entityFile, "		", tmpValName, " := ", s.toNoPointerType(targetType), "(*source.", fieldName, ")")
					fmt.Fprintln(entityFile, "")
					fmt.Fprint(entityFile, "		s.", fieldName, " = &", tmpValName)
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "	}")
				} else {
					fmt.Fprint(entityFile, "	", tmpValName, " := ", s.toNoPointerType(targetType), "(source.", fieldName, ")")
					fmt.Fprintln(entityFile, "")
					fmt.Fprint(entityFile, "	s.", fieldName, " = &", tmpValName)
				}
			} else {
				if strings.HasPrefix(sourceType, "*") {
					fmt.Fprintf(entityFile, "	if source.%s != nil {", fieldName)
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "		s.", fieldName, " = ", targetType, "(*source.", fieldName, ")")
					fmt.Fprintln(entityFile)
					fmt.Fprint(entityFile, "	}")
				} else {
					fmt.Fprint(entityFile, "	s.", fieldName, " = ", targetType, "(source.", fieldName, ")")
				}
			}
		}
		fmt.Fprintln(entityFile, "")
	}

	fmt.Fprintln(entityFile, "}")
}

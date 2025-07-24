package helper

import (
	"fmt"
	"strings"
)

func BuildInClause(values string) string {
	sections := strings.Split(values, ",")
	for i, s := range sections {
		sections[i] = fmt.Sprintf("'%s'", s)
	}
	return strings.Join(sections, ",")
}

func BuildFilterWithQUery(field, value, filter string) string {
	if value != "" {
		return fmt.Sprintf("%s AND (LOWER(%s) LIKE LOWER('%%%s%%'))", filter, field, value)
	}
	return filter
}

func BuildManyColumnFilter(fieldsAndValues map[string]string, filter string) string {
	var conditions []string

	for field, value := range fieldsAndValues {
		if value != "" {
			conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE LOWER('%%%s%%')", field, value))
		}
	}

	if len(conditions) > 0 {
		filter += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	return filter
}

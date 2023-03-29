package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type ConditionSeparator string
type ValuesType string

const (
	SeparatorOR  ConditionSeparator = "OR"
	SeparatorAND ConditionSeparator = "AND"
)

const (
	INTEGER ValuesType = "INTEGER"
	VARCHAR ValuesType = "VARCHAR"
	UUID    ValuesType = "UUID"
	BOOL    ValuesType = "BOOLEAN"
)

func ParseEqual(paramName string, paramValue any, valueType ValuesType) string {
	strValue := pq.QuoteLiteral(strings.Trim(fmt.Sprintf("%v", paramValue), " "))
	if strValue == "" {
		return ""
	}

	preparedParamName := pq.QuoteIdentifier(paramName)

	if strings.Trim(strings.ToLower(strValue), "'") == "null" {
		return fmt.Sprintf(`%s IS NULL`, preparedParamName)
	}
	return fmt.Sprintf(`%s = %s::%s`, preparedParamName, strValue, valueType)
}

func ParseIN(paramName string, paramValues any, valuesType ValuesType) string {
	values := reflect.ValueOf(paramValues)
	if values.Kind() != reflect.Slice {
		return ""
	}

	if values.Len() == 0 {
		return ""
	}
	var strValues []string
	for i := 0; i < values.Len(); i++ {
		strValues = append(strValues, pq.QuoteLiteral(strings.Trim(fmt.Sprintf("%v", values.Index(i).Interface()), " ")))
	}

	var filters []string

	var notNullValues []string
	for _, value := range strValues {
		if strings.Trim(strings.ToLower(value), "'") == "null" {
			filters = append(filters, fmt.Sprintf(`%s IS NULL`, pq.QuoteIdentifier(paramName)))
		} else {
			notNullValues = append(notNullValues, fmt.Sprintf("%s::%s", value, valuesType))
		}
	}

	inFilter := strings.Join(notNullValues, ", ")
	if len(notNullValues) != 0 {
		filters = append(filters, fmt.Sprintf(`%s IN (%s)`, paramName, inFilter))
	}

	return fmt.Sprintf("(%s)", strings.Join(filters, " OR "))
}

func ParseLike(paramName string, paramValue string, useILIKE bool) string {
	like := "LIKE"
	if useILIKE {
		like = "ILIKE"
	}

	preparedValue := pq.QuoteLiteral("%" + strings.Trim(paramValue, " ") + "%")

	return fmt.Sprintf(`%s %s %s`, pq.QuoteIdentifier(paramName), like, preparedValue)
}

func ParseInLike(paramName string, paramValues []string, useILIKE bool) string {
	if len(paramValues) == 0 {
		return ""
	}

	var filters []string
	for _, v := range paramValues {
		filters = append(filters, ParseLike(paramName, v, useILIKE))
	}

	return fmt.Sprintf("(%s)", strings.Join(filters, " OR "))
}

func ParseValuesInParam(paramName string, paramValues any, conditionSeparator ConditionSeparator, valuesType ValuesType) string {
	values := reflect.ValueOf(paramValues)
	if values.Kind() != reflect.Slice {
		return ""
	}

	if values.Len() == 0 {
		return ""
	}

	var strValues []string
	for i := 0; i < values.Len(); i++ {
		strValues = append(strValues, pq.QuoteLiteral(strings.Trim(fmt.Sprintf("%v", values.Index(i).Interface()), " ")))
	}

	var filters []string
	var notNullValues []string

	for _, v := range strValues {
		if strings.Trim(strings.ToLower(v), "'") == "null" {
			filters = append(filters, fmt.Sprintf(`%s IS NULL`, pq.QuoteIdentifier(paramName)))
		} else {
			notNullValues = append(notNullValues, fmt.Sprintf("%s::%s", v, valuesType))
		}
	}

	separator := ""
	switch conditionSeparator {
	case SeparatorOR:
		separator = "&&"
	case SeparatorAND:
		separator = "@>"
	default:
		return ""
	}

	inFilter := fmt.Sprintf(`%s %s ARRAY[%s]`,
		pq.QuoteIdentifier(paramName), separator, strings.Join(notNullValues, ", "))
	filters = append(filters, inFilter)

	return fmt.Sprintf("(%s)", strings.Join(filters, " OR "))
}

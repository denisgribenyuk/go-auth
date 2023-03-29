package dbqueries

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies/models"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func GetCompanyUsers(companyID int, queryData models.GetCompanyUsersRequestQuery) ([]models.DBCompanyUsers, int, error) {
	query := `
        SELECT *
        FROM (
                 WITH roles AS
                          (
                              SELECT user_id,
                                     array_agg(role_id ORDER BY role_id) role_ids
                              FROM roles_users
                              GROUP BY user_id
                          ),
                      structures AS
                          (
                              SELECT user_id,
                                     array_agg(structure_id ORDER BY structure_id) structure_ids
                              FROM users_company_structures
                              GROUP BY user_id
                          ),
                      fields AS
                          (
                              SELECT user_id,
                                     array_agg(field_guid) field_guids
                              FROM users_fields
                              GROUP BY user_id
                          )
                 SELECT u.id,
                        u.ext_id,
                        trim(concat(u.last_name, ' ', u.first_name, ' ', u.middle_name)) user_name,
                        u.first_name,
                        u.last_name,
                        u.middle_name,
                        u.email,
                        u.phone,
                        u.company_id,
                        i.image_path                                                     image_name,
                        u.position_id,
                        cp.name                                                          position_name,
                        u.manager_id,
                        mu.last_name                                                     manager_last_name,
                        COALESCE(ucs.structure_ids, ARRAY []::integer[])                 structure_ids,
                        COALESCE(ru.role_ids, ARRAY []::integer[])                       role_ids,
                        COALESCE(uf.field_guids, ARRAY []::uuid[])                       field_guids,
                        CASE u.active_flag WHEN 1 THEN TRUE ELSE FALSE END AS            active_flag
                 FROM users u
                          LEFT JOIN users mu ON mu.id = u.manager_id
                          LEFT JOIN company_positions cp ON cp.id = u.position_id
                          LEFT JOIN images i ON u.image_id = i.id
                          LEFT JOIN structures ucs ON u.id = ucs.user_id
                          LEFT JOIN roles ru ON u.id = ru.user_id
                          LEFT JOIN fields uf ON u.id = uf.user_id
             ) cte
	`
	filters := []string{
		postgres.ParseEqual("company_id", companyID, postgres.INTEGER),
	}
	if len(queryData.UserIDs) != 0 {
		condition := postgres.ParseIN("id", queryData.UserIDs, postgres.INTEGER)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.ExtIDs) != 0 {
		condition := postgres.ParseInLike("ext_id", queryData.ExtIDs, false)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if queryData.UserName != nil {
		condition := postgres.ParseLike("user_name", *queryData.UserName, true)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.PositionIDs) != 0 {
		condition := postgres.ParseIN("position_id", queryData.PositionIDs, postgres.INTEGER)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if queryData.Email != nil {
		condition := postgres.ParseLike("email", *queryData.Email, true)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if queryData.Phone != nil {
		condition := postgres.ParseLike("phone", *queryData.Phone, true)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.RoleIDs) != 0 {
		condition := postgres.ParseValuesInParam("role_ids", queryData.RoleIDs, postgres.SeparatorOR, postgres.INTEGER)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.ManagerIDs) != 0 {
		condition := postgres.ParseIN("manager_id", queryData.ManagerIDs, postgres.INTEGER)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.StructureIDs) != 0 {
		condition := postgres.ParseValuesInParam("structure_ids", queryData.StructureIDs, postgres.SeparatorOR, postgres.INTEGER)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if len(queryData.FieldGUIDs) != 0 {
		condition := postgres.ParseValuesInParam("field_guids", queryData.FieldGUIDs, postgres.SeparatorOR, postgres.UUID)
		if condition != "" {
			filters = append(filters, condition)
		}
	}
	if queryData.ActiveFlag != nil {
		condition := postgres.ParseEqual("active_flag", *queryData.ActiveFlag, postgres.BOOL)
		if condition != "" {
			filters = append(filters, condition)
		}
	}

	filtersString := fmt.Sprintf("WHERE (%s)", strings.Join(filters, "\nAND "))

	query = fmt.Sprintf("%s\n%s", query, filtersString)

	totalCount, err := postgres.GetCount(nil, query)
	if err != nil {
		return []models.DBCompanyUsers{}, 0, err
	}

	limit := "NULL"
	if queryData.Limit != nil {
		limit = pq.QuoteLiteral(fmt.Sprintf("%d", *queryData.Limit))
	}
	offset := pq.QuoteLiteral(fmt.Sprintf("%d", queryData.Offset))
	limitation := fmt.Sprintf("LIMIT %s OFFSET %s", limit, offset)

	sortColumn := "last_name"
	if queryData.SortColumn != nil {
		sortColumn = *queryData.SortColumn
	}
	sortDirection := "ASC"
	if queryData.SortDirection != nil {
		sortDirection = *queryData.SortDirection
	}
	ordering := fmt.Sprintf("ORDER BY %s %s", pq.QuoteIdentifier(sortColumn), sortDirection)

	query = fmt.Sprintf("%s\n%s\n%s", query, ordering, limitation)

	var res []models.DBCompanyUsers
	err = postgres.DB.Select(&res, query)
	if err != nil {
		return []models.DBCompanyUsers{}, 0, err
	}

	return res, totalCount, nil
}

package delobdriver

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (c *DelobContext) Query(queryBuilder QueryBuilder) ([]Player, error) {
	q, err := queryBuilder.Query()
	if err != nil {
		return []Player{}, err
	}

	return c.getPlayersQuery(q)
}

func (c *DelobContext) getPlayersQuery(expression string) ([]Player, error) {
	jsonResponse, errDelob := c.sendMessage(expression)
	if errDelob != nil {
		return nil, errDelob
	}

	result := []Player{}
	err := json.Unmarshal([]byte(jsonResponse), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Player struct {
	Key     string  `json:"Key"`
	Elo     int16   `json:"Elo,omitempty"`
	Events  []Event `json:"Events,omitempty"`
	records []int16
}

type Event struct {
	Change      int16 `json:"Change,omitempty"`
	DateTime    time.Time
	TeamOne     []string `json:"TeamOne,omitempty"`
	TeamTwo     []string `json:"TeamTwo,omitempty"`
	MatchResult string   `json:"MatchResult,omitempty"`
}

type QueryComponent string

const (
	Key     QueryComponent = "Key"
	Elo     QueryComponent = "Elo"
	Events  QueryComponent = "Events"
	Matches QueryComponent = "Matches"
)

type Catalog string

const (
	Players Catalog = "Players"
)

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)

type QueryBuilder struct {
	catalog         Catalog
	queryComponents []QueryComponent
	orderKey        QueryComponent
	orderDirection  OrderDirection
	query           string
	validQuery      bool
	operationsCount int8
}

func (c *DelobContext) Select(queryComponents ...QueryComponent) QueryBuilder {
	selectCtx := QueryBuilder{
		query: "SELECT ",
	}

	if len(queryComponents) == 0 {
		selectCtx.query += "*"
	} else {
		selectCtx.query += joinComponents(queryComponents, ", ")
	}

	selectCtx.query += " "
	selectCtx.operationsCount++
	return selectCtx
}

func (s QueryBuilder) From(catalog Catalog) QueryBuilder {
	if s.operationsCount != 1 {
		return s
	}

	s.query += "FROM " + string(catalog)
	s.validQuery = true
	s.operationsCount++
	return s
}

func (s QueryBuilder) OrderBy(orderKey QueryComponent, orderDirection OrderDirection) QueryBuilder {
	if s.operationsCount != 2 {
		return s
	}

	if orderKey == Events || orderKey == Matches {
		return s
	}

	s.query += fmt.Sprintf(" ORDER BY %s %s", orderKey, orderDirection)
	s.validQuery = true
	s.operationsCount++
	return s
}

func (s *QueryBuilder) Query() (string, error) {
	if s.operationsCount != 3 || !s.validQuery {
		return "", fmt.Errorf("there wa a problem with queryBuilder - %s", s.query)
	}
	return s.query + ";", nil
}

func joinComponents(comps []QueryComponent, sep string) string {
	parts := make([]string, len(comps))
	for i, c := range comps {
		parts[i] = string(c)
	}
	return strings.Join(parts, sep)
}

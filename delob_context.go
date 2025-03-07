package delobdriver

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DelobContext struct {
	tcpHandler *tcpHandler
}

type Player struct {
	Key string
	Elo int
}

type OrderKey string

const (
	Key OrderKey = "Key"
	Elo OrderKey = "Elo"
)

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)

func NewContext(connectionString string) (DelobContext, error) {
	tcpHandler, err := newTcpHandler(connectionString)
	if err != nil {
		return DelobContext{}, err
	}

	return DelobContext{
		tcpHandler: tcpHandler,
	}, nil
}

func (c *DelobContext) AddPlayer(playerKey string) error {
	expression := fmt.Sprintf("ADD PLAYER '%s';", playerKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) AddPlayers(playerKeys []string) error {
	expression := fmt.Sprintf("ADD PLAYERS %s;", creteCollectionFromArray(playerKeys))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDecisiveTeamMatch(teamWinKeys, teamLoseKeys []string) error {
	expression := fmt.Sprintf("SET WIN FOR %s AND LOSE FOR %s;", creteCollectionFromArray(teamWinKeys), creteCollectionFromArray(teamLoseKeys))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDecisiveMatch(playerOneKey, playerTwoKey string) error {
	expression := fmt.Sprintf("SET WIN FOR '%s' AND LOSE FOR '%s';", playerOneKey, playerTwoKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDrawTeamMatch(teamOne, teamTwo []string) error {
	expression := fmt.Sprintf("SET DRAW BETWEEN %s AND %s;", creteCollectionFromArray(teamOne), creteCollectionFromArray(teamTwo))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDrawMatch(playerOneKey, playerTwoKey string) error {
	expression := fmt.Sprintf("SET DRAW BETWEEN '%s' AND '%s';", playerOneKey, playerTwoKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) GetPlayers() ([]Player, error) {
	expression := "SELECT Players;"

	return c.getPlayersQuery(expression)
}

func (c *DelobContext) GetPlayersOrderBy(orderKey OrderKey, orderDirection OrderDirection) ([]Player, error) {
	expression := fmt.Sprintf("SELECT Players ORDER BY %s %s;", orderKey, orderDirection)

	return c.getPlayersQuery(expression)
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

func creteCollectionFromArray(arr []string) string {
	result := []string{}
	for i := range arr {
		result = append(result, fmt.Sprintf("'%s'", arr[i]))
	}

	return fmt.Sprintf("(%s)", strings.Join(result, ","))
}

func (c *DelobContext) sendMessage(expression string) (string, error) {
	response, err := c.tcpHandler.sendRequest(expression)
	if err != nil {
		return "", err
	}

	return response, nil
}

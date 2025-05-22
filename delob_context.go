package delobdriver

import (
	"fmt"
	"strings"
)

type DelobContext struct {
	tcpHandler *tcpHandler
}

func NewContext(connectionString string) (DelobContext, error) {
	tcpHandler, err := newTcpHandler(connectionString)
	if err != nil {
		return DelobContext{}, err
	}

	return DelobContext{
		tcpHandler: tcpHandler,
	}, nil
}

func (c *DelobContext) Close() {
	c.tcpHandler.conn.Close()
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

func creteCollectionFromArray(arr []string) string {
	result := []string{}
	for i := range arr {
		result = append(result, fmt.Sprintf("'%s'", arr[i]))
	}

	return fmt.Sprintf("(%s)", strings.Join(result, ","))
}

func (c *DelobContext) sendMessage(expression string) (string, error) {
	response, err := c.tcpHandler.sendRequest(expression, 0)
	if err != nil {
		return "", err
	}

	return response, nil
}

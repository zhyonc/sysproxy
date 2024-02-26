package controller

import (
	"sysproxy/service"
)

type ruleController struct {
	menuHandler MenuController
}

func NewRuleController(menuHandler MenuController) RuleController {
	ruleController := &ruleController{
		menuHandler: menuHandler,
	}
	return ruleController
}

func (c *ruleController) GetUserRule() []string {
	return service.PACService.GetUserRule()
}

func (c *ruleController) SaveUserRule(rules string) error {
	err := service.PACService.SaveUserRule(rules)
	if err != nil {
		return err
	}
	return nil
}

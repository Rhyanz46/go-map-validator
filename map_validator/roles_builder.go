package map_validator

func BuildRoles() RulesWrapper {
	return &rulesWrapper{}
}

func (rw *rulesWrapper) SetRule(field string, rule Rules) RulesWrapper {
	if rw.Rules == nil {
		rw.Rules = make(map[string]Rules)
	}
	rw.Rules[field] = rule
	return rw
}

func (rw *rulesWrapper) Done() RulesWrapper {
	return rw
}

func (rw *rulesWrapper) SetManipulator(field string, fun func(data interface{}) (result interface{}, err error)) RulesWrapper {
	rw.manipulator = append(rw.manipulator, manipulator{Field: field, Func: &fun})
	return rw
}

func (rw *rulesWrapper) SetFieldsManipulator(fields []string, fun func(data interface{}) (result interface{}, err error)) RulesWrapper {
	for _, field := range fields {
		rw.manipulator = append(rw.manipulator, manipulator{Field: field, Func: &fun})
	}
	return rw
}

func (rw *rulesWrapper) SetSetting(setting Setting) RulesWrapper {
	rw.Setting = setting
	return rw
}

func (rw *rulesWrapper) getRules() map[string]Rules {
	return rw.Rules
}

func (rw *rulesWrapper) getSetting() Setting {
	return rw.Setting
}

func (rw *rulesWrapper) getManipulator() []manipulator {
	return rw.manipulator
}

func (rw *rulesWrapper) isList() bool {
	return rw.isListRules
}

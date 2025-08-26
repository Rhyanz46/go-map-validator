package map_validator

func BuildListRoles() ListRulesWrapper {
	return &rulesWrapper{
		isListRules: true,
	}
}

func (rw *rulesWrapper) SetListRule(rule ListRules) ListRulesWrapper {
	rw.ListRules = rule
	return rw
}

// func (rw *rulesWrapper) SetManipulator(field string, fun func(data interface{}) (result interface{}, err error)) RulesWrapper {
// 	rw.manipulator = append(rw.manipulator, manipulator{Field: field, Func: &fun})
// 	return rw
// }

// func (rw *rulesWrapper) SetSetting(setting Setting) RulesWrapper {
// 	rw.Setting = setting
// 	return rw
// }

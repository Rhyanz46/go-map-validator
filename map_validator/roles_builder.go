package map_validator

func BuildRoles() *RulesWrapper {
	return &RulesWrapper{}
}

func (rw *RulesWrapper) SetRule(field string, rule Rules) *RulesWrapper {
	if rw.Rules == nil {
		rw.Rules = make(map[string]Rules)
	}
	rw.Rules[field] = rule
	return rw
}

func (rw *RulesWrapper) SetManipulator(field string, fun func(data interface{}) (result interface{}, err error)) *RulesWrapper {
	rw.manipulator = append(rw.manipulator, manipulator{Field: field, Func: &fun})
	return rw
}

func (rw *RulesWrapper) SetFieldsManipulator(fields []string, fun func(data interface{}) (result interface{}, err error)) *RulesWrapper {
	for _, field := range fields {
		rw.manipulator = append(rw.manipulator, manipulator{Field: field, Func: &fun})
	}
	return rw
}

func (rw *RulesWrapper) SetSetting(setting Setting) *RulesWrapper {
	rw.Setting = setting
	return rw
}

func (rw *RulesWrapper) Done() RulesWrapper {
	return *rw
}

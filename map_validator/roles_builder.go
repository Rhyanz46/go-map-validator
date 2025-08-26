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

func (rw *rulesWrapper) getUniqueValues() *map[string]map[string]interface{} {
	return rw.uniqueValues
}

func (rw *rulesWrapper) setUniqueValues(values *map[string]map[string]interface{}) RulesWrapper {
	rw.uniqueValues = values
	return rw
}

func (rw *rulesWrapper) getFilledField() *[]string {
	return rw.filledField
}

func (rw *rulesWrapper) setFilledField(fields *[]string) RulesWrapper {
	rw.filledField = fields
	return rw
}

func (rw *rulesWrapper) appendFilledField(fields string) {
	*rw.filledField = append(*rw.filledField, fields)
}

func (rw *rulesWrapper) getNullFields() *[]string {
	return rw.nullFields
}

func (rw *rulesWrapper) setNullFields(fields *[]string) RulesWrapper {
	rw.nullFields = fields
	return rw
}

func (rw *rulesWrapper) appendNullFields(fields string) {
	*rw.nullFields = append(*rw.nullFields, fields)
}

func (rw *rulesWrapper) getRequiredWithout() *map[string][]string {
	return rw.requiredWithout
}

func (rw *rulesWrapper) setRequiredWithout(req *map[string][]string) RulesWrapper {
	rw.requiredWithout = req
	return rw
}

func (rw *rulesWrapper) getRequiredIf() *map[string][]string {
	return rw.requiredIf
}

func (rw *rulesWrapper) setRequiredIf(req *map[string][]string) RulesWrapper {
	rw.requiredIf = req
	return rw
}

func (rw *rulesWrapper) getManipulator() []manipulator {
	return rw.manipulator
}

func (rw *rulesWrapper) Done() rulesWrapper {
	return *rw
}

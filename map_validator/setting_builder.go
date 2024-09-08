package map_validator

func BuildSetting() *Setting {
	return &Setting{}
}

func (s *Setting) MakeStrict() *Setting {
	s.Strict = true
	return s
}

func (s *Setting) Done() Setting {
	return *s
}

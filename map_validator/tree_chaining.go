package map_validator

import (
	"fmt"
)

type chainState struct {
	key         string
	manipulator *func(interface{}) (interface{}, error)
	value       interface{}
	CustomMsg   *CustomMsg
	uniques     []string
	errs        []error
	parent      *chainState
	children    []*chainState
}

func (cs *chainState) AddError(err error) ChainerType {
	cs.errs = append(cs.errs, err)
	return cs
}

func (cs *chainState) SetCustomMsg(customMsg *CustomMsg) ChainerType {
	cs.CustomMsg = customMsg
	return cs
}

func (cs *chainState) GetErrors() []error {
	return cs.recursiveGetErrors()
}

func (cs *chainState) recursiveGetErrors() []error {
	var errors []error
	errors = append(errors, cs.errs...)
	for _, child := range cs.children {
		errors = append(errors, child.recursiveGetErrors()...)
	}
	return errors
}

func (cs *chainState) SetManipulator(manipulator *func(interface{}) (interface{}, error)) ChainerType {
	cs.manipulator = manipulator
	return cs
}

func (cs *chainState) Next(index int) ChainerType {
	if index < 0 || index >= len(cs.children) {
		panic("invalid child index")
	}
	return cs.children[index]
}

func (cs *chainState) Back() ChainerType {
	if cs.parent == nil {
		panic("can't back")
	}
	return cs.parent
}

func (cs *chainState) Forward(index int) ChainerType {
	if index < 0 || index >= len(cs.children) {
		panic("invalid child index")
	}
	return cs.children[index]
}

func (cs *chainState) GetKey() string {
	return cs.key
}

func (cs *chainState) GetParentKey() string {
	if cs.parent == nil {
		return ""
	}
	return cs.parent.key
}

func (cs *chainState) SetKey(name string) ChainerType {
	cs.key = name
	return cs
}

func (cs *chainState) SetKeyValue(key string, value interface{}) ChainerType {
	cs.value = value
	cs.key = key
	return cs
}

func (cs *chainState) GetParentKeys() []string {
	var keys []string
	current := cs
	for current.parent != nil {
		current = current.parent
		keys = append(keys, current.GetKey())
	}
	return keys
}

func (cs *chainState) AddChild() ChainerType {
	newChild := &chainState{
		parent: cs,
	}
	cs.children = append(cs.children, newChild)
	return newChild
}

func (cs *chainState) PrintHierarchyWithSeparator(separator string, currentPath string) {
	if currentPath == "" {
		currentPath = cs.GetKey()
	} else {
		currentPath = currentPath + separator + cs.GetKey()
	}
	if cs.value != nil {
		fmt.Printf("%s : %v\n", currentPath, cs.value)
	} else {
		fmt.Println(currentPath)
	}
	for _, child := range cs.children {
		child.PrintHierarchyWithSeparator(separator, currentPath)
	}
}

func (cs *chainState) LoadFromMap(data map[string]interface{}) {
	for key, value := range data {
		child := cs.AddChild().(*chainState)
		child.SetKey(key)
		child.SetValue(value)
		if subMap, ok := value.(map[string]interface{}); ok {
			child.LoadFromMap(subMap)
		}
	}
}

func (cs *chainState) SetValue(value interface{}) ChainerType {
	cs.value = value
	return cs
}

func (cs *chainState) GetValue() interface{} {
	return cs.value
}

func (cs *chainState) recursiveToMap() map[string]interface{} {
	result := make(map[string]interface{})
	if len(cs.children) == 0 {
		// If no children, set value directly, or set to nil if value is nil
		if cs.value != nil {
			result[cs.key] = cs.value
		} else {
			result[cs.key] = nil
		}
	} else {
		// If there are children, create a nested map
		childMap := make(map[string]interface{})
		for _, child := range cs.children {
			for k, v := range child.recursiveToMap() {
				childMap[k] = v
			}
		}
		result[cs.key] = childMap
	}
	return result
}

func (cs *chainState) ToMap() map[string]interface{} {
	var ok bool
	result := make(map[string]interface{})
	res := cs.recursiveToMap()
	result, ok = res[chainKey].(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return result
}

func (cs *chainState) GetUniques() []string {
	return cs.uniques
}

func (cs *chainState) GetBrothers() []ChainerType {
	var brothers []ChainerType
	if cs.GetParent() == nil {
		return nil
	}

	for _, child := range cs.GetParent().GetChildren() {
		if child == nil || cs.GetKey() == child.GetKey() {
			continue
		}
		brothers = append(brothers, child)
	}
	return brothers
}

func (cs *chainState) GetParent() ChainerType {
	return cs.parent
}

func (cs *chainState) GetChildren() []ChainerType {
	var children []ChainerType
	if cs == nil || cs.children == nil {
		return []ChainerType{}
	}
	for _, child := range cs.children {
		children = append(children, child)
	}
	return children
}

func (cs *chainState) SetUniques(uniques []string) ChainerType {
	var removedDuplicated []string
	for _, unique := range uniques {
		var found bool
		for _, removed := range removedDuplicated {
			if removed == unique {
				found = true
				break
			}
		}
		if !found {
			removedDuplicated = append(removedDuplicated, unique)
		}
	}
	cs.uniques = removedDuplicated
	return cs
}

func (cs *chainState) RunUniqueChecker() {
	if cs.value != nil && len(cs.uniques) > 0 {
		brothers := cs.GetBrothers()
		for _, bro := range brothers {
			if bro.GetValue() == nil {
				continue
			}
			for _, unique := range cs.uniques {
				originKey := cs.GetKey()
				targetKey := bro.GetKey()
				if targetKey == unique && bro.GetValue() == cs.GetValue() {
					msgError := fmt.Errorf("value of '%s' and '%s' fields must be different", originKey, targetKey)
					if cs.CustomMsg != nil && cs.CustomMsg.uniqueNotNil() {
						msgError = buildMessage(*cs.CustomMsg.OnUnique, MessageMeta{
							Field:        &originKey,
							UniqueOrigin: &originKey,
							UniqueTarget: &targetKey,
						})
					}
					cs.AddError(msgError)
				}
			}
		}
	}
	for _, child := range cs.children {
		child.RunUniqueChecker()
	}
}

func (cs *chainState) RunManipulator() (err error) {
	return cs.runManipulate()
}

func (cs *chainState) runManipulate() (err error) {
	// check if current value is not nil and has manipulator
	if cs.value != nil && cs.manipulator != nil {
		// run manipulator
		cs.value, err = (*cs.manipulator)(cs.value)
		if err != nil {
			return err
		}
	}
	for _, child := range cs.children {
		err = child.runManipulate()
		if err != nil {
			return err
		}
	}
	return
}

func (cs *chainState) GetAllKeys() []string {
	var keys []string
	cs.collectKeys(&keys)
	return keys
}

func (cs *chainState) collectKeys(keys *[]string) {
	*keys = append(*keys, cs.key)
	for _, child := range cs.children {
		child.collectKeys(keys)
	}
}

func (cs *chainState) GetResult() ChainResultType {
	return cs
}

func newChainer() ChainerType {
	return &chainState{}
}

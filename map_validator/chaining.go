package map_validator

import (
	"fmt"
)

type ChainResultType interface {
	GetAllKeys() []string
	PrintHierarchyWithSeparator(separator string, currentPath string)
	ToMap() map[string]interface{}
	RunManipulator() error
}

type ChainerType interface {
	GetParentKey() string
	Next(index int) ChainerType
	Back() ChainerType
	Forward(index int) ChainerType
	SetKey(name string) ChainerType
	GetKey() string
	GetParentKeys() []string
	AddChild() ChainerType
	LoadFromMap(data map[string]interface{})
	SetValue(value interface{}) ChainerType
	GetValue() interface{}
	SetManipulator(manipulator *func(interface{}) (interface{}, error)) ChainerType
	GetResult() ChainResultType
}

type chainState struct {
	key         string
	manipulator *func(interface{}) (interface{}, error)
	value       interface{}
	parent      *chainState
	children    []*chainState
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

func (cs *chainState) RunManipulator() (err error) {
	return cs.runManipulate()
}

func (cs *chainState) runManipulate() (err error) {
	if cs.value != nil && cs.manipulator != nil {
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

package map_validator

import (
	"testing"
)

func TestSetKey(t *testing.T) {
	root := newChainer().SetKey("root")
	if root.GetKey() != "root" {
		t.Errorf("Expected name 'root', got '%s'", root.GetKey())
	}
}

func TestNextAndBack(t *testing.T) {
	root := newChainer().SetKey("root")
	child := root.AddChild().SetKey("child")
	if child.GetParentKey() != "root" {
		t.Errorf("Expected parent name 'root', got '%s'", child.GetParentKey())
	}
	if parent := child.Back().GetKey(); parent != "root" {
		t.Errorf("Expected name 'root' after moving back, got '%s'", parent)
	}
}

func TestMultipleChildren(t *testing.T) {
	root := newChainer().SetKey("root")
	child1 := root.AddChild().SetKey("child1")
	child2 := root.AddChild().SetKey("child2")

	if child1.GetKey() != "child1" {
		t.Errorf("Expected first child 'child1', got '%s'", child1.GetKey())
	}

	if child2.GetKey() != "child2" {
		t.Errorf("Expected second child 'child2', got '%s'", child2.GetKey())
	}

	if child := root.Next(0); child.GetKey() != "child1" {
		t.Errorf("Expected first child 'child1', got '%s'", child.GetKey())
	}
	if child := root.Next(1); child.GetKey() != "child2" {
		t.Errorf("Expected second child 'child2', got '%s'", child.GetKey())
	}
}

func TestGetParentKeys(t *testing.T) {
	root := newChainer().SetKey("root")
	child := root.AddChild().SetKey("child")
	grandchild := child.AddChild().SetKey("grandchild")

	names := grandchild.GetParentKeys()
	if len(names) != 2 {
		t.Errorf("Expected 2 parent names, got %d", len(names))
	}
	if names[0] != "child" {
		t.Errorf("Expected first parent name 'child', got '%s'", names[0])
	}
	if names[1] != "root" {
		t.Errorf("Expected second parent name 'root', got '%s'", names[1])
	}
}

func TestForwardAndBack(t *testing.T) {
	root := newChainer().SetKey("root")
	child := root.AddChild().SetKey("child")
	grandchild := child.AddChild().SetKey("grandchild")

	if grandchild.Back().GetKey() != "child" {
		t.Errorf("Expected to move back to 'child', got '%s'", grandchild.Back().GetKey())
	}

	if child.Forward(0).GetKey() != "grandchild" {
		t.Errorf("Expected to move forward to 'grandchild', got '%s'", child.Forward(0).GetKey())
	}
}

func TestPanicOnInvalidIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on invalid child index, but did not occur")
		}
	}()
	root := newChainer().SetKey("root")
	root.Next(0) // This should cause a panic
}

func TestMultipleLevelChildren(t *testing.T) {
	root := newChainer().SetKey("root")
	child1 := root.AddChild().SetKey("child1")
	child2 := root.AddChild().SetKey("child2")

	child1_1 := child1.AddChild().SetKey("child1_1")
	child2_1 := child2.AddChild().SetKey("child2_1")

	if child1_1.GetParentKey() != "child1" {
		t.Errorf("Expected parent name 'child1' for child1_1, got '%s'", child1_1.GetParentKey())
	}

	if child2_1.GetParentKey() != "child2" {
		t.Errorf("Expected parent name 'child2' for child2_1, got '%s'", child2_1.GetParentKey())
	}

	if child1.Forward(0).GetKey() != "child1_1" {
		t.Errorf("Expected to move forward to 'child1_1', got '%s'", child1.Forward(0).GetKey())
	}

	if child2.Forward(0).GetKey() != "child2_1" {
		t.Errorf("Expected to move forward to 'child2_1', got '%s'", child2.Forward(0).GetKey())
	}

	root.GetResult().GetAllKeys()
}

func TestChainValues(T *testing.T) {
	root := newChainer().SetKey("root")
	childa_1 := root.AddChild().SetKeyValue("childa_1", "value+childa_1")
	childa_2 := root.AddChild().SetKeyValue("childa_2", "value+childa_2")
	root.AddChild().SetKeyValue("childa_3", "value+childa_3")

	childa_1.AddChild().SetKeyValue("childb_1_d", "value+childb_d")
	childa_2.AddChild().SetKeyValue("childa_2_x", "value+a").SetUniques([]string{"childa_2_z", "childa_2_y"})
	childa_2.AddChild().SetKeyValue("childa_2_z", "value+a")
	childa_2.AddChild().SetKeyValue("childa_2_y", "value+ka")
	childa_2.AddChild().SetKeyValue("childa_2_sssy", nil)
	childa_2.AddChild().SetKeyValue("childa_2_m", "value+a")
	childa_2.AddChild().SetKeyValue("childa_2_g", "value+childa_2_g")
	childa_2.AddChild().SetKeyValue("childa_2_s", "value+childa_2_s")
	childa_2.AddChild().SetKeyValue("childa_2_e", "value+childa_2_e")

	childa_1.AddChild().SetKeyValue("childa_1_e", "value+childa_1_e")
	childa_1.AddChild().SetKeyValue("childa_1_f", "value+childa_1_f").SetUniques([]string{"childa_1_t"})
	childa_1.AddChild().SetKeyValue("childa_1_t", "value+childa_1_e").SetUniques([]string{"childa_1_e"})

	root.GetResult().RunUniqueChecker()
	res := root.GetResult()
	errors := res.GetErrors()
	if len(errors) != 2 {
		T.Errorf("Expected have two errors, but we got %d error", len(errors))
	}
}

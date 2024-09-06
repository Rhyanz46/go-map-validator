package map_validator

import "testing"

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

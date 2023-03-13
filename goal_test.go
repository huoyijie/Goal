package goal

import "testing"

func TestGoal(t *testing.T) {
	menuList := groupList()
	if len(menuList) == 0 {
		t.Errorf("invalid menu %d", len(menuList))
	}
}

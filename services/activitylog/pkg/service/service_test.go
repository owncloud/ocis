package service

import (
	"testing"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestAddActivity(t *testing.T) {
	testCases := []struct {
		Name       string
		Tree       map[string]*provider.ResourceInfo
		Activities map[string]string
		Expected   map[string][]Activity
	}{
		{
			Name: "simple",
			Tree: map[string]*provider.ResourceInfo{
				"base":    resourceInfo("base", "parent"),
				"parent":  resourceInfo("parent", "spaceid"),
				"spaceid": resourceInfo("spaceid", "spaceid"),
			},
			Activities: map[string]string{
				"activity": "base",
			},
			Expected: map[string][]Activity{
				"base":    activitites("activity", 0),
				"parent":  activitites("activity", 1),
				"spaceid": activitites("activity", 2),
			},
		},
		{
			Name: "two activities on same resource",
			Tree: map[string]*provider.ResourceInfo{
				"base":    resourceInfo("base", "parent"),
				"parent":  resourceInfo("parent", "spaceid"),
				"spaceid": resourceInfo("spaceid", "spaceid"),
			},
			Activities: map[string]string{
				"activity1": "base",
				"activity2": "base",
			},
			Expected: map[string][]Activity{
				"base":    activitites("activity1", 0, "activity2", 0),
				"parent":  activitites("activity1", 1, "activity2", 1),
				"spaceid": activitites("activity1", 2, "activity2", 2),
			},
		},
		{
			Name: "two activities on different resources",
			Tree: map[string]*provider.ResourceInfo{
				"base1":   resourceInfo("base1", "parent"),
				"base2":   resourceInfo("base2", "parent"),
				"parent":  resourceInfo("parent", "spaceid"),
				"spaceid": resourceInfo("spaceid", "spaceid"),
			},
			Activities: map[string]string{
				"activity1": "base1",
				"activity2": "base2",
			},
			Expected: map[string][]Activity{
				"base1":   activitites("activity1", 0),
				"base2":   activitites("activity2", 0),
				"parent":  activitites("activity1", 1, "activity2", 1),
				"spaceid": activitites("activity1", 2, "activity2", 2),
			},
		},
		{
			Name: "more elaborate resource tree",
			Tree: map[string]*provider.ResourceInfo{
				"base1":   resourceInfo("base1", "parent1"),
				"base2":   resourceInfo("base2", "parent1"),
				"parent1": resourceInfo("parent1", "spaceid"),
				"base3":   resourceInfo("base3", "parent2"),
				"parent2": resourceInfo("parent2", "spaceid"),
				"spaceid": resourceInfo("spaceid", "spaceid"),
			},
			Activities: map[string]string{
				"activity1": "base1",
				"activity2": "base2",
				"activity3": "base3",
			},
			Expected: map[string][]Activity{
				"base1":   activitites("activity1", 0),
				"base2":   activitites("activity2", 0),
				"base3":   activitites("activity3", 0),
				"parent1": activitites("activity1", 1, "activity2", 1),
				"parent2": activitites("activity3", 1),
				"spaceid": activitites("activity1", 2, "activity2", 2, "activity3", 2),
			},
		},
		{
			Name: "different depths within one resource",
			Tree: map[string]*provider.ResourceInfo{
				"base1":   resourceInfo("base1", "parent1"),
				"parent1": resourceInfo("parent1", "parent2"),
				"base2":   resourceInfo("base2", "parent2"),
				"parent2": resourceInfo("parent2", "parent3"),
				"base3":   resourceInfo("base3", "parent3"),
				"parent3": resourceInfo("parent3", "spaceid"),
				"spaceid": resourceInfo("spaceid", "spaceid"),
			},
			Activities: map[string]string{
				"activity1": "base1",
				"activity2": "base2",
				"activity3": "base3",
				"activity4": "parent2",
			},
			Expected: map[string][]Activity{
				"base1":   activitites("activity1", 0),
				"base2":   activitites("activity2", 0),
				"base3":   activitites("activity3", 0),
				"parent1": activitites("activity1", 1),
				"parent2": activitites("activity1", 2, "activity2", 1, "activity4", 0),
				"parent3": activitites("activity1", 3, "activity2", 2, "activity3", 1, "activity4", 1),
				"spaceid": activitites("activity1", 4, "activity2", 3, "activity3", 2, "activity4", 2),
			},
		},
	}

	for _, tc := range testCases {
		alog := &ActivitylogService{
			store: store.Create(),
		}

		getResource := func(ref *provider.Reference) (*provider.ResourceInfo, error) {
			return tc.Tree[ref.GetResourceId().GetOpaqueId()], nil
		}

		for k, v := range tc.Activities {
			err := alog.addActivity(reference(v), k, time.Time{}, getResource)
			require.NoError(t, err)
		}

		for id, acts := range tc.Expected {
			activities, err := alog.Activities(resourceID(id))
			require.NoError(t, err, tc.Name+":"+id)
			require.ElementsMatch(t, acts, activities, tc.Name+":"+id)
		}
	}
}

func activitites(acts ...interface{}) []Activity {
	var activities []Activity
	act := Activity{}
	for _, a := range acts {
		switch v := a.(type) {
		case string:
			act.EventID = v
		case int:
			act.Depth = v
			activities = append(activities, act)
		}
	}
	return activities
}

func resourceID(id string) *provider.ResourceId {
	return &provider.ResourceId{
		StorageId: "storageid",
		OpaqueId:  id,
		SpaceId:   "spaceid",
	}
}

func reference(id string) *provider.Reference {
	return &provider.Reference{ResourceId: resourceID(id)}
}

func resourceInfo(id, parentID string) *provider.ResourceInfo {
	return &provider.ResourceInfo{
		Id:       resourceID(id),
		ParentId: resourceID(parentID),
		Space: &provider.StorageSpace{
			Root: resourceID("spaceid"),
		},
	}
}

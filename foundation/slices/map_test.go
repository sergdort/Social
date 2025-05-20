package slices

import (
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("maps integers to strings", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}

		result := Map(input, func(i int) string {
			return string(rune('0' + i))
		})

		if len(result) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("at index %d: expected %s, got %s", i, expected[i], result[i])
			}
		}
	})

	t.Run("maps strings to integers", func(t *testing.T) {
		input := []string{"1", "2", "3", "4", "5"}
		expected := []int{1, 2, 3, 4, 5}

		result := Map(input, func(s string) int {
			return int(s[0] - '0')
		})

		if len(result) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], result[i])
			}
		}
	})

	t.Run("maps empty slice", func(t *testing.T) {
		input := []int{}
		result := Map(input, func(i int) string {
			return string(rune('0' + i))
		})

		if len(result) != 0 {
			t.Errorf("expected empty slice, got length %d", len(result))
		}
	})

	t.Run("maps with complex transformation", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		input := []Person{
			{"Alice", 25},
			{"Bob", 30},
			{"Charlie", 35},
		}

		expected := []string{"Alice (25)", "Bob (30)", "Charlie (35)"}

		result := Map(input, func(p Person) string {
			return p.Name + " (" + string(rune('0'+p.Age/10)) + string(rune('0'+p.Age%10)) + ")"
		})

		if len(result) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("at index %d: expected %s, got %s", i, expected[i], result[i])
			}
		}
	})
}

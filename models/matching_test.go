package models

import "testing"

func TestMatchingIsMatch_RespectsBlockUsers(t *testing.T) {
	tests := []struct {
		name  string
		user1 Matching
		user2 Matching
		want  bool
	}{
		{
			name: "user1 blocks user2",
			user1: Matching{
				UserID: 1,
				BlockUsers: []BlockUser{
					{UserID: 2},
				},
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			user2: Matching{
				UserID:          2,
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			want: false,
		},
		{
			name: "user2 blocks user1",
			user1: Matching{
				UserID:          1,
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			user2: Matching{
				UserID: 2,
				BlockUsers: []BlockUser{
					{UserID: 1},
				},
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			want: false,
		},
		{
			name: "no block and software compatible",
			user1: Matching{
				UserID:          1,
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			user2: Matching{
				UserID:          2,
				OnlineSoftwares: []OnlineSoftware{{Name: "app", Type: 0}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user1.IsMatch(tt.user2); got != tt.want {
				t.Fatalf("IsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

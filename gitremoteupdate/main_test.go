package gitremoteupdate_test

import (
	"reflect"
	"testing"

	"github.com/shurcooL/go/gitremoteupdate"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		stderr []byte
		want   gitremoteupdate.Result
	}{
		{
			stderr: []byte(`From https://example.com/user/repo.git
   e8569f7..de0ad17  master     -> master
 * [new branch]      new-branch -> new-branch
`),
			want: gitremoteupdate.Result{
				Changes: []gitremoteupdate.Change{
					{Op: gitremoteupdate.Updated, Branch: "master"},
					{Op: gitremoteupdate.New, Branch: "new-branch"},
				},
			},
		},

		{
			stderr: []byte(`From https://example.com/user/repo.git
   990cfc0..a65b539  foo-branch -> foo-branch
   d6d0813..e8569f7  master     -> master
`),
			want: gitremoteupdate.Result{
				Changes: []gitremoteupdate.Change{
					{Op: gitremoteupdate.Updated, Branch: "foo-branch"},
					{Op: gitremoteupdate.Updated, Branch: "master"},
				},
			},
		},

		{
			stderr: []byte(`From https://example.com/user/repo.git
 x [deleted]         (none)     -> master-backup-new-branch
`),
			want: gitremoteupdate.Result{
				Changes: []gitremoteupdate.Change{
					{Op: gitremoteupdate.Deleted, Branch: "master-backup-new-branch"},
				},
			},
		},

		{
			stderr: []byte(`From https://example.com/user/repo.git
 * [new branch]      another-looooooong-branch-wheeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee -> another-looooooong-branch-wheeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
   de0ad17..143ee1d  master     -> master
`),
			want: gitremoteupdate.Result{
				Changes: []gitremoteupdate.Change{
					{Op: gitremoteupdate.New, Branch: "another-looooooong-branch-wheeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"},
					{Op: gitremoteupdate.Updated, Branch: "master"},
				},
			},
		},
	} {
		got, err := gitremoteupdate.Parse(tc.stderr)
		if err != nil {
			t.Errorf("got non-nil error: %v", err)
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %v, want %v", got, tc.want)
			continue
		}
	}
}

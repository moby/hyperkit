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
			stderr: []byte(""),
			want:   gitremoteupdate.Result{},
		},

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

		{
			stderr: []byte(`From https://example.com/user/repo.git
 * [new branch]      gofmt-circleci -> gofmt-circleci
   0bccbc3..fb8ec00  master     -> master
 + ca1c467...939c2da refs/pull/291/merge -> refs/pull/291/merge  (forced update)
 + 2ca958e...2cf60d2 refs/pull/334/merge -> refs/pull/334/merge  (forced update)
 * [new ref]         refs/pull/338/head -> refs/pull/338/head
 * [new ref]         refs/pull/344/head -> refs/pull/344/head
 * [new ref]         refs/pull/344/merge -> refs/pull/344/merge
`),
			want: gitremoteupdate.Result{
				Changes: []gitremoteupdate.Change{
					{Op: gitremoteupdate.New, Branch: "gofmt-circleci"},
					{Op: gitremoteupdate.Updated, Branch: "master"},
					{Op: gitremoteupdate.Updated, Branch: "refs/pull/291/merge"},
					{Op: gitremoteupdate.Updated, Branch: "refs/pull/334/merge"},
					{Op: gitremoteupdate.New, Branch: "refs/pull/338/head"},
					{Op: gitremoteupdate.New, Branch: "refs/pull/344/head"},
					{Op: gitremoteupdate.New, Branch: "refs/pull/344/merge"},
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

package sv_member

import (
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/domain/infra/api_context"
	"github.com/watermint/toolbox/domain/infra/api_parser"
	"github.com/watermint/toolbox/domain/infra/api_test"
	"github.com/watermint/toolbox/domain/model/mo_member"
	"strings"
	"testing"
)

const (
	memberGetInfoSuccessRes = `[
  {
    ".tag": "member_info",
    "profile": {
      "team_member_id": "dbmid:xxxxxxxxxxxxxx-xxxxxxxx-xxxx-xxxxxx",
      "external_id": "xxxxx x+xxxxxxx.xxx-xxxxx@xxxxxxxxx.xxx",
      "account_id": "xxxx:xxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx",
      "email": "xxx+xxx@xxxxxxxxx.xxx",
      "email_verified": true,
      "status": {
        ".tag": "active"
      },
      "name": {
        "given_name": "デモ",
        "surname": "Dropbox",
        "familiar_name": "デモ",
        "display_name": "デモ Dropbox",
        "abbreviated_name": "デD"
      },
      "membership_type": {
        ".tag": "full"
      },
      "joined_on": "2016-01-15T05:42:49Z",
      "groups": [
        "g:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "g:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "g:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "g:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "g:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
      ],
      "member_folder_id": "xxxxxxxxxx"
    },
    "role": {
      ".tag": "team_admin"
    }
  }
]`
)

func TestMemberImpl_Resolve(t *testing.T) {
	member := &mo_member.Member{}
	j := gjson.Parse(memberGetInfoSuccessRes)
	if !j.IsArray() {
		t.Error("invalid")
	}
	a := j.Array()[0]
	err := api_parser.ParseModel(member, a)
	if err != nil {
		t.Error("invalid")
	}
	if member.TeamMemberId != "dbmid:xxxxxxxxxxxxxx-xxxxxxxx-xxxx-xxxxxx" {
		t.Error("invalid")
	}
}

func TestMemberImpl_ResolveByEmail(t *testing.T) {
	api_test.DoTestBusinessInfo(func(ctx api_context.Context) {
		svm := New(ctx)
		members, err := svm.List()
		if err != nil {
			t.Error(err)
		}

		for i, member := range members {
			if i > 10 {
				break
			}
			m1, err := svm.Resolve(member.TeamMemberId)
			if err != nil {
				t.Error(err)
			}
			m2, err := svm.ResolveByEmail(member.Email)
			if err != nil {
				t.Error(err)
			}

			if m1.TeamMemberId != member.TeamMemberId ||
				m2.TeamMemberId != member.TeamMemberId {
				t.Error("invalid")
			}
		}

		_, err = svm.Resolve("dbmid:xxxxxxxxxxxxxx-xxxxxxxx-xxxx-xxxxxx")
		if err == nil {
			t.Error("invalid")
		}

		_, err = svm.ResolveByEmail("non_existent@example.com")
		if err == nil {
			t.Error("invalid")
		}
	})
}

func TestMemberImpl_ListResolve(t *testing.T) {
	api_test.DoTestBusinessInfo(func(ctx api_context.Context) {
		ls := newTest(ctx)
		members, err := ls.List()
		if err != nil {
			t.Error(err)
			return
		}
		if len(members) < 1 {
			t.Error("invalid")
		}
		if !strings.Contains(members[0].Email, "@") {
			t.Error("invalid")
		}

		m, err := ls.Resolve(members[0].TeamMemberId)
		if err != nil {
			t.Error("failed fetch")
		}
		if m.TeamMemberId != members[0].TeamMemberId {
			t.Error("invalid")
		}
		if m.Email != members[0].Email {
			t.Error("invalid")
		}
	})
}

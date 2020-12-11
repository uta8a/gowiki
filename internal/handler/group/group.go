package group

import (
  "database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/session"
  "net/http"
  "io/ioutil"
)

type PostReq struct {
  GroupName string `json:"group_name"`
  GroupMembers []string `json:"group_members"`
}
type PostRes struct {
  GroupName string `json:"group_name"`
}
type GetRes struct {
  UserGroups []string `json:"user_groups"`
}

func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  // GET
  if r.Method == http.MethodGet {
    err := getGroups(db, gs, w, r)
    if err != nil {
      return err
    }
    return nil
  }
	// POST
	if r.Method == http.MethodPost {
		err := registerGroup(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
  }
	http.NotFound(w, r)
	return nil
}

func getGroups(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  // Authentication
  ok := gs.SessionCheck(w, r)
  if !ok {
    http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
  }
  sess := gs.SessionStart(w, r)
  username := sess.Get("username")
  
  // groups
  query := fmt.Sprintf("SELECT group_name FROM group_users WHERE group_user = '%s'", username)
  rows, err := db.Query(query)
  if err != nil {
    return err
  }

  var userGroup []string
  for rows.Next() {
    var groupName string
    if err := rows.Scan(&groupName); err != nil {
      return err
    }
    userGroup = append(userGroup, groupName)
  }

  // response
  response := GetRes{UserGroups: userGroup}
  res, err := json.Marshal(response)
  if err != nil {
    return err
  }
  w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func registerGroup(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
  // Authentication
  ok := gs.SessionCheck(w, r)
  if !ok {
    http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
  }
  sess := gs.SessionStart(w, r)
  // post user is group admin
  username := sess.Get("username")

  // request body groupname, group members validation
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return err
  }
  var req PostReq
  if err := json.Unmarshal(body, &req); err != nil {
    return err
  }
  groupname := req.GroupName
  groupMembers := req.GroupMembers
  err = validateGroupName(groupname)
  if err != nil {
    return err
  }
  for _, m := range groupMembers {
    err := validateUsername(m)
    if err != nil {
      return err
    }
  }
  // groupname identity check
  var exists bool
  query := fmt.Sprintf("SELECT EXISTS (SELECT group_name FROM group_admins WHERE group_name = '%s')", groupname)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("group name %s already exists", groupname)
  }
  
  // username exists check
  // TODO? 仕様の話になりそう。未来のユーザをあらかじめ登録できるのは便利かもしれないけどどうなんだろう
  // DBからとってくるときに困りそうな気もする

  // set group_admin to db
  query = fmt.Sprintf("INSERT INTO group_admins(group_name, group_admin) VALUES('%s', '%s')", groupname, username)
  _, err = db.Exec(query)
  if err != nil {
    return err
  }
  // set group members to db
  query = fmt.Sprintf("INSERT INTO group_users(group_name, group_user) VALUES('%s', '%s')", groupname, username)
  _, err = db.Exec(query)
  if err != nil {
    return err
  }
  for _, member := range groupMembers {
    query = fmt.Sprintf("INSERT INTO group_users(group_name, group_user) VALUES('%s', '%s')", groupname, member)
    _, err = db.Exec(query)
    if err != nil {
      return err
    }
  }
  // response
  response := PostRes{GroupName: groupname}
  res, err := json.Marshal(response)
  if err != nil {
    return err
  }
  w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func validateGroupName(groupname string) error {
	// group name char: [a-zA-Z0-9-_]
	const minGroupnameLength = 2
	const maxGroupnameLength = 20

  // それぞれの文字を調べるコストより、lenを先に調べて文字を減らす
	if (len(groupname) < minGroupnameLength || len(groupname) > maxGroupnameLength) {
		return fmt.Errorf("group name is invalid")
	}
	for _, ch := range groupname {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && (ch != '_') && (ch != '-') {
			return fmt.Errorf("group name is invalid")
		}
	}
	return nil
}

func validateUsername(username string) error {
	// username char: [a-zA-Z0-9-_]
	const minUsernameLength = 2
	const maxUsernameLength = 20

  // それぞれの文字を調べるコストより、lenを先に調べて文字を減らす
	if (len(username) < minUsernameLength || len(username) > maxUsernameLength) {
		return fmt.Errorf("member username is invalid")
	}
	for _, ch := range username {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && (ch != '_') && (ch != '-') {
			return fmt.Errorf("member username is invalid")
		}
	}
	return nil
}

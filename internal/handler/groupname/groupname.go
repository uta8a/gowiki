package groupname

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/session"
	"io/ioutil"
	"net/http"
  "strings"
  "log"
)

type GetRes struct {
	GroupName    string   `json:"group_name"`
	GroupMembers []string `json:"group_members"`
}
type PostReq struct {
	GroupName    string   `json:"group_name"`
	GroupMembers []string `json:"group_members"`
}

type PostRes struct {
	GroupName string `json:"group_name"`
}

func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// Path split, get articleId: string
  groupName, err := getParamFromURL("/groups/", r.URL.Path)
  log.Printf("groupname: %s, %s",groupName, r.URL.Path)
	if err != nil {
		return err
	}
	// GET
	if r.Method == http.MethodGet {
		err := getGroup(db, gs, w, r, groupName)
		if err != nil {
			return err
		}
		return nil
	}
	// POST
	if r.Method == http.MethodPost {
		err := updateGroup(db, gs, w, r, groupName)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

func getParamFromURL(base string, u string) (string, error) {
	rawParam := strings.TrimPrefix(u, base)
	return rawParam, nil
}

// func getArticle(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request, id int) error {
// 	// Authentication
// 	ok := gs.SessionCheck(w, r)
// 	if !ok {
// 		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
// 		return nil
// 	}
// 	sess := gs.SessionStart(w, r)
// 	username := sess.Get("username").(string)

// 	// check user is in group or not
// 	// get users in article's usergroup, compare
// 	var users []string
// 	stmt, err := db.Prepare("SELECT group_user FROM group_users WHERE group_name IN (SELECT group_name FROM articles WHERE article_id = $1)")
// 	if err != nil {
// 		return err
// 	}
// 	rows, err := stmt.Query(id)
// 	if err != nil {
// 		return err
// 	}
// 	for rows.Next() {
// 		var user string
// 		if err := rows.Scan(&user); err != nil {
// 			return err
// 		}
// 		users = append(users, user)
// 	}
// 	_, ok = Find(users, username)
// 	if !ok {
// 		return fmt.Errorf("Forbidden: you cannot see this article")
// 	}
// 	// get article, tag
// 	var (
// 		articleId   int
// 		title       string
// 		articlePath string
// 		groupName   string
// 		body        string
// 		tags        []string
// 	)
// 	stmt, err = db.Prepare("SELECT article_id, title, article_path, group_name, body FROM articles WHERE article_id = $1")
// 	if err != nil {
// 		return err
// 	}
// 	err = stmt.QueryRow(id).Scan(&articleId, &title, &articlePath, &groupName, &body)
// 	if err != nil {
// 		return err
// 	}
// 	stmt, err = db.Prepare("SELECT tag FROM tags WHERE article_id = $1")
// 	if err != nil {
// 		return err
// 	}
// 	rows, err = stmt.Query(articleId)
// 	for rows.Next() {
// 		var tag string
// 		if err := rows.Scan(&tag); err != nil {
// 			return err
// 		}
// 		tags = append(tags, tag)
// 	}
// 	// response
// 	response := GetRes{ArticleId: articleId, Title: title, ArticlePath: articlePath, Tags: tags, GroupName: groupName, Body: body}
// 	res, err := json.Marshal(response)
// 	if err != nil {
// 		return err
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintf(w, string(res))
// 	return nil
// }

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func getGroup(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request, name string) error {
	// Authentication
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	username := sess.Get("username").(string)

	// validate user is in GroupName ?
	stmt, err := db.Prepare("SELECT group_user FROM group_users WHERE group_name = $1")
	if err != nil {
		return err
  }
  log.Printf("%s", name)
	rows, err := stmt.Query(name)
	if err != nil {
		return err
	}
	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil
		}
		users = append(users, user)
  }
  log.Printf("%v", users)
	_, ok = Find(users, username)
	if !ok {
		return fmt.Errorf("You're not in this Group")
	}
	// response
	response := GetRes{GroupName: name, GroupMembers: users}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func updateGroup(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request, name string) error {
	// Authentication
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	username := sess.Get("username").(string)

	// request parsing
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var req PostReq
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return err
	}
	// validate GroupName's admin is username?
	groupname := req.GroupName
	var admin string
	stmt, err := db.Prepare("SELECT group_admin FROM group_admins WHERE group_name = $1")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(groupname).Scan(&admin)
	if err != nil {
		return err
	}
	if username != admin {
		return fmt.Errorf("You're not admin of Group %s", groupname)
	}
	// delete all users which relates GroupName
	stmt, err = db.Prepare("DELETE FROM group_users WHERE group_name = $1 AND group_user <> $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(groupname, admin)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("INSERT INTO group_users(group_name, group_user) VALUES($1, $2)")
	if err != nil {
		return err
  }
  groupMembers := []string{}
	m := make(map[string]struct{})
	for _, mem := range req.GroupMembers {
		m[mem] = struct{}{}
	}
	for i := range m {
    if i != admin {
      groupMembers = append(groupMembers, i)
    }
	}
	for _, m := range groupMembers {
		_, err = stmt.Exec(groupname, m)
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

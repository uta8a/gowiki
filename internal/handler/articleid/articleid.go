package articleid

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/session"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GetRes struct {
	ArticleId   int      `json:"article_id"`
	Title       string   `json:"title"`
	ArticlePath string   `json:"article_path"`
	Tags        []string `json:"tags"`
	GroupName   string   `json:"group_name"`
	Body        string   `json:"body"`
}

type PostReq struct {
	Title       string   `json:"title"`
	ArticlePath string   `json:"article_path"`
	Tags        []string `json:"tags"`
	GroupName   string   `json:"group_name"`
	Body        string   `json:"body"`
}

type PostRes struct {
	ArticleId int `json:"article_id"`
}

func New(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// Path split, get articleId: int
	articleId, err := getParamFromURL("/articles/", r.URL.Path)
	if err != nil {
		return err
	}
	// GET
	if r.Method == http.MethodGet {
		err := getArticle(db, gs, w, r, articleId)
		if err != nil {
			return err
		}
		return nil
	}
	// POST
	if r.Method == http.MethodPost {
		err := updateArticle(db, gs, w, r, articleId)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

func getParamFromURL(base string, u string) (int, error) {
	rawParam := strings.TrimLeft(u, base)
	// db.articles_id is serial, 32 bit 1-index unsigned integer
	param, err := strconv.ParseUint(rawParam, 10, 32)
	if err != nil {
		return -1, err
	}
	// id validation check
	if param == 0 {
		return -1, fmt.Errorf("article_id: 0 is not found")
	}
	return int(param), nil
}

func getArticle(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request, id int) error {
	// Authentication
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	username := sess.Get("username").(string)

	// check user is in group or not
	// get users in article's usergroup, compare
	var users []string
	stmt, err := db.Prepare("SELECT group_user FROM group_users WHERE group_name IN (SELECT group_name FROM articles WHERE article_id = $1)")
	if err != nil {
		return err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return err
		}
		users = append(users, user)
	}
	_, ok = Find(users, username)
	if !ok {
		return fmt.Errorf("Forbidden: you cannot see this article")
	}
	// get article, tag
	var (
		articleId   int
		title       string
		articlePath string
		groupName   string
		body        string
		tags        []string
	)
	stmt, err = db.Prepare("SELECT article_id, title, article_path, group_name, body FROM articles WHERE article_id = $1")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(id).Scan(&articleId, &title, &articlePath, &groupName, &body)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("SELECT tag FROM tags WHERE article_id = $1")
	if err != nil {
		return err
	}
	rows, err = stmt.Query(articleId)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	// response
	response := GetRes{ArticleId: articleId, Title: title, ArticlePath: articlePath, Tags: tags, GroupName: groupName, Body: body}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func updateArticle(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request, id int) error {
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
	// validate GroupName Exists?
	groupname := req.GroupName
	var exists bool
	stmt, err := db.Prepare("SELECT EXISTS (SELECT group_name FROM group_admins WHERE group_name = $1)")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(groupname).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("group name %s does not exist", groupname)
	}
	// validate group name is valid ?
	// check user is in origin Group
	var users []string
	stmt, err = db.Prepare("SELECT group_user FROM group_users WHERE group_name IN (SELECT group_name FROM articles WHERE article_id = $1)")
	if err != nil {
		return err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return err
		}
		users = append(users, user)
	}
	_, ok = Find(users, username)
	if !ok {
		return fmt.Errorf("Forbidden: you cannot see this article")
	}

	// author user is in target group?
	stmt, err = db.Prepare("SELECT group_name FROM group_users WHERE group_user = $1")
	if err != nil {
		return err
	}
	rows, err = stmt.Query(username)
	if err != nil {
		return err
	}
	var groupUsers []string
	if rows.Next() {
		var groupUser string
		if err := rows.Scan(&groupUser); err != nil {
			return err
		}
		groupUsers = append(groupUsers, groupUser)
	}
	_, ok = Find(groupUsers, groupname)
	if !ok {
		return fmt.Errorf("Forbidden: you cannot move this article to specify group")
	}

	// create articles
	title := req.Title
	articlePath := req.ArticlePath
	tags := req.Tags
	body := req.Body
	// insert articles
	var articleId int
	stmt, err = db.Prepare("UPDATE articles SET title = $1, article_path = $2, group_name = $3, body = $4 WHERE article_id = $5 RETURNING article_id")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(title, articlePath, groupname, body, id).Scan(&articleId)
	if err != nil {
		return err
	}
	// insert tags
	// delete all tag which relates article_id, then insert
	stmt, err = db.Prepare("DELETE FROM tags WHERE article_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(articleId)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("INSERT INTO tags(article_id, tag) VALUES($1, $2)")
	if err != nil {
		return err
	}
	for _, tag := range tags {
		_, err = stmt.Exec(articleId, tag)
		if err != nil {
			return err
		}
	}
	// response
	response := PostRes{ArticleId: articleId}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

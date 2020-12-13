package article

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/suburi-dev/gowiki/internal/session"
	"io/ioutil"
	"net/http"
)

type GetRes struct {
	GroupNumber int     `json:"group_number"`
	Groups      []Group `json:"groups"`
}
type Group struct {
	GroupName  string `json:"group_name"`
	ArticlesId []int  `json:"articles_id"`
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
	// GET
	if r.Method == http.MethodGet {
		err := getArticles(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	// POST
	if r.Method == http.MethodPost {
		err := createArticle(db, gs, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	http.NotFound(w, r)
	return nil
}

func getArticles(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// Authentication
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	username := sess.Get("username")
	// get articles
	query := fmt.Sprintf("SELECT group_name FROM group_users WHERE group_user = '%s'", username)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	var groups []Group
	for rows.Next() {
		var groupName string
		if err := rows.Scan(&groupName); err != nil {
			return err
		}
		query = fmt.Sprintf("SELECT article_id FROM articles WHERE group_name = '%s'", groupName)
		rows, err := db.Query(query)
		if err != nil {
			return err
		}
		var articlesId []int
		for rows.Next() {
			var articleId int
			if err := rows.Scan(&articleId); err != nil {
				return err
			}
			articlesId = append(articlesId, articleId)
		}
		groups = append(groups, Group{GroupName: groupName, ArticlesId: articlesId})
	}
	groupNumber := len(groups)
	response := GetRes{GroupNumber: groupNumber, Groups: groups}
	res, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(res))
	return nil
}

func createArticle(db *sql.DB, gs *session.Manager, w http.ResponseWriter, r *http.Request) error {
	// Authentication
	ok := gs.SessionCheck(w, r)
	if !ok {
		http.Error(w, "Unauthorized please login", http.StatusUnauthorized)
		return nil
	}
	sess := gs.SessionStart(w, r)
	_ = sess.Get("username")

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
	// create articles
	title := req.Title
	articlePath := req.ArticlePath
	tags := req.Tags
	body := req.Body
	// insert articles
	var articleId int
	stmt, err = db.Prepare("INSERT INTO articles(title, article_path, group_name, body) VALUES($1, $2, $3, $4) RETURNING article_id")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(title, articlePath, groupname, body).Scan(&articleId)
	if err != nil {
		return err
	}
	// insert tags
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

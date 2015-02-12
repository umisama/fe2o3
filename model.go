package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	sb "github.com/umisama/go-sqlbuilder"
	"net/http"
	"time"
)

var db *sql.DB

type QiitaFeed struct {
	Id       string    `json:"id"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	Body     string    `json:"rendered_body"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

var tbl_qiita = sb.NewTable(
	"QIITAFEED",
	sb.StringColumn("id", &sb.ColumnOption{
		PrimaryKey: true,
	}),
	sb.StringColumn("title", &sb.ColumnOption{
		Size:    511,
		NotNull: true,
	}),
	sb.StringColumn("url", &sb.ColumnOption{
		Size:    511,
		NotNull: true,
	}),
	sb.StringColumn("body", &sb.ColumnOption{
		Size:    10000,
		NotNull: true,
	}),
	sb.DateColumn("create_at", &sb.ColumnOption{
		NotNull: true,
	}),
	sb.DateColumn("update_at", &sb.ColumnOption{
		NotNull: true,
	}),
)

type HatenaFeed struct {
	Id       string  `xml:"guid"`
	Title    string  `xml:"title"`
	Url      string  `xml:"link"`
	Body     string  `xml:"description"`
	CreateAt xmlTime `xml:"pubDate"`
}

var tbl_hatena = sb.NewTable(
	"HATENAFEED",
	sb.StringColumn("id", &sb.ColumnOption{
		PrimaryKey: true,
	}),
	sb.StringColumn("title", &sb.ColumnOption{
		Size:    511,
		NotNull: true,
	}),
	sb.StringColumn("url", &sb.ColumnOption{
		Size:    511,
		NotNull: true,
	}),
	sb.StringColumn("body", &sb.ColumnOption{
		Size:    10000,
		NotNull: true,
	}),
	sb.DateColumn("create_at", &sb.ColumnOption{
		NotNull: true,
	}),
)

func (m *HatenaFeed) Save(tx *sql.Tx) error {
	ff, err := GetHatenaFeed(m.Id, tx)
	if err != nil {
		return err
	}
	if ff != nil {
		return nil
	}

	query, args, err := sb.Insert(tbl_hatena).
		Values(m.Id, m.Title, m.Url, m.Body, time.Time(m.CreateAt)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetHatenaFeed(id string, tx *sql.Tx) (*HatenaFeed, error) {
	query, args, err := sb.Select(sb.Star).
		From(tbl_hatena).
		Where(tbl_hatena.C("id").Eq(id)).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(query, args...)

	m := new(HatenaFeed)
	err = row.Scan(&m.Id, &m.Title, &m.Url, &m.Body, &m.CreateAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return m, nil
}

func GetHatenaFeedList() ([]HatenaFeed, error) {
	query, args, err := sb.Select(sb.Star).
		From(tbl_hatena).
		Limit(5).
		OrderBy(true, tbl_hatena.C("create_at")).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	l := make([]HatenaFeed, 0)
	for rows.Next() {
		m := HatenaFeed{}
		err = rows.Scan(&m.Id, &m.Title, &m.Url, &m.Body, &m.CreateAt)
		if err != nil {
			return nil, err
		}
		l = append(l, m)
	}
	return l, nil

}

func GetHatenaFeedListFromInternet() ([]HatenaFeed, error) {
	resp, err := http.Get("http://umisama.hatenablog.com/rss")
	if err != nil {
		return nil, err
	}

	type xmlType struct {
		XMLName xml.Name     `xml:"rss"`
		List    []HatenaFeed `xml:"channel>item"`
	}

	var doc xmlType
	err = xml.NewDecoder(resp.Body).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.List, nil
}

func (m *QiitaFeed) Save(tx *sql.Tx) error {
	ff, err := GetQiitaFeed(m.Id, tx)
	if err != nil {
		return err
	}
	if ff != nil {
		return nil
	}

	query, args, err := sb.Insert(tbl_qiita).
		Values(m.Id, m.Title, m.Url, m.Body, m.CreateAt, m.UpdateAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetQiitaFeed(id string, tx *sql.Tx) (*QiitaFeed, error) {
	query, args, err := sb.Select(sb.Star).
		From(tbl_qiita).
		Where(tbl_qiita.C("id").Eq(id)).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(query, args...)

	m := new(QiitaFeed)
	err = row.Scan(&m.Id, &m.Title, &m.Url, &m.Body, &m.CreateAt, &m.UpdateAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return m, nil
}

func GetQiitaFeedList() ([]QiitaFeed, error) {
	query, args, err := sb.Select(sb.Star).
		From(tbl_qiita).
		Limit(5).
		OrderBy(true, tbl_qiita.C("create_at")).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	l := make([]QiitaFeed, 0)
	for rows.Next() {
		m := QiitaFeed{}
		err = rows.Scan(&m.Id, &m.Title, &m.Url, &m.Body, &m.CreateAt, &m.UpdateAt)
		if err != nil {
			return nil, err
		}
		l = append(l, m)
	}
	return l, nil

}

func GetQiitaFeedListFromInternet() ([]QiitaFeed, error) {
	resp, err := http.Get("http://qiita.com/api/v2/users/umisama/items")
	if err != nil {
		return nil, err
	}

	var list []QiitaFeed
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	sb.SetDialect(sb.SqliteDialect{})

	query, args, err := sb.CreateTable(tbl_qiita).IfNotExists().ToSql()
	db.Exec(query, args...)
	query, args, err = sb.CreateTable(tbl_hatena).IfNotExists().ToSql()
	db.Exec(query, args...)
}

type xmlTime time.Time

func (m *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	buf := ""
	d.DecodeElement(&buf, &start)
	t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", buf)
	if err != nil {
		return err
	}
	*m = xmlTime(t)
	return nil
}

func (m *xmlTime) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if ok {
		*m = xmlTime(t)
	} else {
		return errors.New("type invalid")
	}
	return nil
}

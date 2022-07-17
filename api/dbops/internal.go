package dbops

import (
	"strconv"
)

func InsertSession(sid string, ttl int64, uid int) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := db.Prepare("INSERT INTO sessions (session_id, ttl, user_id) VALUE ($1, $2, $3)")
	CheckErr(err)
	_, err = stmtIns.Exec(sid, ttlstr, uid)
	CheckErr(err)
	defer stmtIns.Close()
	return nil
}

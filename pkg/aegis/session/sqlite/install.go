package sqlite

import "fmt"

func (ss *AegisSqliteSessionStore) Install() error {
	tx, err := ss.connection.Begin()
	if err != nil { return err }
	defer tx.Rollback()
	_, err = tx.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %ssession (
    user_name TEXT,
    value TEXT,
    reg_timestamp INTEGER
)`, ss.config.Session.TablePrefix))
	if err != nil { return err }
	err = tx.Commit()
	if err != nil { return err }
	return nil
}


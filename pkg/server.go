package pkg

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Breaking this endpoint down into two SQL queries because it'll be simpler
// to build out the resulting JSON structure required. This is likely to be
// possible with a single SQL query, but I'm unsure how exactly to structure
// it.
func (c *Conn) ListShips(w http.ResponseWriter, r *http.Request) {
	stmt := `
		SELECT
			s.id,
			s.name,
			s.class,
			s.crew,
			s.image,
			s.value,
			status.value
		FROM spacecraft s
		LEFT JOIN status ON s.status = status.id`
	rows, err := c.Db.Query(stmt)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
	defer rows.Close()

	ships := make([]*Ship, 0)

	for rows.Next() {
		tmp := &Ship{}
		err = rows.Scan(
			&tmp.Id,
			&tmp.Name,
			&tmp.Class,
			&tmp.Crew,
			&tmp.Img,
			&tmp.Value,
			&tmp.Status,
		)
		if err != nil {
			c.Zlog.Warn().Msg(err.Error())
		}

		ships = append(ships, tmp)
	}

	if err = rows.Err(); err != nil {
		c.Zlog.Warn().Msg("Error in database query: ListShips: " + err.Error())
	}

	for _, ship := range ships {
		stmt = "SELECT name, qty FROM armament a WHERE a.ship_id = ?"

		// Passing the ship ID here instead of straight into the statement is
		// best practice to avoid SQL injection
		rows, err := c.Db.Query(stmt, ship.Id)
		if err != nil {
			c.Zlog.Warn().Msg(err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Armament{}
			err = rows.Scan(&tmp.Name, &tmp.Qty)
			if err != nil {
				c.Zlog.Warn().Msg(err.Error())
			}

			ship.Arms = append(ship.Arms, tmp)
		}

	}

	ret, err := json.Marshal(ships)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	_, err = w.Write(ret)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
}

// While very similar to the above ListShips function, this functions just
// differently enough to validate being a separate function.
func (c *Conn) SingleShip(w http.ResponseWriter, r *http.Request) {
	shipID, err := strconv.Atoi(chi.URLParam(r, "shipID"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
	stmt := `
		SELECT
			s.id,
			s.name,
			s.class,
			s.crew,
			s.image,
			s.value,
			status.value
		FROM spacecraft s
		LEFT JOIN status ON s.status = status.id
		WHERE s.id = ?`
	row := c.Db.QueryRow(stmt, shipID)

	var ship Ship
	err = row.Scan(
		&ship.Id,
		&ship.Name,
		&ship.Class,
		&ship.Crew,
		&ship.Img,
		&ship.Value,
		&ship.Status,
	)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	stmt = "SELECT name, qty FROM armament a WHERE a.ship_id = ?"

	// Passing the ship ID here instead of straight into the statement is
	// best practice to avoid SQL injection
	rows, err := c.Db.Query(stmt, ship.Id)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &Armament{}
		err = rows.Scan(&tmp.Name, &tmp.Qty)
		if err != nil {
			c.Zlog.Warn().Msg(err.Error())
		}

		ship.Arms = append(ship.Arms, tmp)
	}

	ret, err := json.Marshal(ship)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	_, err = w.Write(ret)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
}

func (c *Conn) DeleteShip(w http.ResponseWriter, r *http.Request) {
	shipID, err := strconv.Atoi(chi.URLParam(r, "shipID"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
	stmt := "DELETE FROM spacecraft WHERE id = ?"

	_, err = c.Db.Exec(stmt, shipID)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	_, err = w.Write([]byte("Success"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
}

// Required body example:
//
//	curl -X POST -H "Content-Type: application/json" -d '{
//	 "Name": "Spacecraft 1",
//	 "Class": "Class A",
//	 "Crew": 10,
//	 "Img": "image_url",
//	 "Value": 1000000,
//	 "Status": "Operational"
//	}' localhost:7333/ships/new
//
// "Status" must be of a valid value: Operational, Under Repair, Destroyed
// TODO: Handle adding armaments to a spacecraft on insert
func (c *Conn) NewShip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var ship Ship
	err := json.NewDecoder(r.Body).Decode(&ship)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	// This should insert the id of the status value based on the string given
	// in the request body
	stmt := `
		INSERT INTO spacecraft (
			name,
			class,
			crew,
			image,
			value,
			status)
		SELECT ?, ?, ?, ?, ?, id FROM status
		WHERE status.value = ?`
	_, err = c.Db.Exec(stmt, ship.Name, ship.Class, ship.Crew, ship.Img, ship.Value, ship.Status)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	_, err = w.Write([]byte("Success"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
}

// Required body example:
// curl -X PUT -H  "Content-Type: application/json" -d '{"Crew": 50}' localhost:7333/ships/3
// Note that the given column name must be exact and cannot contain anything
// else, otherwise that will cause issues
func (c *Conn) UpdateShip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	shipID, err := strconv.Atoi(chi.URLParam(r, "shipID"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	body := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}

	var key string
	for k := range body {
		key = k
	}

	// TODO: Add field name validation
	stmt := "UPDATE spacecraft SET " + key + " = ? WHERE id = ?"

	// In theory this should only be a length of 1. I'm using this map approach
	// As I don't know what value the user wants to update, so I don't have to
	// Decode the body into a struct. While I could store both the key and value
	// here as struct members, this appears to be the smoother approach to take
	for _, v := range body {
		_, err = c.Db.Exec(stmt, v, shipID)
		if err != nil {
			c.Zlog.Warn().Msg(err.Error())
		}
	}

	_, err = w.Write([]byte("Success"))
	if err != nil {
		c.Zlog.Warn().Msg(err.Error())
	}
}

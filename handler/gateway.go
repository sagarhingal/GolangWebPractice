package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"question/config"
	"question/model"
	"strconv"
)

// internal function to fetch gateway data
func getGateway(db *sql.DB, gatewayId int64) []model.Gateway {
	sql_getGateway := `SELECT * FROM gateway WHERE id = ` + strconv.FormatInt(gatewayId, 10)
	rows, err := db.Query(sql_getGateway)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []model.Gateway
	for rows.Next() {
		gateway := model.Gateway{}
		err2 := rows.Scan(&gateway.ID, &gateway.Name, &gateway.IpAddress)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, gateway)
	}

	return result
}

// CreateGateway: base function to create a new gateway
func CreateGateway(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "POST":

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		reqBody, err := ioutil.ReadAll(req.Body)

		formGate := model.Gateway{}
		json.Unmarshal(reqBody, &formGate)

		// Preparing SQL query
		sql_additem := `
			INSERT INTO gateway(
			Name,
			IpAddress
			) values(?, ?)
		`

		// Inititalising the db object
		db := config.InitDB(config.Dbpath)

		stmt, err := db.Prepare(sql_additem)
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		tempResponse, err2 := stmt.Exec(formGate.Name, formGate.IpAddress)
		if err2 != nil {
			// w.Header().Set("content-type", "application/json; charset=UTF-8")
			// w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("gateway with same name already exists"))
			panic(err2)
		}
		lastID, err := tempResponse.LastInsertId()
		// payloadString := "Newly created ID is " + strconv.FormatInt(lastID, 10)
		payloadObject := model.Gateway{}
		payloadObject.ID = lastID
		payloadObject.Name = formGate.Name
		payloadObject.IpAddress = formGate.IpAddress

		finalGateway := new(bytes.Buffer)
		json.NewEncoder(finalGateway).Encode(payloadObject)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(finalGateway.Bytes())

	case "GET":
		// logic of fetching a gateway
		query := req.URL.Query()
		gatewayId := query.Get("id")

		db := config.InitDB(config.Dbpath)
		var result []model.Gateway
		gid, _ := strconv.ParseInt(gatewayId, 10, 64)
		result = getGateway(db, gid)

		//Preparing the final JSON payload
		finalGateway := new(bytes.Buffer)
		json.NewEncoder(finalGateway).Encode(result[0])

		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(finalGateway.Bytes())
	}
}

// CreateRoute: base function to create a new route
func CreateRoute(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		reqBody, err := ioutil.ReadAll(req.Body)
		formRoute := model.CustomRoute{}
		json.Unmarshal(reqBody, &formRoute)
		db := config.InitDB(config.Dbpath)
		// Step 1 - Fetching the required gateway data
		//-------------------------------------------------
		var result []model.Gateway
		result = getGateway(db, formRoute.GatewayId)
		//-----------------------------------------------------
		// Step 2 - Post the new route
		// ----------------------------------------------------
		// Preparing new payload
		// Preparing SQL query
		sql_additem := `
			  INSERT OR REPLACE INTO route(
			  Prefix,
			  GatewayId
			  ) values(?, ?)
		  `

		stmt, err := db.Prepare(sql_additem)
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		tempResponse, err2 := stmt.Exec(formRoute.Prefix, formRoute.GatewayId)
		lastID, err := tempResponse.LastInsertId()
		payloadObject := model.CustomRouteResponse{}
		payloadObject.ID = lastID
		payloadObject.Prefix = formRoute.Prefix
		payloadObject.Gateway = result[0]

		if err2 != nil {
			panic(err2)
		}

		finalRoute := new(bytes.Buffer)
		json.NewEncoder(finalRoute).Encode(payloadObject)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(finalRoute.Bytes())
		//----------------------------------------------------------------
	case "GET":
		query := req.URL.Query()
		routeId := query.Get("id")

		// Step 1 - Get the route object
		// --------------------------------------------------------
		sql_getRoute := `SELECT * FROM route WHERE id = ` + routeId
		db := config.InitDB(config.Dbpath)
		rows, err := db.Query(sql_getRoute)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var result []model.Route
		for rows.Next() {
			item := model.Route{}
			err2 := rows.Scan(&item.ID, &item.Prefix, &item.GatewayId)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, item)
		}
		// ----------------------------------------------------------
		// Step 2 - Get the Gateway object
		// ----------------------------------------------------------
		var gatewayResult []model.Gateway
		gatewayResult = getGateway(db, result[0].GatewayId)
		// ----------------------------------------------------------
		// Step 3 - Prepare the final payload
		payloadObject := model.CustomRouteResponse{}
		pid, _ := strconv.ParseInt(result[0].ID, 10, 64)
		payloadObject.ID = pid
		payloadObject.Prefix = result[0].Prefix
		payloadObject.Gateway = gatewayResult[0]
		finalRoute := new(bytes.Buffer)
		json.NewEncoder(finalRoute).Encode(payloadObject)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(finalRoute.Bytes())
	}
}

// SearchRoute: Base function to seach a route given a number
func SearchRoute(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	baseNumber := query.Get("number")
	var prefix string

	// Step 1 -  Parse the number and find out the prefix
	areaCodes := [3]string{"123", "1234", "9194"}
	basePrefix := baseNumber[:4]
	for i := 1; i < 2; i++ {
		if basePrefix == areaCodes[i] {
			prefix = areaCodes[i]
		}
	}

	if prefix == "" {
		if basePrefix == areaCodes[0] {
			prefix = areaCodes[0]
		}
	}
	// -----------------------------------------------------
	// Step 2 - Search the route using prefix
	sql_getRoute := `SELECT * FROM route WHERE prefix LIKE '` + prefix + `%'`
	db := config.InitDB(config.Dbpath)
	rows, err := db.Query(sql_getRoute)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []model.Route
	for rows.Next() {
		item := model.Route{}
		err2 := rows.Scan(&item.ID, &item.Prefix, &item.GatewayId)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	// ----------------------------------------------------------
	// Step 3 - Get the Gateway object
	// ----------------------------------------------------------
	var gatewayResult []model.Gateway
	gatewayResult = getGateway(db, result[0].GatewayId)
	// ----------------------------------------------------------
	// Step 4 - Prepare the final payload
	payloadObject := model.CustomRouteResponse{}
	pid, _ := strconv.ParseInt(result[0].ID, 10, 64)
	payloadObject.ID = pid
	payloadObject.Prefix = result[0].Prefix
	payloadObject.Gateway = gatewayResult[0]
	finalRoute := new(bytes.Buffer)
	json.NewEncoder(finalRoute).Encode(payloadObject)
	w.Header().Set("content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(finalRoute.Bytes())

}

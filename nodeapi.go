package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	common "zetanet.io/common"
	util "zetanet.io/utils"
)

//NodeAPI is a REST api for nodes
type NodeAPI struct {
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

//Init initialize REST API
func (na *NodeAPI) Init() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/nodes", getNodes).Methods("GET")
	router.HandleFunc("/contents", getContents).Methods("GET")

	fmt.Println("API Started..")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	if db, err := common.NewDb(util.LoadConfigByKey("DB_NODES")); err == nil {
		iterator := db.NewIterator(nil, nil)
		var node common.Node
		var nodes []common.Node
		for iterator.Next() {
			if err := json.Unmarshal(iterator.Value(), &node); err == nil {
				nodes = append(nodes, node)
			} else {
				fmt.Println("Get:" + err.Error())
			}
		}
		db.Close()
		json.NewEncoder(w).Encode(nodes)
	} else {
		fmt.Println("OpenFile:" + err.Error())
	}
}

func getContents(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	if db, err := common.NewDb(util.LoadConfigByKey("DB_CONTENTS")); err == nil {
		iterator := db.NewIterator(nil, nil)
		var desc common.Desc
		var descs []common.Desc
		for iterator.Next() {
			if err := json.Unmarshal(iterator.Value(), &desc); err == nil {
				descs = append(descs, desc)
			} else {
				fmt.Println("Get:" + err.Error())
			}
		}
		db.Close()
		json.NewEncoder(w).Encode(descs)
	} else {
		fmt.Println("OpenFile:" + err.Error())
	}
}

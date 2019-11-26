package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	util "zetanet.io/utils"
	common "zetanet.io/common"
	p2p "zetanet.io/p2p"
)

//DiscoveryNode is the first node where the newly created storage nodes listens
type DiscoveryNode struct {
	config common.Config
}

func (dn *DiscoveryNode) listen() {
	config, err := common.InitConfig("./config/" + util.LoadEnv() + ".yml")
	
	dn.config = config

	if err != nil {
		fmt.Println(err)
	}

	if l, err := net.Listen(dn.config.Type, dn.config.Host+":"+dn.config.Port); err == nil {
		defer l.Close()

		fmt.Println("Listening on " + dn.config.Type + " endpoint: " + dn.config.Host + ":" + dn.config.Port)

		for {
			if conn, err := l.Accept(); err == nil {
				go dn.handleRequest(conn)
			} else {
				fmt.Println("Error accepting:", err.Error())
				os.Exit(1)
			}
		}
	}
}

func (dn *DiscoveryNode) handleRequest(conn net.Conn) {
	data, _, _ := bufio.NewReader(conn).ReadLine()

	cmd, body := extractRequest(data)

	fmt.Println("REQUEST:", string(data))

	switch string(cmd) {
		
	case string(common.Command.Reg):

		//register node by adding it to leveldb
		register(conn, body)
		dn.GetNodes(conn)
		break

	case string(common.Command.Get):

		//get the content based on content-address hash
		get(conn, body)
		break

	case string(common.Command.Add):

		//add the content description and hash to leveldb
		add(conn, body)
		break
	}

	

}

func  register(conn net.Conn, data []byte){
	//validate json
	var node common.Node
	if err := json.Unmarshal([]byte(data), &node); err == nil {

		//save json to leveldb
		if db, err := common.NewDb(util.LoadConfigByKey("DB_NODES")); err == nil {
			defer db.Close()

			if value, err := json.Marshal(node); err == nil {

				//add bytes converted node information to leveldb
				db.Put([]byte(node.Host+":"+node.Port), value, nil)

				//log newly added node
				if data, err := db.Get([]byte(node.Host+":"+node.Port), nil); err == nil {
					fmt.Println("REGISTERED: " + string(data))
				} else {
					fmt.Println("Error:" + err.Error())

				}
			} else {
				fmt.Println("Error:" + err.Error())
			}
		} else {
			fmt.Println("Error:" + err.Error())
		}
	} else {
		fmt.Println(err)
	}
}
func add(conn net.Conn, data []byte){

	//data is the self-describing info of the content

	//validate json
	var desc common.Desc
	if err := json.Unmarshal([]byte(data), &desc); err == nil {
		//save json to leveldb
		if db, err := common.NewDb(util.LoadConfigByKey("DB_CONTENTS")); err == nil {
			defer db.Close()

			if err := db.Put([]byte(desc.Hs), data, nil); err != nil {
				fmt.Println("Put:" + err.Error())
			} else {
				response, _ := json.Marshal(p2p.CreateResponse(200, "Success"))
				conn.Write(response)
				fmt.Println("Success!")
			}
		} else {
			fmt.Println("OpenFile:" + err.Error())
		}
	} else {
		fmt.Println(err)
	}
}


func get(conn net.Conn, hash []byte){
	if db, err := common.NewDb(util.LoadConfigByKey("DB_CONTENTS")); err == nil {
		
		data, _ := db.Get(hash, nil)
		var desc common.Desc

		if err := json.Unmarshal(data, &desc); err == nil {

			//TODO: return a formatted data
			fmt.Println("RESPONDED:", string(data))
			conn.Write(append(data, '\n'))
		} else {
			fmt.Println("Error:" + err.Error())
		}
		
		db.Close()
		//json.NewEncoder().Encode(descs)
	} else {
		fmt.Println("OpenFile:" + err.Error())
	}
}



//extractRequest returns command,body
func extractRequest(data []byte) ([]byte, []byte){
	return data[0:common.HeaderCmdSize], data[common.HeaderCmdSize:len(data)]
}


//GetNodes returns all the registered nodes
func (dn *DiscoveryNode) GetNodes(conn net.Conn) {

	if db, err := common.NewDb(util.LoadConfigByKey("DB_NODES")); err == nil {
		iterator := db.NewIterator(nil, nil)

		for iterator.Next() {
			conn.Write(append(iterator.Value(), '\n'))
		}

		iterator.Release()
		db.Close()
		conn.Close()
	} else {
		fmt.Println("OpenFile:" + err.Error())
	}
}

func (dn *DiscoveryNode) stop() {
	fmt.Println("Stopping Discovery Node")
	os.Exit(0)
}

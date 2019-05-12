package settings

import (
    "../utils"
)

func DeleteLocalMessages(slice []string) {
	var node_address map[string]string
	if User.ModeF2F {
		node_address = Node.Address.F2F
	} else {
		node_address = Node.Address.P2P
	}
    Mutex.Lock()
    for _, user := range slice {
        if _, ok := node_address[user]; ok {
            _, err := DataBase.Exec("DELETE FROM Local" + user)
            utils.CheckError(err)
        }
    }
    Mutex.Unlock()
}

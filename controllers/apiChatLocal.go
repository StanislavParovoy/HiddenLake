package controllers

import (
    "time"
    "strings"
    "net/http"
    "encoding/json"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func apiChatLocal(w http.ResponseWriter, r *http.Request) {
    var result = strings.TrimPrefix(r.URL.Path, "/api/chat/")
    w.Header().Set("Content-Type", "application/json")

    var slice = strings.Split(result, "/")

    switch slice[1] {
        case "update": 
            updateChatLocal(w, slice[0])
        default: 
            json.NewEncoder(w).Encode(models.PackageHTTP{Head: settings.HEAD_WARNING})
    } 
}

func updateChatLocal(w http.ResponseWriter, user string) {
    var (
        mode = settings.CurrentMode()
        node_address = settings.CurrentNodeAddress()
    )

    if _, ok := node_address[user]; !ok {
        json.NewEncoder(w).Encode(models.PackageHTTP{Exists:false})
        return
    }

    timeout := make(chan bool)
    go func() {
        time.Sleep(time.Second * 30)
        timeout <- true
    }()

    select {
        case <-timeout:
            json.NewEncoder(w).Encode(models.PackageHTTP{Exists:false})
            return

        case <-settings.Messages.NewDataExistLocal[user]:
            rows, err := settings.DataBase.Query(
                "SELECT Body FROM Local" + user + " WHERE Mode = $1 ORDER BY Id DESC LIMIT $2",
                mode,
                settings.Messages.CurrentIdLocal[user],
            )
            utils.CheckError(err)

            var messages, message string
            for rows.Next() {
                rows.Scan(&message)
                messages += crypto.Decrypt(settings.User.Password, message)
            }

            var data = models.PackageHTTP {
                Exists: true,
                Head: settings.HEAD_MESSAGE,
                Body: messages,
            }

            settings.Mutex.Lock()
            settings.Messages.CurrentIdLocal[user] = 0
            settings.Mutex.Unlock()

            json.NewEncoder(w).Encode(data)
            return
    }
}

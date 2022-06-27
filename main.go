package main

import "fmt"
import "net/http"
import "encoding/json"
import "strings"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"



func connect() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/svgodb")
    if err != nil {
        return nil, err
    }

    return db, nil
}

func handleSave(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        payload := struct {
            Title       string `json:"title"`
            Content     string `json:"content"`
            Category    string `json:"category"`
            Status      string `json:"status"`
        }{}
        if err := decoder.Decode(&payload); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        message := fmt.Sprintf(
            "hello, here is my title %s. I'm %s in category %s with %s", 
            payload.Title, 
            payload.Content, 
            payload.Category,
            payload.Status,
        )
        w.Write([]byte(message))
        return
    }

    http.Error(w, "Only accept POST request", http.StatusBadRequest)
}

func ActionIndex(w http.ResponseWriter, r *http.Request) {
    data := [] struct {
        Name string
        Age  int
    } {
        { "Richard Grayson", 24 },
        { "Jason Todd", 23 },
        { "Tim Drake", 22 },
        { "Damian Wayne", 21 },
    }

    jsonInBytes, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonInBytes)
}

func insertDB(title string, content string, category string, status string)(string){
    db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return  ""
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?);", title, content, category, status)
    if err != nil {
        fmt.Println(err.Error())
        return ""
    }

    return "Berhasil menyimpan data"
}

func selectDB(sqlString string) (string) {
    db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return ""
    }
    defer db.Close()
    
    rows, err := db.Query(sqlString)
    if err != nil {
        return ""
    }
    defer rows.Close()
    columns, err := rows.Columns()
    if err != nil {
        return ""
    }
    count := len(columns)
    tableData := make([]map[string]interface{}, 0)
    values := make([]interface{}, count)
    valuePtrs := make([]interface{}, count)
    for rows.Next() {
        for i := 0; i < count; i++ {
          valuePtrs[i] = &values[i]
        }
        rows.Scan(valuePtrs...)
        entry := make(map[string]interface{})
        for i, col := range columns {
            var v interface{}
            val := values[i]
            b, ok := val.([]byte)
            if ok {
                v = string(b)
            } else {
                v = val
            }
            entry[col] = v
        }
        tableData = append(tableData, entry)
    }
    jsonData, err := json.Marshal(tableData)
    if err != nil {
        return ""
    }
    // fmt.Println(string(jsonData))
    return string(jsonData)
}

func updateDB(title string, content string, category string, status string, id string)(string){
    db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return  ""
    }
    defer db.Close()

    _, err = db.Exec("UPDATE posts SET title = ?, content = ?, category = ?, status = ?  WHERE id=?;", title, content, category, status, id)
    if err != nil {
        fmt.Println(err.Error())
        return ""
    }

    return "Berhasil merubah data"
}

func deleteDB(id string)(string){
    db, err := connect()
    if err != nil {
        fmt.Println(err.Error())
        return  ""
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM posts WHERE id=?;", id)
    if err != nil {
        fmt.Println(err.Error())
        return ""
    }

    return "Berhasil menghapus data"
}

func writeRes(code string, message string, data string)([]byte){
    resp := make(map[string]string)
    resp["code"] = code
    resp["message"] = message
    resp["data"] = data
    jsonResp, err := json.Marshal(resp)
    if err != nil {}
    return jsonResp
}

func main(){    
    http.HandleFunc("/article/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            decoder := json.NewDecoder(r.Body)
            payload := struct {
                Title       string `json:"title"`
                Content     string `json:"content"`
                Category    string `json:"category"`
                Status      string `json:"status"`
            }{}
            if err := decoder.Decode(&payload); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if (len(payload.Title) < 20){
                    w.Write(writeRes("400","Title minimal 20 karakter",""))
            } else {
                if (len(payload.Content) < 200){
                    w.Write(writeRes("400","Content minimal 200 karakter",""))
                } else {
                    if (len(payload.Category) < 3){
                        w.Write(writeRes("400","Category minimal 3 karakter",""))
                    } else {
                        switch payload.Status {
                        case "publish", "draft", "trash":
                            messageDB := insertDB(payload.Title, payload.Content, payload.Category, payload.Status)
                            w.Write(writeRes("200",messageDB,""))
                        default:
                            {
                                w.Write(writeRes("400","Status invalid",""))
                            }
                        }
                    }
                }
            }
        case "GET":
                sUrl := strings.Split(r.URL.Path[1:], "/")
                var query string
                
                lenUrl := len(sUrl)
                if (lenUrl >= 3){
                    query =  "SELECT * FROM posts LIMIT "+sUrl[1]+","+sUrl[2]
                }   else {
                    query =  "SELECT * FROM posts WHERE id="+sUrl[1]
                }
                w.Write(writeRes("200","Berhasil mendapatkan data",selectDB(query)))
        case "PUT":
            sUrl := strings.Split(r.URL.Path[1:], "/")
            decoder := json.NewDecoder(r.Body)
            payload := struct {
                Title       string `json:"title"`
                Content     string `json:"content"`
                Category    string `json:"category"`
                Status      string `json:"status"`
            }{}
            if err := decoder.Decode(&payload); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if (len(payload.Title) < 20){
                w.Write(writeRes("400","Title minimal 20 karakter",""))
            } else {
                if (len(payload.Content) < 200){
                    w.Write(writeRes("400","Content minimal 200 karakter",""))
                } else {
                    if (len(payload.Category) < 3){
                        w.Write(writeRes("400","Category minimal 3 karakter",""))
                    } else {
                        switch payload.Status {
                        case "publish", "draft", "trash":
                            messageDB := updateDB(payload.Title, payload.Content, payload.Category, payload.Status, sUrl[1])
                            w.Write(writeRes("200",messageDB,""))
                        default:
                            {
                                w.Write(writeRes("400","Status invalid",""))
                            }
                        }
                    }
                }
            }
        case "DELETE":
            sUrl := strings.Split(r.URL.Path[1:], "/")

            messageDB := deleteDB(sUrl[1])
            w.Write(writeRes("200",messageDB,""))     
        default:
            http.Error(w, "", http.StatusBadRequest)
        }
    })
    http.ListenAndServe(":9000", nil)
}
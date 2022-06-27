package main

import "fmt"
import "net/http"
import _ "html/template"
import "encoding/json"
import _ "reflect"
import "strings"
import _ "path"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

// type posts struct{
// 	id string
// 	title string
// 	content string
// 	category string
// 	status string
// }

func connect() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/svgodb")
    if err != nil {
        return nil, err
    }

    return db, nil
}

// func sqlQuery(w http.ResponseWriter, r *http.Request) {
//     db, err := connect()
//     if err != nil {
//         fmt.Println(err.Error())
//         return
//     }
//     defer db.Close()

//     // var age = 27
//     rows, err := db.Query("select id, title, content, category, status from posts")
//     if err != nil {
//         fmt.Println(err.Error())
//         return
//     }
//     defer rows.Close()

//     var result []posts

//     for rows.Next() {
//         var each = posts{}
//         var err = rows.Scan(&each.id, &each.title, &each.content, &each.category, &each.status)

//         if err != nil {
//             fmt.Println(err.Error())
//             return
//         }

//         result = append(result, each)
//     }

//     if err = rows.Err(); err != nil {
//         fmt.Println(err.Error())
//         return
//     }

//     for _, each := range result {
//         fmt.Println(each.title)
//     }
// }

// func routeIndexGet(w http.ResponseWriter, r *http.Request) {
//     if r.Method == "GET" {
//         var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
//         var err = tmpl.Execute(w, nil)

//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//         return
//     }

//     http.Error(w, "", http.StatusBadRequest)
// }

// func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
//     if r.Method == "POST" {
//         var tmpl = template.Must(template.New("form").ParseFiles("view.html"))

//         if err := r.ParseForm(); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }

//         var name = r.FormValue("name")
//         var message = r.Form.Get("message")

//         var data = map[string]string{"name": name, "message": message}

//         if err := tmpl.Execute(w, data); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//         return
//     }

//     http.Error(w, "", http.StatusBadRequest)
// }

// func handleIndex(w http.ResponseWriter, r *http.Request) {
//     tmpl := template.Must(template.ParseFiles("view.html"))
//     if err := tmpl.Execute(w, nil); err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//     }
// }

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

func main(){

    // http.HandleFunc("/", handleIndex)
    // http.HandleFunc("/save", handleSave)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
            messageDB := insertDB(payload.Title, payload.Content, payload.Category, payload.Status)
            resp := make(map[string]string)
            resp["code"] = "200"
            resp["message"] = messageDB
            resp["data"] = ""
            jsonResp, err := json.Marshal(resp)
            if err != nil {}
            w.Write(jsonResp)
        default:
            http.Error(w, "", http.StatusBadRequest)
        }
    })
    
    http.HandleFunc("/article/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            decoder := json.NewDecoder(r.Body)
            payload := struct {
                Title       string `json:"title" validate:"required,min=20"`
                Content     string `json:"content" validate:"required,min=200"`
                Category    string `json:"category" validate:"required,min=3"`
                Status      string `json:"status" validate:"required"`
            }{}
            if err := decoder.Decode(&payload); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if (payload.Status == "publish" && payload.Status != "draft" && payload.Status != "trash"){
                messageDB := insertDB(payload.Title, payload.Content, payload.Category, payload.Status)
                resp := make(map[string]string)
                resp["code"] = "200"
                resp["message"] = messageDB
                resp["data"] = ""
                jsonResp, err := json.Marshal(resp)
                if err != nil {}
                w.Write(jsonResp)
            } else {
                resp := make(map[string]string)
                resp["code"] = "400"
                resp["message"] = "Status invalid"
                resp["data"] = ""
                jsonResp, err := json.Marshal(resp)
                if err != nil {}
                w.Write(jsonResp)
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
                
                resp := make(map[string]string)
                resp["code"] = "200"
                resp["message"] = "Berhasil mendapatkan data"
                resp["data"] = selectDB(query)
                jsonResp, err := json.Marshal(resp)
                if err != nil {}
                w.Write(jsonResp)
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

            messageDB := updateDB(payload.Title, payload.Content, payload.Category, payload.Status, sUrl[1])
            resp := make(map[string]string)
            resp["code"] = "200"
            resp["message"] = messageDB
            resp["data"] = ""
            jsonResp, err := json.Marshal(resp)
            if err != nil {}
            w.Write(jsonResp)
        case "DELETE":
            sUrl := strings.Split(r.URL.Path[1:], "/")

            messageDB := deleteDB(sUrl[1])
            resp := make(map[string]string)
            resp["code"] = "200"
            resp["message"] = messageDB
            resp["data"] = ""
            jsonResp, err := json.Marshal(resp)
            if err != nil {}
            w.Write(jsonResp)        
        default:
            http.Error(w, "", http.StatusBadRequest)
        }
    })
    http.ListenAndServe(":9000", nil)
}

// func HelloServer(w http.ResponseWriter, r *http.Request) {
//     s := strings.Split(r.URL.Path[1:], "/")
//     fmt.Fprintf(w, "Hello, LIMIT %s OFFSET %s !", s[1], s[2])
// }
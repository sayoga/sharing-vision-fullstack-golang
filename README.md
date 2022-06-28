# sharing_vision_be_golang


Back End (Golang):
- Untuk semua endpoint sudah di masukan masing-masing route, dan port localhost yang digunakan adalah 9000. Semua route baik BE maupun FE ada di file main.go
  1)  Membuat article baru
      POST
      http://localhost:9000/article/
     
  2)  Menampilkan seluruh article di database dengan paging pada parameter limit & offset
      GET
      http://localhost:9000/article/1/10
      
  3)  Menampilkan article dengan id yang di-request
      GET
      http://localhost:9000/article/1
      
  4)  Merubah data article dengan id yang di-request
      PUT
      http://localhost:9000/article/1
      
  5)  Menghapus data article dengan id yang di request
      DELETE
      http://localhost:9000/article/1
      
      
Front End (Bootstrap & JQuery)
- Untuk file FE-nya ada di folder public, dan route sama seperti BE ada di file main.go

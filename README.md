# sharing_vision_be_golang


Back End (Golang):
- Untuk semua endpoint yang dibutuhkan sudah dimasukan masing-masing route, dan port localhost yang digunakan adalah 9000. Semua route baik BE maupun FE ada di file main.go
  1)  Membuat article baru <br>
      POST <br>
      http://localhost:9000/article/
     
  2)  Menampilkan seluruh article di database dengan paging pada parameter limit & offset <br>
      GET <br>
      http://localhost:9000/article/1/10
      
  3)  Menampilkan article dengan id yang di-request <br>
      GET <br>
      http://localhost:9000/article/1
      
  4)  Merubah data article dengan id yang di-request <br>
      PUT <br>
      http://localhost:9000/article/1
      
  5)  Menghapus data article dengan id yang di request <br>
      DELETE <br>
      http://localhost:9000/article/1
      
      
Front End (Bootstrap & JQuery)
- Untuk file FE-nya ada di folder public, dan route sama seperti BE ada di file main.go

SQL
- Terlampir juga mengenai file sql sesuai dengan kebutuhan tabel yang dimaksud

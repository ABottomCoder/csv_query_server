# csv_query_server

csv_query_server is a server application built with Go that provides two APIs for querying and modifying CSV files.

## Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/ABottomCoder/csv_query_server.git
   
2. Change to the project directory:

    ```shell
   cd csv_query_server

3. Build the server

    ```shell
   go build -o csv_query_server.exe

4. Run the server

    ```shell
   ./csv_query_server

## Usage
The CSV Query Server provides the following 2 APIs:

1. Query API: Allows you to query the CSV file based on specific criteria. The API accepts query parameters and returns the matching records from the CSV file.
   
   Method: GET

   column: The name of the column to query, default column is [C1 C2 C3]
   
   value: The value to search for in the specified column.

   Example Request: http://127.0.0.1:9527/?query=C3%3D%3Daa  (query C3==aa)

   example Response:
   ```json
   {
    "result": [
        [
            "a1",
            "ah",
            "aa"
        ]
    ],
    "msg": ""
    }

2. Modify API: Allows you to modify the CSV file with INSERT, DELETE, UPDATE. The API accepts query parameters and returns modify result.

   Method: GET

   column: The name of the column to query, default column is [C1 C2 C3]

   value: The value to search for in the specified column.


   - Example INSERT Request: http://127.0.0.1:7259/?job=INSERT%20b1%2Cb2%2Cb3  (insert [b1 b2 b3])
   - This will insert a row [b1 b2 b3] to csv if it not in csv

      example Response:
      ```json
      {
      "msg": "Modification successful"
      }
     
      // run second time, will get this response, and terminal will output "Execute modify fail, err: data already exist"
     {
     "msg": "Execute modify fail"
     }

   - Example DELETE Request: http://127.0.0.1:7259/?job=DELETE%20b1%2Cb2  (delete [b1 b2])
   - this will delete all rows where C1==b1 and C2==b2
       
        example Response:
        ```json
        {
        "msg": "Modification successful"
        }

   - Example UPDATE Request: http://127.0.0.1:7259/?job=UPDATE%20a1%2Ca3%2CC3%2Ctt  (update [a1 a3 C3->tk])
   - this will update all rows where C1==a1 and C2==a3, it will update C3 to tk

      example Response:
      ```json
     {
     "msg": "Modification successful"
     }
   
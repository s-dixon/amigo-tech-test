# Amigo Tech Test

**Disclaimer: I haven't used Go before so please excuse any bad coding practices :)**

With the requirements for adding and retrieving a single message defined within the specification of the test, I decided to make the following elaborations:

  - Tracking "users" by means of recording their IP address (if present within the request)
  - Recording when messages are created
  - Functionality to delete messages
  - View messages: enabling messages to be retrieved in "pages" (supporting pagination); filter by message content; and filter by IP address
  - Utilising PostgreSQL for data persistence - for consistency; given that it's the database technology already adopted by Amigo

# Running Locally

Setup the database - I've included a docker file for convenience:
``` sh
$ cd amigo-tech-test/docker
$ ./docker-compose.exe up
```
  
Build the app using Go
``` sh
$ cd amigo-tech-test
$ go get
$ go build
```

Run the web service
``` sh
$ ./amigo-tech-test.exe
```
  
# API
### **Create Message**

* **URL**

  /messages/

* **Method:**

  `POST`

* **Data Params**

  Request body: message value

* **Success Response:**

  * **Code:** 200
    **Content:** `{ "id" : 12 }`
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST
    **Content:** `{ "error" : "Invalid message payload" }`

  OR

  * **Code:** 500 INTERNAL SERVER ERROR
    **Content:** `{ "error" : "<error_description>" }`

* **Sample Call:**

  ```
     curl $domain/messages/ -d 'my test message to store'
  ```
### **Delete Message**

* **URL**

  /messages/:id

* **Method:**

  `DELETE`

*  **URL Params**

   **Required:**
 
   `id=[integer]`

* **Success Response:**

  * **Code:** 200
    **Content:** `{"result":"success"}`
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST
    **Content:** `{ error : "Invalid message ID" }`

  OR

  * **Code:** 500 INTERNAL SERVER ERROR
    **Content:** `{ error : "<error_description>" }`

* **Sample Call:**

  ```
     curl $domain/messages/12 -X DELETE
  ```
### **Get Message**

* **URL**

  /messages/:id

* **Method:**

  `GET`

*  **URL Params**

   **Required:**
 
   `id=[integer]`

* **Success Response:**

  * **Code:** 200
    **Content:** `<message>`
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST
    **Content:** `{ error : "Invalid message ID" }`

  OR
  
  * **Code:** 404 NOT FOUND
    **Content:** `{ error : "Message not found" }`
 
  OR

  * **Code:** 500 INTERNAL SERVER ERROR
    **Content:** `{ error : "<error_description>" }`

* **Sample Call:**

  ```
     curl $domain/messages/12
  ```

### **Get Messages**

Retrieve a single page of messages using default/provided parameters

* **URL**

  /messages/

* **Method:**

  `GET`

*  **Query String Params**

   **Optional:**
 
   `offset=[integer]`  (offset of page results - default: 0)
   `limit=[integer]` (the maximum number of results - default/maximum: 20)
   `message=[string]` (filter by messages containing supplied value)
   `ip=[string]` (filter messages by matching supplied IP address/pattern)
    
* **Success Response:**

  * **Code:** 200
    **Content:** 
    ``` json
    {
        "offset":0,
        "limit":20,
        "total_count":3,
        "results":[
            {"id":1,"value":"message1","ip_address":"::1","date_created":"2017-06-25T14:11:57.663843Z"},
            {"id":2,"value":"message2","ip_address":"::1","date_created":"2017-06-25T14:22:12.296925Z"},
            {"id":3,"value":"message3","ip_address":"::1","date_created":"2017-06-25T14:30:51.457346Z"}
        ]
    }
    ```
 
* **Error Response:**

  * **Code:** 500 INTERNAL SERVER ERROR
    **Content:** `{ error : "<error_description>" }`

* **Sample Call:**

  ```
     curl $domain/messages/?ip=192.168.200.201&message=foo&offset=20&limit=5
  ```

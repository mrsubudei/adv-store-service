**Overview**
----
This project is implimintation of RESTful api which allow to store, get, update and delete adverts.

**Content**
----  

- [Status codes](#status-codes)
- [Create advert](#create-advert)
- [Get advert](#get-advert)
- [Get all adverts](#get-all-adverts)
- [Update advert](#update-advert)
- [Delete advert](#delete-advert)  

- [Usage](#usage)
  
**Status codes**
----

The following table gives an overview of how the API functions generally behave.

| Request type | Description |
| ------------ | ----------- |
| `GET`   | Access one or more adverts and return the result as JSON. |
| `POST`  | Return `201 Created` if the resource is successfully created and the ID of created advert returned. |
| `GET` / `PUT` / `DELETE` | Return `200 OK` if the resource is accessed or modified or deleted successfully. |  

The following table shows the possible return codes for API requests.

| Return values | Description |
| ------------- | ----------- |
| `200 OK` | The `GET`, `PUT` or `DELETE` request was successful. |
| `201 Created` | The `POST` request was successful and the ID of created advert returned. |
| `400 Bad Request` | A required attribute of the API request is missing. |
| `404 Not Found` | A resource could not be accessed, e.g., an ID for a resource could not be found. |
| `405 Method Not Allowed` | The request is not supported. |
| `409 Conflict` | A conflicting advert's name already exists |
| `500 Server Error` | While handling the request something went wrong server-side. |  

**Create advert**
----
  Return ID of created advert.  
  Adverts should have unique names. Name length limit is 200 symbols. Description length limit is 1000 symbols.  
  Number of photo urls links minimum 1, maximum 3. First url will become main url. Required fields should not be empty.

* **URL**

  /v1/adverts

* **Method:**

  `POST`
  
*  **URL Params**

    None

* **Data Params**

```json
{
  "name": "some name",
  "description": "some description",
  "price": 150,
  "photo_urls": [
    "http://files.com/12",
    "http://files.com/13"
  ]
}
```

* **Success Response:**

  * **Code:** 201 <br />
    **Content:** 

```json
{
 "data": [
        {
            "id": 5
        }
    ]
}
```
 
* **Error Response:**

  * *Request has wrong json fields*
    **Code:** 400 BAD REQUEST <br />
    **Content:** 
```json
{
    "error": "json format is not correct"
}
```

  OR with details

```json
{
    "error": "Request Entity Too Large",
    "detail": "'name:' field's length exceeded"
}
```  

  * *There is already advert with same name*
    **Code:** 409 STATUS CONFLICT <br />
    **Content:** 
```json
{
    "error": "item with name 'some name' already exists"
}
```

**Get advert**
----
  Return JSON with advert's information.

* **URL**

  /v1/adverts/{id}

* **Method:**

  `GET`
  
*  **URL Params**

   **Optional:**
 
   `fields=[true]`

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 

```json
{
  "data": [
    {
      "name": "advert1",
      "price": 100,
      "main_photo_url": "http:fileserver.com/125"
    }
  ]
}
```

 With 'fields' param:

 ```json
{
  "data": [
    {
      "id": 14,
      "name": "advert1",
      "description": "some description",
      "price": 100,
      "main_photo_url": "http:fileserver.com/125",
      "photo_urls": [
        "http:fileserver.com/125",
        "http:fileserver.com/126"
        ]
    }
  ]
}
```

* **Error Response:**

  * *There is no item with that ID*
    **Code:** 404 NOT FOUND <br />
    **Content:** 
```json
{
    "error": "no content found with id: 12345"
}
```

**Get all adverts**
----
  Return JSON with list of all adverts with metadata of max page. By default, 10 adverts will be given per page.  
  It is possible to sort adverts by 'created date' or 'price'.

* **URL**

  /v1/adverts

* **Method:**

  `GET`
  
*  **URL Params**

   **Optional:**
 
   `limit=[integer]`  
   `offset=[integer]`  
   `sort_by=[created_at] or [price]`  
   `order_by=[asc] or [desc]`

* **Data Params**

  None

* **Success Response:**

* **Code:** 200 OK <br />
  **Content:** 

```json
{
  "meta_data": {
    "max_page": 11
    },
    "data": [
      {
        "name": "advert1",
        "price": 30,
        "main_photo_url": "http:fileserver.com/125"
      },
      {
        "name": "advert2",
        "price": 40,
        "main_photo_url": "http:fileserver.com/126"
      },
      {
        "name": "advert3",
        "price": 50,
        "main_photo_url": "http:fileserver.com/127"
      }
   ]
}
```

* **Error Response:**

  * *There is no items*
  * **Code:** 200 OK <br />
  **Content:** 

```json
{}
```

**Update advert**
----
  Return status code and empty JSON. Only fields are passed will be changed. If in 'photo_urls' field at least one url is sent, all urls of this advert will be deleted and replaced with sent urls.

* **URL**

  /v1/adverts/{id}

* **Method:**

  `PUT`
  
*  **URL Params**

   None

* **Data Params**

```json
{
  "description": "new description"
}
```

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 

```json
{}
```

* **Error Response:**

  * *There is no item with that ID*
    **Code:** 404 NOT FOUND <br />
    **Content:** 
```json
{
    "error": "no content found with id: 12345"
}
```

  * *There is item with given name already exists*
    **Code:** 409 STATUC CONFLICT <br />
    **Content:** 
```json
{
    "error": "item with name 'car' already exists"
}
```

**Delete advert**
----
  Return status code and empty JSON.

* **URL**

  /v1/adverts/{id}

* **Method:**

  `DELETE`
  
*  **URL Params**

   None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** 

```json
{}
```

* **Error Response:**

  * *There is no item with that ID*
    **Code:** 404 NOT FOUND <br />
    **Content:** 
```json
{
    "error": "no content found with id: 12345"
}
```

**Usage**
----
Run app
```
make go
```  
Run in docker container
```
make start
```  
Stop container
```
make stop
```  
Begin tests
```
make test
```  
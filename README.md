### Run
#### PostgreSQL build and run
```shell
  make dc_db_build
```

#### PostgreSQL run
```shell
  make dc_db
```

#### In-memory build and run
```shell
  make dc_inmem_build
```

#### In-memory run
```shell
  make dc_inmem_run
```

### Tests
#### Tests with coverage .html file
```shell
    make test
```

### Endpoints
#### Default endpoint

*0.0.0.0:8000*


### Methods
#### **POST** /api/http
##### Accepts
```json
{
  "url": "https://ozon.ru"
}
```
##### Returns
```json
{
  "status": "OK",
  "alias": "http://0.0.0.0:8000/Y2Wv0YnlFe"
}
```
or
```json
{
  "status": "Error",
  "error": "error description"
}
```


#### **GET** /api/http
##### Accepts
```json
{
  "alias": "http://0.0.0.0:8000/Y2Wv0YnlFe"
}
```
##### Returns
```json
{
  "status": "OK",
  "url": "https://ozon.ru"
}
```
or
```json
{
  "status": "Error",
  "error": "error description"
}
```
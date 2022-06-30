# JWT Authentication service
Простой JWT сервис аутентификации, написанный на Go.

## Запустить сервис
`docker-compose up`


## REST API
Описание API представлено ниже.

### Generate
Сгенерировать пару `AccessToken`, `RefreshToken` для предоставленного в http запросе UUID.  
`POST /generate?uuid={uuid}`
```
curl --request POST \
  --url 'http://localhost:8080/generate?uuid=79c62d874a7249be844cb9292388b6ee'
```
#### Response
```
{
	"accessToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhkWFJvYjNKcGVtVmtJanAwY25WbExDSmxlSEFpT2pFMk5UWTJNVEkzTlRNc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLlpBMXdDQUZqMlhjdThZZUhkYk5jZk1ISWZZUGM2eVFpbEdmVWZaY3NYUFduc3BmelB3WEhzeUlZNmI5bGIyMG8xUVY3OVVXUkFqN0FnenpIcnVoTExB",
	"refreshToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk5UWTJPVGd5TlRNc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLkx2NkltTWZ2OEVkUnFlUm84SDhXSmVtZUJPbVdrbWFKUHVJUDhxUG9zZmNYeGJVS2otUUtfbVlMenZrUUpObjExZE10MlAwNVBqV1I1VFM0NUZNTlZn"
}

```


### Refresh
Обновить пару `Access`, `Refresh` токенов. Операция успешна только при использовании `Refresh` токена, выданного с `Access` токеном.  
`POST /refresh`
```
curl --request POST \
  --url http://localhost:8080/refresh \
  --header 'Content-Type: application/json' \
  --data '{
	"accessToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhkWFJvYjNKcGVtVmtJanAwY25WbExDSmxlSEFpT2pFMk5UWTJNVEkzTlRNc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLlpBMXdDQUZqMlhjdThZZUhkYk5jZk1ISWZZUGM2eVFpbEdmVWZaY3NYUFduc3BmelB3WEhzeUlZNmI5bGIyMG8xUVY3OVVXUkFqN0FnenpIcnVoTExB",
	"refreshToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk5UWTJPVGd5TlRNc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLkx2NkltTWZ2OEVkUnFlUm84SDhXSmVtZUJPbVdrbWFKUHVJUDhxUG9zZmNYeGJVS2otUUtfbVlMenZrUUpObjExZE10MlAwNVBqV1I1VFM0NUZNTlZn"
}'
```
#### Response
```
{
	"accessToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmhkWFJvYjNKcGVtVmtJanAwY25WbExDSmxlSEFpT2pFMk5UWTJNVEkzTmpJc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLkVQY3Q5OHk1SDNSLU5MdnNtUkFSZlhxVjY5SmF1cjdUU2h1MWl5LXRmVTNtbkN0SmpMV09hN0N2Z3VHMndyNkczTEtUNUtfRktzb2k5blZGRzljS1Nn",
	"refreshToken": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk5UWTJPVGd5TmpJc0luVjFhV1FpT2lJM09XTTJNbVE0TnkwMFlUY3lMVFE1WW1VdE9EUTBZeTFpT1RJNU1qTTRPR0kyWldVaWZRLkhDMkY2bk82Zkhzc3ZmM25HN1VLSXBCTHRSZmp1REtzeGJpcE5GOTBTQkFucDdtdjNpMjZRbDNyd3JzZzdkdm5lMnNLMFpRRkhGX3F6dmdKcEcyVG5R"
}
```
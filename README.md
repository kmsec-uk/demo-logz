# Demo log panel

How easy is it to start stealing data? Created this for a lunch and learn to demonstrate the risk of infostealers.

Setup:
```bash
# Download picocss stylesheet
curl 'https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css' > ./static/pico.min.css
```

Run server:
```bash
AUTH_HEADER="foo" go run main.go
```

Send:

```bash
curl -H "authorization: foo" -d '{"something" : "valuabledata"}' localhost:8080/send
```

View:

```bash
curl -H "authorization: foo"  localhost:8080/view | jq
# Response:
#{
#  "something": "valuabledata"
#}
```
(alternatively you can load up localhost:8080 on your browser at and hit "refresh" -- don't forget to provide your authentication token!)

Clear:

```bash
curl -H "authorization: foo"  localhost:8080/clear
```
# WORD物協

割と真面目につくってしまった

## install

1. git clone & cd 
2. go build bin/main.go
3. touch items.json
4. touch users.json

## 説明

### items.json

販売している商品が格納されています

例:

```json
[
    {
        "name":"ジョージア 500ml",
        "id":"0000XSNJG0MQJHBF4QX1EFD6Y3",
        "price":69
    },
    {
        "name":"マックスコーヒー 350ml",
        "id":"01H4EDPJ3GZTGAK1PDFM3MTJNP",
        "price":90
    }
]
```

### users.json

ユーザ名（ここではUTID-13を想定）と、残額を格納しています

```json
{
    "0010122135035": 720,
    "001012213502X": 0
}

```
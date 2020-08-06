# saucenao-service-go
SauceNao service based on golang

Click [here](https://saucenao.com/user.php?page=search-api) get the official API key of saucenao and view the documents

![](https://tuchuang.laji.blog/imgs/2020/08/e6fff1389ce57cfb.jpg)

===========

[ GET ]

![](https://tuchuang.laji.blog/imgs/2020/08/ebd60049e95a4b92.png)
```
{
    api_key: xxxxxxxxxxxxxxxxxx,
    url: xxxxxxxxxxxxxxxxxxx,
    output_type: 2,
    testmode: 1,
    numres: 10,
    db: 999,
    minsim: 80,
}
```
===========

[ POST ]

![](https://tuchuang.laji.blog/imgs/2020/08/505c8e2ae74ce4a1.jpg)
```
let formData = new FormData();
formData.append("api_key", "xxxxxxxxxxxxxxxx");
formData.append("output_type", "2");
formData.append("testmode", "1");
formData.append("numres", "10");
formData.append("db", "999");
formData.append("minsim", "80");
formData.append("file", your_file);
```
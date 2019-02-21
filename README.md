# geoip2redis

 Loader of Multiple GeoIP providers to Redis

We currently support the following providers (* partially)

* [IP2Location](https://lite.ip2location.com/database/ip-country)
* [Software77](http://software77.net/geo-ip/)


geoip2redis primarilly supports DB1 from IP2Location, but with Auto mode enabled can load any of their standard IPv4 databases ~~, including their ASN database in DB1 format~~. (ASN is currently broken)
It can also load Software77's database either in it's native format, or convert it on the fly to IP2Location format, making them interchangeable.

### Examples:

![Use example1](https://user-images.githubusercontent.com/691270/53105684-8b38b400-356c-11e9-8cdd-ac0c76a7b64a.png)

![Redis output](https://user-images.githubusercontent.com/691270/53105706-92f85880-356c-11e9-9c2d-83b6c88f4a76.png)


Please install the [golang/vgo](https://github.com/golang/go/wiki/vgo) package before running ./build.sh to build the app.


Please check out the [the Wiki](https://github.com/ConsulTent/geoip2redis/wiki) for more info.

### TODO
A LOT
1. testing - need to add golang tests
2. Deal with conversion formats
3. Add support for Redis passwords
4. Add more GEOIP sources.   Suggestions welcome!



(c) 2019 ConsulTent Ltd.  http://consultent.ltd

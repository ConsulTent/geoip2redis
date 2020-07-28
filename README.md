![GeoIP2Redis](https://user-images.githubusercontent.com/691270/73528553-0e86a900-4450-11ea-80a8-5d603ddfbfd7.png)

![Go](https://github.com/ConsulTent/geoip2redis/workflows/Go/badge.svg?branch=master)

Loader of Multiple GeoIP providers to Redis.  Now with **[LIVE MIGRATION](https://github.com/ConsulTent/geoip2redis/wiki/Live-Migration)**!

##### The goal of GeoIP2Redis is to standardise all GeoIP formats into a standard CSV format (based on Ip2Location) that can be queried using Redis via subkeys.


We currently support the following providers (* partially)

* [IP2Location](https://lite.ip2location.com/database/ip-country)
* [Software77](http://software77.net/geo-ip/)
* [Maxmind](https://www.maxmind.com/en/geoip2-databases)*  (See tools/tools/maxmind-ip2location and [Wiki](https://github.com/ConsulTent/GeoIP2Redis/wiki))

GeoIP2Redis primarily supports DB1 from IP2Location, but with Auto mode enabled can load any of their standard IPv4 databases.
It can also load Software77's database either in it's native format, or convert it on the fly to IP2Location format, making them interchangeable.  *
Please see the Wiki for Maxmind support.


### Examples:

![Use example1](https://user-images.githubusercontent.com/691270/53105684-8b38b400-356c-11e9-8cdd-ac0c76a7b64a.png)

![Redis output](https://user-images.githubusercontent.com/691270/53105706-92f85880-356c-11e9-9c2d-83b6c88f4a76.png)


Please check out the [the Wiki](https://github.com/ConsulTent/GeoIP2Redis/wiki) for more info.

### TODO
1. Testing - need to add golang tests
2. ASN support for ip2location.
3. Add more GEOIP sources.   Suggestions welcome!



(c) 2020 ConsulTent Ltd.  http://consultent.ltd

<a href="https://www.ip2location.com/?rid=1415"><img src="https://www.ip2location.com/assets/img/affiliate_728x90.jpg" width="728" height="90" alt="Leading IP Geolocation solution provider to pinpoint the location of an IP address" /></a>

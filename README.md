![GeoIP2Redis](https://user-images.githubusercontent.com/691270/73528553-0e86a900-4450-11ea-80a8-5d603ddfbfd7.png)

![Go](https://github.com/ConsulTent/geoip2redis/workflows/Go/badge.svg?branch=master)

Loader of Multiple GeoIP providers to Redis.  Now with **[LIVE MIGRATION](https://github.com/ConsulTent/geoip2redis/wiki/Live-Migration)**!

![geoip-docker2](https://user-images.githubusercontent.com/691270/166111897-80b09371-3eaa-498a-adc8-77e6e5c4cff0.gif)

##### The goal of GeoIP2Redis is to standardise all GeoIP formats into a standard CSV format (based on Ip2Location) that can be queried using Redis via subkeys.


**Why GeoIP via Redis?**  

* It's FAST!  
* Much faster than a REST API.
* Low Latency:
  * Local vs Internet traffic
  * Much less overhead via Redis protocol vs HTTP

We currently fully support the following providers:

* [IP2Location](https://lite.ip2location.com/database/ip-country)
* ~~[Software77](http://software77.net/geo-ip/)~~
* [MaxMind](https://www.maxmind.com/en/geoip2-databases)  (See [Wiki](https://github.com/ConsulTent/GeoIP2Redis/wiki))

---

+ GeoIP2Redis primarily supports DB1 & DB3 DB formats from IP2Location.  Through conversion, we also fully support MaxMind's Block and Location databases.  We also provide a tool to convert MaxMind's Blocks and Location tables to a unified IP2Location format.
+ ~~It can also load Software77's database either in it's native format, or convert it on the fly to IP2Location format, making them interchangeable.  *~~  Software77 is gone, and support is now deprecated.



Please see the Wiki for usage and additional MaxMind support.

---

Basic **Telemetry** is included.  Please see the [Wiki](https://github.com/ConsulTent/geoip2redis/wiki/Telemetry).

**Coming Soon**
* [ipgeolocation](https://ipgeolocation.io/)
* [dbip](https://db-ip.com/db/)




### Examples:

![Use example1](https://user-images.githubusercontent.com/691270/53105684-8b38b400-356c-11e9-8cdd-ac0c76a7b64a.png)

![Redis output](https://user-images.githubusercontent.com/691270/53105706-92f85880-356c-11e9-9c2d-83b6c88f4a76.png)


Please check out the [the Wiki](https://github.com/ConsulTent/GeoIP2Redis/wiki) for more info.

### TODO
1. Testing - need to add golang tests
2. ASN support for ip2location.
3. Add more GEOIP sources.   **Suggestions welcome!**


### Disclaimer
IP2Location, Maxmind, Software77 and Redis are trademarks of their respective owners.


(c) 2021 ConsulTent Pte. Ltd.  https://consultent.ltd

<a href="https://www.ip2location.com/?rid=1415"><img src="https://www.ip2location.com/assets/img/affiliate_728x90.jpg" width="728" height="90" alt="Leading IP Geolocation solution provider to pinpoint the location of an IP address" /></a>

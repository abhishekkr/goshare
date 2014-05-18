## Changelog

### v0.8.2

* [push] encoding for ns/tsds may contain a root parent-ns config to pre-pend for all ASCII multi-val

---

#### v0.8.1

* added json support
* fixed multi-val support for non-tsds key-types

#### v0.8.0

* huge refactor and design fixes

#### *Good To Have*

>
> * [read,delete] enable {def,ns,tsds}-{,csv,json} to get multiple keys, this also define multi-val type for "read" output
>
> * (logging) "debug:default" log triggers only at failure status; "verbose" with only DBRest call logged again
>
> * [read] features for tsds like: {latest,this}_{year,month,week,day,hour,hour,min,dot}
>


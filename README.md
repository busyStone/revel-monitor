revel-monitor
=============

Use  beego / toolbox to add revel healthcheck, profile and statistics feature.

It works as a revel module.

Usage
=====

Add a line to app.conf

    module.monitor=github.com/busyStone/revel-monitor
    
Add route to routes

    module:monitor
    
Import in you app/controllers/init.go

    _ "github.com/busyStone/revel-monitor/app/controllers"
    
Use routes

* /@qps
* /@prof

    

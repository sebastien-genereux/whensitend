# When does it end?

A small web app that tells you when an activity ends, based on a start date and length. It will do this while excluding all those dates specified in the configuration.

Specify days off for the activity in configs/conf.json

Configuration parsing supports specific dates in dd/mm/yyyy format or any day of the week as per the golang time package weekday [type](https://pkg.go.dev/time#Weekday)

Note the start and end date are inclusive!
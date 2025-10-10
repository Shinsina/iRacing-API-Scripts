# iRacing API Scripts

Scripts for authenticating and processing data from the iRacing Web API

## Before Proceeding

You will need an iRacing `client_id` and `client_secret` the means of obtaining one is outlined in this [iRacing forum](https://forums.iracing.com/discussion/84226/legacy-authentication-removal-dec-9-2025/) thread.

Prior to anything run `./auth.sh` from the root of the project with a .env file containing the following

`IRACING_USERNAME=IRACING_USERNAME_GOES_HERE`
`IRACING_PASSWORD=IRACING_PASSWORD_GOES_HERE`
`IRACING_CLIENT_ID=IRACING_CLIENT_ID_GOES_HERE`
`IRACING_CLIENT_SECRET=IRACING_CLIENT_SECRET_GOES_HERE`

This will then create `token.json` in the root of the project for the current authentication session by which all other API requests require.

## Gathering Subsession IDs

Following authentication the first script to be run is `/search-series/main.go` at time of writing this will require all permutations of `season_quarter`, `season_year` and `cust_id` parameters to be set as an object within `/search-series/customer-id-season-quarter-season-year-mappings.json`. This will create the input file(s) that depending on how they are being used will need to be merged or used individually within `/subsession/main.go`.

## Gathering Subsession Data

Following runs of `/search-series/main.go`  for all `season_quarter`, `season_year` and `cust_id` permutations next we will be using the output file(s) and using them as the input file for `/subsession/main.go`, the file output here will be used [within this folder](https://github.com/Shinsina/Stat-N-Track/tree/master/Stat-N-Track) and should ultimately be a 1 dimensional array of Subsessions and saved as `1-subsessions-output.json` in the aforementioned folder

## Gathering Past Seasons Data

Following runs of `/search-series/main.go` for all `season_quarter`, `season_year` and `cust_id` permutations next we will be using the output file(s) and using them as the input file for `/past-seasons/main.go` as each result in the files has the `series_id` for each series participated in, these `series_id` values are to be gathered running `/utils/distinct-ids/main.go` for all subsessions exported and then saved within `/past-seasons` as `distinct-series-ids-output.json`, the output of this will ultimately go [within this folder](https://github.com/Shinsina/Stat-N-Track/tree/master/Stat-N-Track) saved as it is outputted named `past-seasons-output.json`.

## Gathering Standings Data

Following runs of `/search-series/main.go` for all `season_quarter`, `season_year` and `cust_id` permutations next will be using the output file(s) and using them as the input file for `/standings-jake/main.go` (for the iRacing account that the `client_id` and `client_secret` were generated) or `/standings-other/main.go` (for any other driver whose data is being exported) respectively. As each result in the files has the `season_id` and `car_class_id` for each series participated in. The file output here will be used [within this folder](https://github.com/Shinsina/Stat-N-Track/tree/master/Stat-N-Track) and should ultimately be a 1 dimensional array of Subsessions and saved as `standings-output.json` in the aforementioned folder

## Loading Car Class Data

Use the response from [this location](https://members-ng.iracing.com/data/carclass/get) (`/data/carclass/get`).

This can be accessed easiest by logging into the iRacing /data API [via this link](https://members-login.iracing.com/?ref=https%3A%2F%2Fmembers-ng.iracing.com%2Fdata%2Fdoc&signout=true) and then using the link provided above.

## Video Guide

COMING SOON

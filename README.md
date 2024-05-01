# iRacing API Scripts

Scripts for authenticating and processing data from the iRacing Web API

## Before Proceeding

Prior to anything run `./auth.sh` from the root of the project with a .env file containing the following

`IRACING_USERNAME=IRACING_USERNAME_GOES_HERE`
`IRACING_PASSWORD=IRACING_PASSWORD_GOES_HERE`

This will then create `cookie.txt` in the root of the project for the current authentication session by which all other API requests require.

## Gathering Subsession IDs

Following authentication the first script to be run is `iracingSearchSeriesExport.py` at time of writing this will require changing the query string on `query_string` for the respective `season_quarter`, `season_year` and `cust_id` parameters. This will create the input file for `iracingSubsessionExport.py`

## Gathering Subsession Data

Following runs of `iracingSearchSeriesExport.py` for all `season_quarter`, `season_year` and `cust_id` permutations next will be using the output file(s) and using them as the input file for `iracingSubsessionExport.py` the file output here will be used by `load-iracing-data/load-subsessions.js` to load subsession data into MongoDB.

## Gathering Past Seasons Data

Following runs of `iracingSearchSeriesExport.py` for all `season_quarter`, `season_year` and `cust_id` permutations next will be using the output file(s) and using them as the input file for `iracingPastSeasonsExport.py` as each result in the files has the `series_id` for each series participated in, these entries will additionally have the `season_id` for each series which can be used as `seasonIds` within `load-iracing-data/load-past-season-data.js` assuming they're mapped to an array and constructed as a unique set like so:

```js
  new Set(series.map((series) => series.season_id))
```

## Gathering Standings Data

Following runs of `iracingSearchSeriesExport.py` for all `season_quarter`, `season_year` and `cust_id` permutations next will be using the output file(s) and using them as the input file for `iracingStandingsReqJake.py` or `iracingStandingsReqJack.py` respectively. As each result in the files has the `season_id` and `car_class_id` for each series participated in. The file output here will be used by `load-iracing-data/load-standings.js` to load standings data into MongoDB. (NOTE: These will likely need to be filtered into unique set of string values of the format `SEASON-ID_CAR-CLASS-ID` and will need to be generated for each `cust_id` that needs to be processed).

## Loading Other Data

The following scripts within `load-iracing-data` use direct iRacing API responses and simply upload them to a respective MongoDB collection:

- `load-car-classes-data.js` Uses the response from [here](https://members-ng.iracing.com/data/carclass/get) (`/data/carclass/get`).
- `load-cars-data.js` Uses the response from [here](https://members-ng.iracing.com/data/car/get) (`/data/car/get`).
- `load-season-data.js` Uses the response from [here](https://members-ng.iracing.com/data/series/seasons) (`/data/series/seasons`).
- `load-users-data.js` Uses the response from [here](https://members-ng.iracing.com/data/member/info) (`/data/member/info`).

`/data/member/info` Needs to be run by the user the information is for and therefore requires the user to authenticate against the API to get their own user information.

`/data/carclass/get`, `/data/car/get` and `/data/series/seasons` should be run at the start of a subsequent season in order to get the most update to date information regarding each, the last of which (`/data/series/seasons`) is crucial as this changes every quarterly iRacing season and needs to be pulled prior to the season ending and is what is used for the scheduling portion of Stat 'n' Track.

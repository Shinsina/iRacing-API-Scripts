# iRacing API Scripts

Scripts for authenticating and processing data from the iRacing Web API

## Before Proceeding

Prior to anything run `./auth.sh` from the root of the project with a .env file containing the following

`IRACING_USERNAME=IRACING_USERNAME_GOES_HERE`
`IRACING_PASSWORD=IRACING_PASSWORD_GOES_HERE`

This will then create `cookie.txt` in the root of the project for the current authentication session by which all other API requests require.

## Gathering Subsession IDs

Following authentication the first script to be run is `iracingSearchSeasonsExport.py` at time of writing this will require changing the query string on `query_string` for the respective `season_quarter`, `season_year` and `cust_id` parameters. This will create the input file for `iracingSubsessionExport.py`

## Gathering Subsession Data

Following runs of `iracingSearchSeasonsExport.py` for all `season_quarter`, `season_year` and `cust_id` permutations next will be using the output file(s) and using them as the input file for `iracingSubsessionExport.py` the file output here will be used by `load-iracing-data/load-subsessions.js` to load subsession data into MongoDB.

## Gathering Past Seasons Data

Using the following query `await collection.distinct("series_id", {})` on [Stat 'n' Track](https://github.com/Shinsina/Stat-N-Track) within `pages/user/[id]/subsessions/index.astro` against the `subsessions` collection, this will return what should be saved as `past-season-series-ids-input.json` in the project root, then you can run `iracingPastSeasonsExport.py` which will output a file that is used by `loading-iracing-data/load-past-season-data.js` in conjunction with a respective `past-season-season-ids-input.json` (which is generated likewise to `past-season-series-ids-input.json` albeit using the following query `await collection.distinct("season_id", {})`) file to load past season data into MongoDB. (NOTE: All queries listed should be run following all subsessions being loaded into the `subsessions` collection)

## Gathering Standings Data

Using the results of the query on [Stat 'n' Track](https://github.com/Shinsina/Stat-N-Track) within `pages/user/[id]/subsessions/index.astro` against the `subsessions` collection, create a Set of string values of the following shape `SEASON-ID_CAR-CLASS-ID` and save this result as a JSON file for each `cust_id` that needs to be processed and utilize the respective file for each `cust_id` in `iracingStandingsReqJake.py` or `iracingStandingsReqJack.py` respectively. (NOTE: All queries listed should be run following all subsessions being loaded into the `subsessions` collection)

The outlined operation can be accomplished like so:

```js
Array.from(new Set(subsessions.map((subsession) => {
  const { session_results, season_id } = subsession;
  const user = session_results[2].results.find((v) => v.cust_id === Number(id));
  return `${season_id}_${user.car_class_id}`;
})));
```

## Loading Other Data

The following scripts within `load-iracing-data` use direct iRacing API responses and simply upload them to a respective MongoDB collection:

- `load-car-classes-data.js` Uses the response from [here](https://members-ng.iracing.com/data/carclass/get) (`/data/carclass/get`).
- `load-cars-data.js` Uses the response from [here](https://members-ng.iracing.com/data/car/get) (`/data/car/get`).
- `load-season-data.js` Uses the response from [here](https://members-ng.iracing.com/data/series/seasons) (`/data/series/seasons`).
- `load-users-data.js` Uses the response from [here](https://members-ng.iracing.com/data/member/info) (`/data/member/info`).

`/data/member/info` Needs to be run by the user the information is for and therefore requires the user to authenticate against the API to get their own user information.

`/data/carclass/get`, `/data/car/get` and `/data/series/seasons` should be run at the start of a subsequent season in order to get the most update to date information regarding each, the last of which (`/data/series/seasons`) is crucial as this changes every quarterly iRacing season and needs to be pulled prior to the season ending and is what is used for the scheduling portion of Stat 'n' Track.

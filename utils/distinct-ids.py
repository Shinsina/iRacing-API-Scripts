import json

subsessions = []
with open('./subsessions-output.json', 'r') as input_file:
  json_data = json.load(input_file)
  for value in json_data:
    subsessions.append(value)

series_ids = set()
season_ids = set()
for subsession in subsessions:
  if subsession['series_id'] not in series_ids:
    series_ids.add(subsession['series_id'])
  if subsession['season_id'] not in season_ids:
    season_ids.add(subsession['season_id'])

with open('distinct-series-ids-output.json', 'w') as series_ids_output_file:
  json.dump(list(series_ids), series_ids_output_file)

with open('distinct-season-ids-output.json', 'w') as season_ids_output_file:
  json.dump(list(season_ids), season_ids_output_file)

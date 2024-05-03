import json

subsessions = []
with open('../5-1-2024-jake-search-series-output.json', 'r') as input_file:
  json_data = json.load(input_file)
  for key in json_data.keys():
    subsessions.extend(json_data[key])

unique_season_car_class_id_mappings = set()
for subsession in subsessions:
  mapping = str(subsession['season_id']) + '_' + str(subsession['car_class_id'])
  if mapping not in unique_season_car_class_id_mappings:
    unique_season_car_class_id_mappings.add(mapping)

with open('standings-input.json', 'w') as output_file:
  json.dump(list(unique_season_car_class_id_mappings), output_file)

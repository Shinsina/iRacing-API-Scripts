import json
import requests

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

past_series = []
with open ('distinct-series-ids-output.json', 'r') as input_file:
    json_data = json.load(input_file)
    past_series.extend(json_data)

output = {}
query_string = 'https://members-ng.iracing.com/data/series/past_seasons?series_id={}'
for series in past_series:
    response = session.get(query_string.format(series))
    response_json = response.json()
    series_response = session.get(response_json['link'])
    series_response_json = series_response.json()
    output[series] = series_response_json

with open('past-seasons-output.json', 'w') as output_file:
  json.dump(output, output_file)

import json
import requests

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

subsessions = []
with open ('5-1-2024-search-series-output.json', 'r') as input_file:
    json_data = json.load(input_file)
    for key in json_data.keys():
        subsessions.extend(json_data[key])

output = {}
query_string = 'https://members-ng.iracing.com/data/results/get?subsession_id={}'
for index, subsession in enumerate(subsessions):
    print(str(index + 1) + ' of ' + str(subsessions.__len__()))
    response = session.get(query_string.format(subsession['subsession_id']))
    response_json = response.json()
    subsession_response = session.get(response_json['link'])
    subsession_response_json = subsession_response.json()
    output[subsession['subsession_id']] = subsession_response_json

with open('subsessions-output.json', 'w') as output_file:
    json.dump(output, output_file)

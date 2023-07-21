import json
import requests

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

subsessions = []
with open ('6-7-2023-jake-input.json', 'r') as input_file:
    json_data = json.load(input_file)
    for result in json_data['d']['r']:
        subsessions.append(result['41'])

output = {}
query_string = 'https://members-ng.iracing.com/data/results/get?subsession_id={}'
for subsession in subsessions:
    response = session.get(query_string.format(subsession))
    response_json = response.json()
    subsession_response = session.get(response_json['link'])
    subsession_response_json = subsession_response.json()
    output[subsession] = subsession_response_json

with open('output.json', 'w') as output_file:
    json.dump(output, output_file)
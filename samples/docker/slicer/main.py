import json

print('hola')

with open('/share/input.json') as data_file:
    inputdata = json.load(data_file)

print('-- input --')
print(inputdata)
print(inputdata['size'])
print(inputdata['slice'])

for i in inputdata['slice']:
    print(i)

outputdata = {"nextoutput": [1,2,3,4]}

with open('/share/output.json', 'w') as outfile:
    json.dump(outputdata, outfile)

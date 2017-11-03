
import os
import json
import subprocess
from subprocess import call

# Read input data
with open('/share/input.json') as data_file:
    inputdata = json.load(data_file)

# Folder
FOLDER = '/share/'

# Input
INPUT_IMAGE = inputdata["image"]
INPUT_BVALUES = inputdata["bvalues"]

# Output
OUTPUT_IMAGE = os.environ["image"]

numBvalues = len(INPUT_BVALUES)

# Se debe cambiar el nombre del archivo a .nii porque sino no reconoce el header
os.rename(FOLDER + INPUT_IMAGE, FOLDER + INPUT_IMAGE + ".nii")
INPUT_IMAGE = INPUT_IMAGE + ".nii"

print 'FOLDER', FOLDER
print 'INPUT IMAGE', INPUT_IMAGE
print 'INPUT BVALUES', INPUT_BVALUES
print 'OUTPUT IMAGE', OUTPUT_IMAGE

bvalues = [str(i) for i in INPUT_BVALUES]
print 'INPUT BVALUES PLAIN => ', bvalues

# cmd call

f = open(FOLDER+OUTPUT_IMAGE, 'w')
cmd = ["/home/models/modeloIVIM", FOLDER+INPUT_IMAGE, str(numBvalues)] + bvalues
print(cmd)

call(cmd, stdout=f)

# Guardar el output
outputdata = {"image": OUTPUT_IMAGE}

with open('/share/output.json', 'w') as outfile:
    json.dump(outputdata, outfile)

#0,10,20,30,40,50,100,150,200,400,600,800,1000,1700

#0 10 20 30 40 50 100 150 200 400 600 800 1000 1700

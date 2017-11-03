
import os
import json
import subprocess
from subprocess import call
import shutil

# Read input data
with open('/share/input.json') as data_file:
    inputdata = json.load(data_file)

# Folder
FOLDER = '/share/'

# Input
INPUT_IMAGE = inputdata["image"]

# Output
OUTPUT_IMAGE = os.environ["image"]

# Se debe cambiar el nombre del archivo a .nii porque sino no reconoce el header
os.rename(FOLDER + INPUT_IMAGE, FOLDER + INPUT_IMAGE + ".nii")
INPUT_IMAGE = INPUT_IMAGE + ".nii"

print 'FOLDER', FOLDER
print 'INPUT IMAGE', INPUT_IMAGE
print 'OUTPUT IMAGE', OUTPUT_IMAGE

# ConfidenceConnected3D necesita tener la extension .nii
TEMP_IMAGE="temp.nii"

# cmd call

f = open(FOLDER+OUTPUT_IMAGE, 'w')
cmd = ["/home/itk/bin/ConfidenceConnected3D", FOLDER+INPUT_IMAGE, TEMP_IMAGE]
print(cmd)

call(cmd, stdout=f)

shutil.move(TEMP_IMAGE, FOLDER+OUTPUT_IMAGE)

# Guardar el output
outputdata = {"image": OUTPUT_IMAGE}

with open('/share/output.json', 'w') as outfile:
    json.dump(outputdata, outfile)

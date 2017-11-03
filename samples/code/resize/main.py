
import os
import json

# Import the modules
from PIL import Image, ImageFilter

# Read input data
with open('/share/input.json') as data_file:
    inputdata = json.load(data_file)

# Folder
FOLDER = '/share/'

# Input
INPUT_IMAGE = inputdata["image"]
INPUT_SIZE = inputdata["size"]

# Output
OUTPUT_IMAGE = os.environ["image"]

# Logs
print 'Image', INPUT_IMAGE
print 'SIZE', INPUT_SIZE
print 'OUTPUT', OUTPUT_IMAGE

try:
    # Load an image from the hard drive
    original = Image.open(FOLDER + INPUT_IMAGE)

    # Resize the image
    new_img = original.resize((INPUT_SIZE, INPUT_SIZE))
    
    # save the new image
    new_img.save(FOLDER + OUTPUT_IMAGE, "JPEG")
    
except Exception as e:
    print "Unable to load image"
    print e

# Guardar el output
outputdata = {"image": OUTPUT_IMAGE}

with open('/share/output.json', 'w') as outfile:
    json.dump(outputdata, outfile)


import os
import colorsys
import json

# Import the modules
from PIL import Image, ImageFilter

def hueChange(img, hue):

    img.load()
    r, g, b = img.split()
    r_data = []
    g_data = []
    b_data = []

    for rd, gr, bl in zip(r.getdata(), g.getdata(), b.getdata()):
        h, s, v = colorsys.rgb_to_hsv(rd / 255., bl / 255., gr / 255.) 
        rgb = colorsys.hsv_to_rgb(hue/360., s, v)
        rd, gr, bl = [int(x*255.) for x in rgb]
        r_data.append(rd)
        g_data.append(gr)
        b_data.append(bl)

    r.putdata(r_data)
    g.putdata(g_data)
    b.putdata(b_data)

    return Image.merge('RGB',(r,g,b))

# Read input data
with open('/share/input.json') as data_file:
    inputdata = json.load(data_file)

# Folder
FOLDER = '/share/'

# Input
INPUT_IMAGE = inputdata["image"]
INPUT_HUE = int(inputdata["hue"])

# Output
OUTPUT_IMAGE = os.environ["image"]

# Logs
print 'Image', INPUT_IMAGE
print 'HUE', INPUT_HUE
print 'OUTPUT', OUTPUT_IMAGE

try:
    # Load an image from the hard drive
    original = Image.open(FOLDER + INPUT_IMAGE).convert('RGB')

    # Change the hue
    new_img = hueChange(original, INPUT_HUE)

    # save the new image
    new_img.save(FOLDER + OUTPUT_IMAGE, "JPEG")

except Exception as e:
    print "Unable to load image"
    print e

# Guardar el output
outputdata = {"image": OUTPUT_IMAGE}

with open('/share/output.json', 'w') as outfile:
    json.dump(outputdata, outfile)


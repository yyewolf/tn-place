import json

f = open('place.json')
p = json.load(f)
f.close()

# p is a 2d array of the form [[email, email], [email, email], ...]
counts = {}
for i in p:
    for j in i:
        if j in counts:
            counts[j] += 1
        else:
            counts[j] = 1

# sort the counts
sorted_counts = sorted(counts.items(), key=lambda x: x[1], reverse=True)

# save the sorted counts
f = open('sorted_counts.txt', 'w')
for i in sorted_counts:
    f.write(str(i) + '\n')
f.close()


# Create a gif from log.txt
f = open('log.txt')

logs = []

# Only read lines with 'User 114832833234792776424 placed pixel at (78, 3) with color {255 0 255 255}'
for line in f:
    if 'placed pixel' not in line:
        continue
    # We get the color and the coordinates
    color = line.split('with color ')[1].strip()
    rgba = color[1:-1].split(' ')
    position = line.split('at ')[1].split(' with')[0][1:-1].split(', ')
    logs.append([int(position[0]), int(position[1]), int(rgba[0]), int(rgba[1]), int(rgba[2]), int(rgba[3])])
    
    
# Gif settings
import glob
from PIL import Image

current_frame = 0
max_frame = len(logs)
real_count = 0

# Original size is 200x200 but we upscale for better quality
frame = Image.new('RGBA', (1000, 1000), (255, 255, 255, 255))

for i in range(max_frame):
    # Draw a square of 5x5 pixels
    frame.paste((logs[i][2], logs[i][3], logs[i][4], logs[i][5]), (logs[i][0] * 5, logs[i][1] * 5, logs[i][0] * 5 + 5, logs[i][1] * 5 + 5))
    # Save every 4 frames
    if i % 4 == 0:
        frame.save('frames/' + str(real_count).zfill(len(str(max_frame))) + '.png')
        real_count += 1
    current_frame += 1
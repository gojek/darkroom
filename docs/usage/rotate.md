---
id: rotate
title: Rotate
---


## Rotate

The `rot` parameter can be used to rotate the image clockwise for a certain degree.
Image can be rotated upto any degree till 360 and `rot` can have any float value.

| `?w=500&h=250` | `?w=500&h=250&rot=90` |`?w=500&h=250&rot=180` |
|:---:|:---:|:---:|
| {@injectImage: sample-image.jpg?w=500&h=250} | {@injectImage: sample-image.jpg?w=500&h=250&rot=90} | {@injectImage: sample-image.jpg?w=500&h=250&rot=180} |


## Flip

The `flip` parameter can be used to flip the image vertically or horizontally using values `v` or `h` respectively.

| `?w=500&h=250` | `?w=500&h=250&flip=v` | `?w=500&h=250&flip=h` |
|:---:|:---:|:---:|
| {@injectImage: sample-image.jpg?w=500&h=250} | {@injectImage: sample-image.jpg?w=500&h=250&flip=v} | {@injectImage: sample-image.jpg?w=500&h=250&flip=h} |
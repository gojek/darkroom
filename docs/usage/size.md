---
id: size
title: Size
---

The size parameters allow you to resize, crop and fit-to-crop your image.

## Crop Mode | `crop`

Crop mode controls how the image is aligned when `fit=crop` is set. The `w` and `h` parameters should also be set, so that the crop behavior is defined within specific image dimensions.

Valid values are `top`, `bottom`, `left` and `right`. Multiple values can be used by separating them with a comma `,`. If no value is explicitly set, the default behavior is to crop from the center of the image.

#### Directional Cropping
The `top`, `bottom`, `left`, and `right` values allow you to specify the starting location of the crop. Image dimensions will be calculated from this starting point outward. These values can be combined by separating with commas, e.g. `crop=top,left`.

- `top`: Crop from the top of the image, down.
- `bottom`: Crop from the bottom of the image, up.
- `left`: Crop from the left of the image, right.
- `right`: Crop from the right of the image, left.  


| `?w=500&h=250&fit=crop` | `?w=500&h=250` |
|:---:|:---:|
| ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=500&h=250&fit=crop) | ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=500&h=250) |

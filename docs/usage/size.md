---
id: size
title: Size
---

The size parameters allow you to resize, crop and fit-to-crop your image.

## Fit

Fit mode can be used to enforce crop on an image. If this is not set, the default behaviour is to resize the image while maintaing original aspect ratio. The `w` and `h` parameters should also be set, so that the crop is defined within specific image dimensions.

| `?w=500&h=250&fit=crop` | `?w=500&h=250` |
|:---:|:---:|
| ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=500&h=250&fit=crop) | ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=500&h=250) |

## Crop
Crop mode controls the focus point of image when `fit=crop` is set. The `w` and `h` parameters should also be set, so that the crop is defined within specific image dimensions.

Available values are `top`, `bottom`, `left` and `right`. More than one value can be used by separating them with a comma `,`. If crop mode is not set and `fit=crop` is set, it'll crop from the center of the image.

#### Changing Focus Point
The `top`, `bottom`, `left`, and `right` values allow you to specify the starting location of the crop. Image dimensions will be calculated from this starting point outward. These values can be combined by separating with commas, e.g. `crop=top,left`.

- `top`: Crop from the top of the image, down.
- `bottom`: Crop from the bottom of the image, up.
- `left`: Crop from the left of the image, right.
- `right`: Crop from the right of the image, left.  

| `?w=250&h=250&fit=crop&crop=left` | `?w=250&h=250&fit=crop&crop=right` |
|:---:|:---:|
| ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=250&h=250&fit=crop&crop=left) | ![image](https://kdarkroom.herokuapp.com/sample-image.jpg?w=250&h=250&fit=crop&crop=right) |

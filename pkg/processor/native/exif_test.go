// This code is modified from https://github.com/disintegration/imageorient

package native

import (
	_ "image/jpeg"
	"os"
	"testing"
)

var testFiles = []struct {
	path        string
	orientation int
}{
	{"./_testdata/overlay.png", 0},
	{"./_testdata/exif_orientation/expected.jpg", 0},
	{"./_testdata/exif_orientation/f2t.jpg", 2},
	{"./_testdata/exif_orientation/f3t.jpg", 3},
	{"./_testdata/exif_orientation/f4t.jpg", 4},
	{"./_testdata/exif_orientation/f5t.jpg", 5},
	{"./_testdata/exif_orientation/f6t.jpg", 6},
	{"./_testdata/exif_orientation/f7t.jpg", 7},
	{"./_testdata/exif_orientation/f8t.jpg", 8},
}

func TestReadOrientation(t *testing.T) {
	for _, tf := range testFiles {
		f, err := os.Open(tf.path)
		if err != nil {
			t.Fatalf("os.Open(%q): %v", tf.path, err)
		}

		o := readOrientation(f)
		if o != tf.orientation {
			t.Fatalf("expected orientation=%d but got %d (%s)", tf.orientation, o, tf.path)
		}
	}
}

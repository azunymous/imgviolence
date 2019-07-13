# imgviolence
Looks in a directory (deeply) in order to number, unify and resize all images. Renames images from 0.jpg, 1.jpg...N.jpg and resizes images to 100px, 200px and 500px. e.g the third processed image would become 2.jpg, 2_100.jpg, 2_200.jpg and 2_500.jpg

# Usage
./imgviolence <path to input directory> <path to output directory>

If the output directory does not exist, it will be created.
*Make sure the output directory is not inside the input directory at all*

package main

import (
	"fmt"
	"path"
)

/* Image size - power of two, degree - the number of G components of pixels,
into which the message size will be written.
A message is written to the rest of the pixels using the Koch algorithm.
On the decryption side, a similar algorithm */



func main(){
	img, pixels := selectFile(path.Join("pic/normal.bmp"))
	message := makeMessage(100)
	fmt.Println("Message := ", message)
	dct := makeDCT(pixels) // Двумерный массив зеленого цвета.

	dct = insertMessage(message, dct, img, 20, pixels)
	res := makeBackDCT(dct)

	createPicture(img, pixels, res, path.Join("pic/newImage20.bmp"))
	fmt.Println("\n \t \t \t Decoding \t \t \t ")

	img, pixels = selectFile(path.Join("pic/newImage20.bmp"))
	dct = makeDCT(pixels)
	messageRec := extractMessage(pixels, img, dct, 0.0)
	printBits(messageRec)
}

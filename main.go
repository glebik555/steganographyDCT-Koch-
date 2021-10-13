package main

import (
	"fmt"
	"path"
)


/*Размер картинки - степень двойки,
степень - кол-во G компонент пикселей, в которые запишется размер сообщения.
В остальные пиксели записывается сообщение алгоритмом Коха.
На расшифровывающей стороне аналогичный алгоритм*/


func main(){
	img, pixels := selectFile(path.Join("pic\\normal.bmp"))
	message := makeMessage(100)
	fmt.Println("Message := ", message)
	dct := makeDCT(pixels) // Двумерный массив зеленого цвета.

	dct = insertMessage(message, dct, img, 20, pixels)
	res := makeBackDCT(dct)

	createPicture(img, pixels, res, path.Join("pic\\newImage20.bmp"))
	fmt.Println("\n \t \t \t Decoding \t \t \t ")

	img, pixels = selectFile(path.Join("pic\\newImage20.bmp"))
	dct = makeDCT(pixels)
	messageRec := extractMessage(pixels, img, dct, 0.0)
	printBits(messageRec)
}

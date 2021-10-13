package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)


type Pixel struct {
	R int
	G int
	B int
	A int
}

const pi = 3.1415926535


func createPicture(img image.Image, pixels [][]Pixel, res [][]float64, name string) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	upLeft := image.Point{}

	lowRight := image.Point{X: width, Y: height}
	newImg := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cyan := color.RGBA{R: uint8(pixels[y][x].R), G: uint8(res[y][x]), B: uint8(pixels[y][x].B), A: uint8(pixels[y][x].A)}
			newImg.Set(x, y, cyan)
		}
	}

	s, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}

	err = bmp.Encode(s, newImg)
	if err != nil {
		return
	}
	err = s.Close()
	if err != nil {
		return
	}
	defer s.Close()
}

func makeBackDCT(arr [][]float64) [][]float64 {
	backDct := make([][]float64, len(arr))
	for i := range backDct {
		backDct[i] = make([]float64, len(arr[0]))
		for j := range arr[0] {
			backDct[i][j] = arr[i][j]
		}
	}

	for h := 0; h < len(arr); h += 8 {
		if h+8 >= len(arr) {
			break
		}
		for w := 0; w < len(arr[0]); w += 8 {
			if w+8 >= len(arr[0]) {
				break
			}
			for i := 0; i < 8; i++ {
				for j := 0; j < 8; j++ {
					sum := 0.0
					Ck, Cl := 0.0, 0.0
					for k := 0; k < 8; k++ {
						for l := 0; l < 8; l++ {
							if k == 0 {
								Ck = 1.0 / 8.0
							} else {
								Ck = 2.0 / 8.0
							}
							if l == 0 {
								Cl = 1.0 / 8.0
							} else {
								Cl = 2.0 / 8.0
							}

							sum += math.Sqrt(Cl) * math.Sqrt(Ck) * arr[h+k][w+l] *
								math.Cos(((2*float64(i)+1)*float64(k)*pi)/(2.0*8.0)) *
								math.Cos(((2*float64(j)+1)*float64(l)*pi)/(2.0*8.0))
						}
					}
					backDct[h+i][w+j] = math.Round(sum)
				}
			}
		}
	}
	return backDct
}

func makeDCT(arr [][]Pixel) [][]float64 {
	dct := make([][]float64, len(arr))
	for i := range dct {
		dct[i] = make([]float64, len(arr[0]))
		for j := range arr[0] {
			dct[i][j] = float64(arr[i][j].G)
		}
	}

	for h := 0; h < len(arr); h += 8 {
		if h+8 >= len(arr) {
			break
		}
		for w := 0; w < len(arr[0]); w += 8 {
			if w+8 >= len(arr[0]) {
				break
			}
			for k := 0; k < 8; k++ {
				for l := 0; l < 8; l++ {
					sum := 0.0

					for i := 0; i < 8; i++ {
						for j := 0; j < 8; j++ {
							sum += float64(arr[h+i][w+j].G) *
								math.Cos(((2.0*float64(i)+1.0)*float64(k)*pi)/(2.0*8.0)) *
								math.Cos(((2.0*float64(j)+1.0)*float64(l)*pi)/(2.0*8.0))
						}
					}
					Ck, Cl := 0.0, 0.0
					if k == 0 {
						Ck = 1.0 / 8.0
					} else {
						Ck = 2.0 / 8.0
					}
					if l == 0 {
						Cl = 1.0 / 8.0
					} else {
						Cl = 2.0 / 8.0
					}
					res := math.Sqrt(Ck) * math.Sqrt(Cl) * sum
					dct[k+h][l+w] = res
				}
			}
		}
	}

	return dct
}

func rgbaToPixel(r, g, b, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func getPixels(file io.Reader) ([][]Pixel, image.Image, error) {
	img, err := bmp.Decode(file)

	if err != nil {
		return nil, nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}
	return pixels, img, nil
}

func selectFile(path string) (image.Image, [][]Pixel) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	pixels, img, err := getPixels(f)

	if err != nil {
		log.Println(err)
		panic(1)
	}
	return img, pixels
}

func ConvertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(i, toBase), nil
}

func calculateSize(multiplication int) (degree, number int) {
	i := 0.0

	for {
		if (int(math.Pow(2, i)) < multiplication) && ((int)(math.Pow(2, i+1)) >= multiplication) {
			return int(i), int(math.Pow(2, i))
		}
		i++
	}
}

func insertMessage(message []uint, dct [][]float64, img image.Image, epsilon float64, pixels [][]Pixel) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	placeForMessageSize, pow := calculateSize(width * height)
	countOfBitPlace := 0
	binSizeM, _ := ConvertInt(strconv.Itoa(len(message)), 10, 2)
	countOfMessage := 0

	if len(message) > height*width/64 {
		panic("Size of message too large")
	}
	once := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if countOfBitPlace < len(binSizeM) { // Writing a binary number whose value is equal to the size of the message immediately following the size.
				if placeForMessageSize-(y*x+x+len(binSizeM)) > 0 {
					pixels[y][x].B &= 0xFE
					fmt.Print("0")
				} else {
					if binSizeM[countOfBitPlace] == '1' {
						pixels[y][x].B |= 0x01
						fmt.Print("1")
						countOfBitPlace++
					} else {
						pixels[y][x].B &= 0xFE
						fmt.Print("0")
						countOfBitPlace++
					}
				}
			} else {
				if (countOfMessage < len(message)) && countOfMessage <= pow { // Write message w/ Koch's algorithm
					if once == 0 {
						fmt.Println("")
						once++
					}
					if message[countOfMessage] == 0x01 { // 1
						k := math.Abs(dct[y+1][x+3]) - math.Abs(dct[y+2][x+5])
						if k >= -epsilon {
							if countOfMessage > 120 && countOfMessage < 130 {
								fmt.Println("dct[y+2][x+5] was: ", dct[y+2][x+5])
							}

							dct[y+2][x+5] = math.Abs(dct[y+1][x+3]) + epsilon + 1
							if countOfMessage > 120 && countOfMessage < 130 {
								fmt.Println("dct[y+2][x+5] Now: ", dct[y+2][x+5])
							}
						}

						fmt.Println("#", countOfMessage, "(1): dct[y+1][x+3] (", dct[y+1][x+3], ") - dct[y+2][x+5] (", dct[y+2][x+5], ") = (k)", math.Abs(dct[y+1][x+3])-math.Abs(dct[y+2][x+5]), "< ", -epsilon, "(epsilon)")

						x += 8
						if x+8 >= width {
							y += 7
							x = width
							countOfMessage++
							continue
						}
						x--
						countOfMessage++
					} else {
						k := math.Abs(dct[y+1][x+3]) - math.Abs(dct[y+2][x+5])

						if k <= epsilon {
							dct[y+1][x+3] = math.Abs(dct[y+2][x+5]) + epsilon + 1
						}
						fmt.Println("#", countOfMessage, "(0): dct[y+1][x+3](", dct[y+1][x+3], ") - dct[y+2][x+5](", dct[y+2][x+5], ") = ", math.Abs(dct[y+1][x+3])-math.Abs(dct[y+2][x+5]), "> ", epsilon)
						x += 8
						if x+8 >= width {
							y += 7
							x = width
							countOfMessage++
							continue
						}
						x--
						countOfMessage++
					}
				}
			}
		}
	}

	return dct
}

func extractMessage(pixelsRec [][]Pixel, imgRec image.Image, dct [][]float64, epsilon float64) []bool {

	boundsRec := imgRec.Bounds()
	widthRec, heightRec := boundsRec.Max.X, boundsRec.Max.Y
	placeForMessageRec, _ := calculateSize(widthRec * heightRec)
	countOfSymbolRec := 0   // The number of bits in the encrypted message
	countOfSizeMessage := 0 // Pure bit of space for message

	var lengthOfMessageInt int64

	var lengthOfMessageBin []string

	var containerText []bool // Decrypted message

	flagOfCountingValue := true

	for y := 0; y < heightRec; y++ {
		for x := 0; x < widthRec; x++ {
			if countOfSizeMessage < placeForMessageRec {
				if pixelsRec[y][x].B%2 == 1 {
					lengthOfMessageBin = append(lengthOfMessageBin, "1")
					countOfSizeMessage++
				} else {
					lengthOfMessageBin = append(lengthOfMessageBin, "0")
					countOfSizeMessage++
				}
			} else {
				if flagOfCountingValue {
					lengthOfMessageInt, _ = strconv.ParseInt(strings.Join(lengthOfMessageBin, ""), 2, 64)
					fmt.Println("Length: ", strings.Join(lengthOfMessageBin, ""), " -> ", lengthOfMessageInt)
					flagOfCountingValue = false
				}
				if countOfSymbolRec < int(lengthOfMessageInt) {
					k := math.Abs(dct[y+1][x+3]) - math.Abs(dct[y+2][x+5])

					if k > epsilon {
						fmt.Println("#", countOfSymbolRec, "(0): dct[y+1][x+3](", dct[y+1][x+3], ") - dct[y+2][x+5](", dct[y+2][x+5], ") = ", math.Abs(dct[y+1][x+3])-math.Abs(dct[y+2][x+5]), "> ", epsilon)
						containerText = append(containerText, false)
						countOfSymbolRec++
						x += 8
						if x+8 >= widthRec {
							y += 7
							x = widthRec
							continue
						}
						x--
					} else {
						containerText = append(containerText, true)
						fmt.Println("#", countOfSymbolRec, "(1): dct[y+1][x+3](", dct[y+1][x+3], ") - dct[y+2][x+5](", dct[y+2][x+5], ") = (k)", math.Abs(dct[y+1][x+3])-math.Abs(dct[y+2][x+5]), "< ", -epsilon, "(epsilon)")
						countOfSymbolRec++
						x += 8
						if x+8 >= widthRec {
							y += 7
							x = widthRec
							continue
						}
						x--
					}
				}
			}
		}
	}
	return containerText
}

func makeMessage(lenOfMessage int) []uint {
	message := make([]uint, lenOfMessage) // Example of message
	for i := 0; i < cap(message); i++ {
		if i%2 == 1 {
			message[i] = 0x00
		} else {
			message[i] = 0x01
		}
	}

	return message
}

func printBits(slice []bool) {
	for i := 0; i < len(slice); i++ {
		if slice[i] {
			fmt.Print(1)
		} else {
			fmt.Print(0)
		}
	}
	fmt.Println()
}

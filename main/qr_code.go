package main

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"project/main/model"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"github.com/tealeg/xlsx"
)

var (
	err    error
	inFile = "/Users/kevin/Downloads/H_Room.xlsx"
)

func main() {
	excel := readExcel()
	for _, v := range excel {
		createQRCodeWithBg("/Users/kevin/Documents/A5.png", v.Name, v.Code)
	}
}

// 读取Excel
func readExcel() []model.ClassroomModel {
	var classModelList []model.ClassroomModel
	// 打开文件
	xlFile, _ := xlsx.OpenFile(inFile)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	// 遍历sheet页读取
	for _, sheet := range xlFile.Sheets {
		fmt.Println("sheet name: ", sheet.Name)
		//遍历行读取
		for i := 1; i < len(sheet.Rows); i++ {
			var classModel model.ClassroomModel
			for j := 0; j < len(sheet.Rows[i].Cells); j++ {
				text := sheet.Rows[i].Cells[j].String()
				if j == 0 {
					classModel.SetCode(text)
				} else if j == 1 {
					classModel.SetName(text)
				}
			}
			classModelList = append(classModelList, classModel)
		}
	}
	fmt.Println("\n\nimport success")
	return classModelList
}

func createQRCodeWithBg(bgPath string, fileName string, content string) {
	var (
		bgFile    *os.File
		bgImg     image.Image
		qrCodeImg image.Image
		offset    image.Point
	)

	// 01: 打开背景图片
	bgFile, err = os.Open(bgPath)
	if err != nil {
		fmt.Println("打开背景图片失败", err)
		return
	}

	defer bgFile.Close()

	// 02: 编码为图片格式
	bgImg, err = png.Decode(bgFile)
	if err != nil {
		fmt.Println("背景图片编码失败:", err)
		return
	}

	// 03: 生成二维码
	qrCodeImg, err = createAvatar(content)
	if err != nil {
		fmt.Println("生成二维码失败:", err)
		return
	}

	offset = image.Pt(339, 90)

	b := bgImg.Bounds()

	m := image.NewRGBA(b)

	draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)

	draw.Draw(m, qrCodeImg.Bounds().Add(offset), qrCodeImg, image.Point{X: 0, Y: 0}, draw.Over)

	i, _ := os.Create("/Users/kevin/qrcode3/" + fileName + ".png")

	_ = png.Encode(i, m)

	defer i.Close()

}

// 生成头像
func createAvatar(content string) (image.Image, error) {
	var (
		bgImg      image.Image
		offset     image.Point
		avatarFile *os.File
		avatarImg  image.Image
	)

	bgImg, err = createQrCode(content)

	if err != nil {
		fmt.Println("创建二维码失败:", err)
		return nil, errors.New("创建二维码失败")
	}
	avatarFile, err = os.Open("/Users/kevin/Documents/logo.png")
	avatarImg, err = png.Decode(avatarFile)
	avatarImg = ImageResize(avatarImg, 40, 40)
	b := bgImg.Bounds()

	// 设置为居中
	offset = image.Pt((b.Max.X-avatarImg.Bounds().Max.X)/2, (b.Max.Y-avatarImg.Bounds().Max.Y)/2)

	m := image.NewRGBA(b)

	draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(m, avatarImg.Bounds().Add(offset), avatarImg, image.Point{X: 0, Y: 0}, draw.Over)

	return m, err
}

// 生成二维码
func createQrCode(content string) (img image.Image, err error) {
	var qrCode *qrcode.QRCode

	qrCode, err = qrcode.New(content, qrcode.Highest)

	if err != nil {
		return nil, errors.New("创建二维码失败")
	}
	qrCode.DisableBorder = false

	img = qrCode.Image(200)

	return img, nil
}

func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}

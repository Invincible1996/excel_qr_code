package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"project/main/model"
	"time"

	"github.com/tealeg/xlsx"
)

var (
	err    error
	inFile = "/Users/kevin/Downloads/H_Room.xlsx"
)

func main() {
	excel := readExcel()
	fmt.Println("开始生成", time.Now())
	for _, v := range excel {
		if v.Code != "" && v.Name != "" {
			var (
				bgFile    *os.File
				bgImg     image.Image
				qrCodeImg image.Image
				offset    image.Point
			)

			// 01: 打开背景图片
			bgFile, err = os.Open("./bg3.png")
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
			qrCodeImg, err = createAvatar(v.Code, v.Name)
			if err != nil {
				fmt.Println("生成二维码失败:", err)
				return
			}

			offset = image.Pt(206, 522)

			b := bgImg.Bounds()

			m := image.NewRGBA(b)

			draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)

			draw.Draw(m, qrCodeImg.Bounds().Add(offset), qrCodeImg, image.Point{X: 0, Y: 0}, draw.Over)
			addLabel(m, 206, 600, "Hello Go")
			// 上传至oss时这段要改
			i, _ := os.Create("/Users/kevin/qr-code/" + v.Name + ".png")

			_ = png.Encode(i, m)

			defer i.Close()
		}
	}
	fmt.Println("生成结束", time.Now())
}

// 读取Excel
func readExcel() []model.ClassroomModel {
	var classModelList []model.ClassroomModel
	// 打开文件
	xlFile, _ := xlsx.OpenFile(inFile)
	if err != nil {
		fmt.Println(err.Error())
		return []model.ClassroomModel{}
	}
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

//
func createQrCode(content string, name string) (image.Image, error) {
	var (
		qr  *qrcode.QRCode
		img image.Image
	)

	qr, err = qrcode.New(content, qrcode.High)

	if err != nil {
		return nil, err
	}
	qr.DisableBorder = true
	img = qr.Image(750)

	return img, nil
}

func createAvatar(content string, name string) (image.Image, error) {
	var (
		bgImg image.Image
		//offset     image.Point
		avatarFile *os.File
		avatarImg  image.Image
	)

	bgImg, err = createQrCode(content, name)

	if err != nil {
		fmt.Println("创建二维码失败", err)
		return nil, err
	}

	avatarFile, err = os.Open("./avatar.png")

	avatarImg, err = png.Decode(avatarFile)

	avatarImg = ImageResize(avatarImg, 40, 40)
	b := bgImg.Bounds()

	// 设置头像居中
	//offset = image.Pt((b.Max.X-avatarImg.Bounds().Max.X)/2, (b.Max.Y-avatarImg.Bounds().Max.Y)/2)
	m := image.NewRGBA(b)

	// 绘制二维码
	draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)

	// 将头像绘制在二维码中间
	//draw.Draw(m, avatarImg.Bounds().Add(offset), avatarImg, image.Point{X: 0, Y: 0}, draw.Src)

	return m, err
}

func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 200, G: 100, A: 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

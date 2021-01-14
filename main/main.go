package main

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/tealeg/xlsx"
	"image/png"
	"os"
	"project/main/model"
)

var (
	inFile = "/Users/kevin/Downloads/H_Room.xlsx"
)

func main() {
	excel := readExcel()
	fmt.Println(len(excel))
	for _, v := range excel {
		createQRCode(v.Name, v.Code)
	}
}

func createQRCode(filePath string, content string) {
	qrCode, _ := qr.Encode(content, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	file, _ := os.Create("/Users/kevin/qrcode/" + filePath + ".png")
	defer file.Close()
	_ = png.Encode(file, qrCode)
}

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

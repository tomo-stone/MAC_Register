package main

import (
	"log"
	"flag"
	"os"
	"encoding/csv"
	"io"
	"strings"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

// データの変換(必要に応じて増やしてください)
func replaceMAC(str string) string {
	temp := strings.Replace(str,"-",":",-1)
	return strings.ToLower(temp)
}

func replaceSlackID(str string) string {
	return strings.Replace(str," ","_",-1)
}

func checkRecord(record []string) bool {
	// ヘッダ部分を飛ばす
	if record[0]=="タイムスタンプ"{
		return true
	}
	// 空の行を飛ばす
	if record[0]==""{
		return true
	}
	// 追加済みレコードを飛ばす
	if record[3]=="TRUE"{
		return true
	}
	return false
}

func main(){
	flag.Parse()
	//読み込みファイル準備
	file1, err := os.Open(flag.Args()[0])
	failOnError(err)
	defer file1.Close()

	//書き込みファイル準備
	writeFile, err := os.OpenFile(flag.Args()[1],os.O_APPEND|os.O_WRONLY, 0755)
	failOnError(err)
	defer writeFile.Close()

	reader := csv.NewReader(file1) //utf8
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない

	writer := csv.NewWriter(writeFile)
	writer.UseCRLF = true //デフォルトはLFのみ

	// 読み書き
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}
		if checkRecord(record) == true{
			continue
		}

		// 書き込み
		newRecords := make([]string,2)
		newRecords[0] = replaceMAC(record[2])     //MAC address
		newRecords[1] = replaceSlackID(record[1]) //slack ID
		log.Printf("Input:%+v\n",newRecords)
		err1 := writer.Write(newRecords) // 1行書き出す
		failOnError(err1)
		writer.Flush()
	}
}

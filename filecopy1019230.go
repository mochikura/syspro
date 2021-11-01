package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func main() {
	// open input file and make a buffered reader
	//問題1：osの部分をsyscallに変えればOK
	//ただしArgsのところのみosを用いる
	fi, err := syscall.Open(os.Args[1], syscall.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Close(fi); err != nil {
			panic(err)
		}
	}()
	// open output file and make a buffered writer
	fo, err := syscall.Open(os.Args[2],
		syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Close(fo); err != nil {
			panic(err)
		}
	}()
	// make a buffer to read data
	buf := make([]byte, 1024)
	// copy the whole content of
	// the input file to the output file
	//問題2：bufにまず結果を入れて、次にtempに読み飛ばした結果を入れる
	temp := make([]byte, 1024)
	for {
		n, err := syscall.Read(fi, buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		//fmt.Println(buf)
		if n == 0 {
			break
		}
		//文字コードが指定の範囲に入っていれば、tempにその文字が入る
		//iはbufのindex、xはtempのindexを指す
		x := 0
		for i := 0; i < n; i++ {
			if buf[i] == 10 || buf[i] == 32 || buf[i] >= 48 && buf[i] <= 57 || buf[i] >= 65 && buf[i] <= 90 || buf[i] >= 97 && buf[i] <= 122 {
				temp[x] = buf[i]
				x++
			}
		}
		//xというtempのindexを利用して配列をfoに入れる
		if _, err := syscall.Write(fo, temp[:x]); err != nil {
			panic(err)
		}
	}
	//問題3：syscall.Fstatを用いてStat_tに出力を入れる。
	//入れたStat_tのサイズを見ることで大きさがわかる
	var stat syscall.Stat_t
	if err := syscall.Fstat(fo, &stat); err != nil {
		panic(err)
	}
	fmt.Printf("%dbyte\n", stat.Size)
	//少し詰まったところ：FstatとStat_tがWindowsでは動かない
	//(便利なのでWinのターミナルを使っていたが、動かず。Ubuntuから実行したら動いた)
}

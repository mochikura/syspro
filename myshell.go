package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	scanner := bufio.NewReader(os.Stdin)
	ia := 0 //履歴番号
	for {
		fmt.Printf("./myshell[%02d]>", ia)  //shell
		in, err := scanner.ReadString('\n') //テキスト読み取り
		if err == io.EOF {                  //Ctrl+Dのとき
			fmt.Printf("\n")
			os.Exit(0)
		}
		in = strings.Trim(in, "\n")
		if in == "bye" { //byeのとき
			break
		}

		com := strings.Split(in, " ")

		thflag := 0         //三項演算子の判定
		onenum := 0         //?の位置
		twonum := 0         //:の位置
		var onecom []string //?の後
		var twocom []string //:の後
		for n, a := range com {
			//fmt.Println(n)
			if a == "?" && thflag == 0 { //?があったら、flag更新して場所を入れる
				thflag = 1
				onenum = n
			}
			if a == ":" && thflag == 1 { //:があったらflag更新して場所を入れる
				thflag = 2
				twonum = n
				//com = com[:(onenum - 1)]
				onecom = com[(onenum + 1):(twonum)] //三項演算子に基づいて3つのコマンドを取得
				twocom = com[(twonum + 1):]
				com = com[:onenum]
			}
		}

		comnum := comgo(com) //コマンド実行

		if thflag == 2 && comnum == 1 { //失敗時
			//fmt.Println("twocom")
			comgo(twocom)
		}
		if thflag == 2 && comnum == 2 { //成功時
			//fmt.Println("onecom")
			comgo(onecom)
		}
		ia++ //番号更新
	}
}

func comgo(com []string) int {
	//fmt.Println(com)
	cpath, err := exec.LookPath(com[0]) //コマンドのパス検索
	//fmt.Println(cpath)
	if err != nil { //存在しないコマンドなら

		fmt.Println(com[0] + ": No such file or directory")
		//log.Fatalf("%s not found in $PATH.", com[0])
		return 1
	} else { //存在したら

		args := com[0:]
		attr := syscall.ProcAttr{Files: []uintptr{0, 1, 2}}
		pid, err := syscall.ForkExec(cpath, args, &attr)
		//pid取得
		if err != nil {
			fmt.Println(err) //Permisson denied
		}

		proc, err := os.FindProcess(pid) //プロセス取得
		//fmt.Println("dmt")
		status, err := proc.Wait() //プロセス実行待ち
		//fmt.Println(status)
		if err != nil {
			fmt.Println(err) //waitid: invalid argument
		}
		//プロセス成功しなければ
		if !status.Success() { //error
			fmt.Println(status.String())
		}
		return 2
	}
}

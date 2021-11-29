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

		comnum := -1
		thflag := 0 //三項演算子の判定
		onenum := 0 //?の位置
		twonum := 0 //:の位置
		rdflag := 0 //redirectとpipeの判定
		//rdnum := 0          //><の位置
		pipenum := 0        //|の位置
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
			if a == ">" && rdflag == 0 { //リダイレクト1の時
				rdflag = 1
				//rdnum = n
			}
			if a == "<" && rdflag == 0 { // リダイレクト2のとき
				rdflag = 2
				//rdnum = n
			}
			if a == "|" && rdflag == 0 { //パイプの時
				rdflag = 3
				pipenum = n
			}

		}
		if rdflag == 0 { //パイプリダイレクトなし
			comnum = comgo(com) //コマンド実行
		} else if rdflag == 1 { //リダイレクト1時
			rdlslice := strings.Split(in, " > ")
			comnum = rdcomgo_l(rdlslice[0], rdlslice[1])
		} else if rdflag == 2 {
			rdrslice := strings.Split(in, " < ")
			comnum = rdcomgo_r(rdrslice[0], rdrslice[1])
		} else if rdflag == 3 { //パイプ時
			pipecomgo(com[:(pipenum-1)], com[pipenum+1:])
			ia++
			continue
		}

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
		fmt.Println("OK")
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

func rdcomgo_l(com string, f string) int {
	// open a file which will be connected to stdout
	//fmt.Println("OK")
	coms := strings.Split(com, " ")
	fw, err := os.OpenFile(f,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := fw.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	//fmt.Println("OK")
	cmdpath, err := exec.LookPath(coms[0])
	if err != nil {
		fmt.Println(coms[0] + ": No such file or directory")
		return 1
	} else {
		//fmt.Println("OK")
		args := coms[0:]
		attr := syscall.ProcAttr{
			Files: []uintptr{0, fw.Fd(), fw.Fd()}}
		pid, err := syscall.ForkExec(cmdpath, args, &attr)
		if err != nil {
			fmt.Println(err)
		}
		proc, err := os.FindProcess(pid)
		status, err := proc.Wait()
		if err != nil {
			fmt.Println(err)
		}
		if !status.Success() {
			fmt.Println(status.String())
		}
		return 2
	}
}

func rdcomgo_r(com string, f string) int {
	// open a file which will be connected to stdout
	coms := strings.Split(com, " ")
	fw, err := os.OpenFile(f,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := fw.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	cmdpath, err := exec.LookPath(coms[0])
	if err != nil {
		fmt.Println(coms[0] + ": No such file or directory")
		return 1
	} else {
		//fmt.Println("OK")
		args := coms[0:]
		attr := syscall.ProcAttr{
			Files: []uintptr{fw.Fd(), 1, 2}}
		pid, err := syscall.ForkExec(cmdpath, args, &attr)
		if err != nil {
			fmt.Println(err)
		}
		proc, err := os.FindProcess(pid)
		status, err := proc.Wait()
		if err != nil {
			fmt.Println(err)
		}
		if !status.Success() {
			fmt.Println(status.String())
		}
		return 2
	}

}

func pipecomgo(com1 []string, com2 []string) int {
	// make a pipe
	pin, pout, err := os.Pipe()
	// execute the 1st command
	cmdpath1, err := exec.LookPath(com2[0])
	if err != nil {
		fmt.Println(com2[0] + ": No such file or directory")
		return 0
	}
	attr1 := syscall.ProcAttr{Env: os.Environ(),
		Files: []uintptr{pin.Fd(), 1, 2}}
	pid, err := syscall.ForkExec(cmdpath1, com2, &attr1)
	if err != nil {
		fmt.Println(err)
	}
	pin.Close()
	// execute the 2nd command
	cmdpath2, err := exec.LookPath(com1[0])
	if err != nil {
		fmt.Println(com1[0] + ": No such file or directory")
		return 0
	}
	attr2 := syscall.ProcAttr{Env: os.Environ(),
		Files: []uintptr{0, pout.Fd(), 2}}
	_, err = syscall.ForkExec(cmdpath2, com1, &attr2)
	if err != nil {
		fmt.Println(err)
	}
	pout.Close()
	// wait for the 2nd command to complete
	proc, err := os.FindProcess(pid)
	status, err := proc.Wait()
	if err != nil {
		fmt.Println(err)
	}
	if !status.Success() {
		fmt.Println(status.String())
	}
	return 0
}

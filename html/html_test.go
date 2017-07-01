package html

import (
	"fmt"
	"testing"
)

const markdown_1 = "# OK\n\n" +
	"```\n" +
	"# TOMLの例\n" +
	"[DB]\n" +
	"user = \"fkm\"\n" +
	"pass = \"moke\"\n" +
	"db = \"tech\"\n" +
	"```\n"

const markdown_2 = `
## 起動モード

[タスク](./task.html)で、生成したActivityはタスクとよばれるスタックに積まれることがわかりました。

IntentをActivityに対して投げるとActivityが起動しますが、Androidでは次の4つの起動モードが存在します。

 * Standard(デフォルト)
 * SingleTop
 * SingleTask
 * SingleInstance

ここでは、この起動モードについて説明します。

## Standard

これはデフォルトの起動モードです。Activityインスタンスを必ず生成し、Intentを投げたタスクの上に積みます。Standard設定しているActivityで、自分自身へのIntentを投げると同じActivityがどんどんスタックに積まれ、バックボタンを押すと1つずつ戻る動作をします。

## SingleTop

Standardと同じように、Intentを投げたタスクの上にActivityを積もうとしますが、もし該当Activityがタスクの一番上にいた場合(=自分自身を呼び出した場合)は、新しくActivityインスタンスを生成せずにonNewIntent()が代わりに呼ばれます。

タスクの一番上にいない場合は新しいActivityインスタンスを生成しタスクに積むので、A->B->A->BのようなActivity遷移を行うと、Activityインスタンスは複数生成できます。

## SingleTask

ここまでの2つはIntentを投げたタスクの上にActivityを積みましたが、SingleTaskなActivityは

 * 同じタスクがなければ新しいタスクを作ってそれをフォアグラウンドにする
 * 同じタスクがいれば、そのActivityの上にのっているActivityを全部破棄してフォアグラウンドにする

という動作をします。Androidの標準カメラアプリなどはホームから起動するActivityがこのSingleTaskに設定されています。なぜなら

 * カメラアプリからプレビューに遷移し、タスク切り替えで戻ってくるとプレビューがちゃんと表示される
 * ホームから起動すると、常に撮影のActivityが起動する（たとえプレビュー画面で中断していたとしても）

という動作をしているからです。カメラ撮影はなるべく早くやりたい動作ですからね。

ホームから起動されるActivityにこのSingleTaskをつけるときは、「ユーザーはすぐこのタスクを始めることができるべきか」に注意しましょう。

## SingleInstance

SingleTaskの動作に加え、スタックの上にActivityを1つも置けなくしたのがこのSingleInstanceです。ブラウザあたりがこのSingleInstanceに設定されています。なぜなら

 * 他のアプリからブラウザを起動しようとすると、必ずタスク切り替えになる
 * ブラウザから他のActivityを起動しようとすると、必ず別タスクでActivityが起動する

という動作をしているからです。

4つの起動モードを説明しましたが、イメージできたでしょうか？
`

const markdown_3 = `
Androidのアカウント管理の仕組みを使おう

 - [Authenticatorを実装する](./step1_authenticator.html)
 - [AuthenticationServiceを作成する](./step2_service.html)
 - [Authenticator用のXMLを作成する](./step3_xml.html)
 - [AndroidManifest.xmlにServiceを追加する](./step4_manifest.html)
 - [アカウント追加用のログイン画面を作成する](./step5_login_page.html)
 - [追加したアカウントを取得する](./step6_get_account.html)
 - [アカウント選択ダイアログを表示する](./step7_choose_account.html)
 - [アクセストークン取得を実装する](./step8_get_token.html)
`

func Test_OK(t *testing.T) {
	src := "# OK\n\nThis is sample text\n\n![img](./image.webp)"
	m := NewMarkdown()
	out, err := m.Compile(src)
	if err != nil {
		t.Errorf("error : %s", err)
		return
	}
	fmt.Printf(out)
}

func Test_1(t *testing.T) {
	m := NewMarkdown()
	out, err := m.Compile(markdown_1)
	if err != nil {
		t.Errorf("error : %s", err)
		return
	}
	fmt.Printf(out)
}

func Test_2(t *testing.T) {
	m := NewMarkdown()
	out, err := m.Compile(markdown_2)
	if err != nil {
		t.Errorf("error : %s", err)
		return
	}
	fmt.Printf(out)
}

func Test_3(t *testing.T) {
	m := NewMarkdown()
	out, err := m.Compile(markdown_3)
	if err != nil {
		t.Errorf("error : %s", err)
		return
	}
	fmt.Printf(out)
}

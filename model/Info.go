package model

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type Info struct {
	Email        string
	Nome         string
	Valor        float64
	Tipo         string
	Destinatario string
}

func (i Info) Send() {

	t := template.New("template.html")

	var err error
	t, err = t.ParseFiles("template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Print(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "Ccoin@compasso.com.br")
	m.SetHeader("To", i.Email)
	if i.Tipo == "transferencia" {
		m.SetHeader("Subject", "Transferencia efetuada")
	}
	m.SetHeader("Subject", "Compra efetuada")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "joaoterceiro366@gmail.com", "magica123")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

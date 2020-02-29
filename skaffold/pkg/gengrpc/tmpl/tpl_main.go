package tmpl

const MainCodeTpl = `package main

import "{{ . }}/cmd"

func main(){

	cmd.Execute()

}`

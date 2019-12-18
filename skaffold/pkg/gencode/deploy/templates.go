package deploy

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/golang/glog"
)

func genTemplates(companyName, srvName string) error {

	//1: service file
	genServiceFile(srvName)
	//2: deployment file
	genDeploymentFile(companyName, srvName)
	//3: configMap file
	genConfigMapFile(srvName)

	return nil
}

const serviceFile = `
apiVersion: v1
kind: Service
metadata:
  name: {{ . }}
  labels:
    app: {{ . }}-injector
spec:
  ports:
  - port: 443
    targetPort: 443
  selector:
    app: {{ . }}-injector`
const servicePath = "deployment/templates/service.yaml"

func genServiceFile(srvName string) {
	tpl, err := template.New("service").Parse(serviceFile)
	if err != nil {
		glog.Error(err)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	if err = tpl.Execute(buf, &srvName); err != nil {
		glog.Error(err)
		return
	}
	if err := ioutil.WriteFile(servicePath, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return
	}
}

const deploymentFile = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .SrvName }}-deploy
  labels:
    app: {{ .SrvName }}-injector
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .SrvName }}-injector
  template:
    metadata:
      annotations:
        timstamp: "\{\{ date "20060102150405" .Release.Time \}\}"
      labels:
        app: {{ .SrvName }}-injector
    spec:
      containers:
        - name: {{ .SrvName }}-injector
          image: {{ .CompanyName }}/{{ .SrvName }}:\{\{ .Values.app.version \}\}
          resources:
            requests:
              memory: "50Mi"
              cpu: "300m"
            limits:
              memory: "50Mi"
              cpu: "300m"
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: {{ .SrvName }}-certs
              mountPath: /etc/hook/certs
              readOnly: true
            - name: {{ .SrvName }}-config
              mountPath: /etc/hook/config
      volumes:
        - name: {{ .SrvName }}-certs
          secret:
            secretName: {{ .SrvName }}-certs
        - name: {{ .SrvName }}-config
          configMap:
            name: {{ .SrvName }}-injector-config`
const deployPath = "deployment/templates/deployment.yaml"

func genDeploymentFile(companyName, srvName string) {
	tpl, err := template.New("deployment").Parse(deploymentFile)
	if err != nil {
		glog.Error(err)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	arg := struct {
		CompanyName string
		SrvName     string
	}{
		CompanyName: companyName,
		SrvName:     srvName,
	}
	if err := tpl.Execute(buf, &arg); err != nil {
		glog.Error(err)
		return
	}
	if err := ioutil.WriteFile(deployPath, []byte(strings.Replace(buf.String(), "\\", "", -1)), 0666); err != nil {
		glog.Error(err)
		return
	}
}

const configmapFile = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ . }}-injector-config
data:
  injector.yaml: |
    initContainers:
      - name: kms-letter
        image: alpine:latest
        imagePullPolicy: IfNotPresent
        command: ["/bin/sleep","1d"]`

const configmapPath = "deployment/templates/configmap.yaml"

func genConfigMapFile(srvName string) {
	tpl, err := template.New("configmap").Parse(configmapFile)
	if err != nil {
		glog.Error(err)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, &srvName); err != nil {
		glog.Error(err)
		return
	}
	if err := ioutil.WriteFile(configmapPath, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return
	}
}

//go:build js && wasm
// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	document = js.Global().Get("document")
)

func getElementByID(id string) js.Value {
	return document.Call("getElementById", id)
}

func renderTextArea(s string) string {
	return fmt.Sprintf(`<textarea style="width: 100%%;height: 400px">%s</textarea>`, s)
}

func renderErrorTextArea(err error) string {
	return fmt.Sprintf(`<textarea style="width: 100%%;height: 400px; border: 1px solid red;">%s</textarea>`, err.Error())
}

func deserialize(data []byte) (runtime.Object, error) {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func getDiffPatch(current, updated runtime.Object) (string, error) {
	currentBytes, err := json.Marshal(current)
	if err != nil {
		return "", err
	}
	updatedBytes, err := json.Marshal(updated)
	if err != nil {
		return "", err
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(currentBytes, updatedBytes, current)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(patch, &m); err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func main() {
	quit := make(chan struct{})

	originYaml := getElementByID("origin")
	desiredYaml := getElementByID("desired")
	preview := getElementByID("preview")
	diffButton := getElementByID("diff")
	diffButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		origin := originYaml.Get("value").String()
		originObj, err := deserialize([]byte(origin))
		if err != nil {
			preview.Set("innerHTML", renderErrorTextArea(err))
			return nil
		}
		desired := desiredYaml.Get("value").String()
		desiredObj, err := deserialize([]byte(desired))
		if err != nil {
			preview.Set("innerHTML", renderErrorTextArea(err))
			return nil
		}
		diff, err := getDiffPatch(originObj, desiredObj)
		if err != nil {
			preview.Set("innerHTML", renderErrorTextArea(err))
			return nil
		}
		preview.Set("innerHTML", renderTextArea(diff))
		return nil
	}))

	<-quit
	originYaml.Call("remove")
	desiredYaml.Call("remove")
	preview.Call("remove")
}

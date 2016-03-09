package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	/*"log"*/
	"net/http"

	"strconv"

	"github.com/kubeexcutor/api"
	"github.com/kubeexcutor/g"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.VERSION))
	})
	//deal with the imformation of replication controll
	http.HandleFunc("/replication", func(w http.ResponseWriter, r *http.Request) {
		//receive the data of kubecontr assembly and judgment processing
		//judge the parameter if nil print the log
		err := r.ParseForm()
		if err != nil {
			log.Fatalln("read replication info fail:", err)
		}
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			if len(r.Form["strurl"]) != 0 && len(r.Form["rcnumber"]) != 0 {

				strurl := r.Form["strurl"][0]
				rcnumber, error := strconv.Atoi(r.Form["rcnumber"][0])
				if error != nil {
					fmt.Println("string to int failed")
				}
				fmt.Println("the strurl is:", r.Form["strurl"][0])
				fmt.Println("the rcnumber is:", r.Form["rcnumber"][0])

				var rc api.ReplicationController = getRCJson(strurl)
				//change the num of replication
				rc.Spec.Replicas = rcnumber
				fmt.Println("the rc number is change to:", r.Form["rcnumber"][0])
				//change the resource version set the version nil
				rc.ObjectMeta.ResourceVersion = ""
				sendRcJson(strurl, rc)

				w.Write([]byte("rc number is change to:" + r.Form["rcnumber"][0] + "\n"))

			} else {
				log.Fatalln("rc info error, strurl:", r.Form["strurl"], " rcnumber: ", r.Form["rcnumber"])
			}
		}

	})

	//change the prefence of the rc
	http.HandleFunc("/rcprefence", func(w http.ResponseWriter, r *http.Request) {
		//receive the data of kubecontr assembly and judgment processing
		//judge the parameter if nil print the log
		err := r.ParseForm()
		if err != nil {
			log.Fatalln("read replication info fail:", err)
		}
		if r.Method == "GET" {

		} else if r.Method == "POST" {
			if len(r.Form["strurl"]) != 0 {

				strurl := r.Form["strurl"][0]
				limcpu := r.Form["limcpu"][0]
				limmem := r.Form["limmem"][0]
				reqcpu := r.Form["reqcpu"][0]
				reqmem := r.Form["reqmem"][0]
				fmt.Println(limcpu, limmem, reqcpu, reqmem)
				
				var rc api.ReplicationController = getRCJson(strurl)
				//change the prefence of replication

				limits := make(map[string]string)
				// 然后赋值
				limits["cpu"] = limcpu
				limits["memory"] = limmem
				requests := make(map[string]string)
				requests["cpu"] = reqcpu
				requests["memory"] = reqmem

				rc.Spec.Template.Spec.Containers[0].Resources.Limits = limits
				rc.Spec.Template.Spec.Containers[0].Resources.Requests = requests
				//change the resource version set the version nil
				rc.ObjectMeta.ResourceVersion = ""
				sendRcJson(strurl, rc)

				w.Write([]byte("rc number prefence change to:" + limcpu + limmem + reqcpu + reqmem + "\n"))

			} else {
				log.Fatalln("rc info error, strurl:", r.Form["strurl"], " rcnumber: ", r.Form["rcnumber"])
			}
		}

	})

}
func getRCJson(strurl string) api.ReplicationController {
	resp, err := http.Get(strurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", data)

	//analysis the  json format data to the struct
	var rc api.ReplicationController
	err = json.Unmarshal(data, &rc)
	if err != nil {
		panic(err)
	}
	return rc
}
func sendRcJson(strurl string, rc api.ReplicationController) {
	//convert the struct to the json formt
	b, err := json.Marshal(rc)
	if err != nil {
		panic(err)
	}

	body := bytes.NewBuffer([]byte(b))

	fmt.Printf("%s", body)

	//strurl := "http://192.168.11.58:8080/api/v1/namespaces/default/replicationcontrollers/spark-master-controller"
	var request1 *http.Request
	request1, err = http.NewRequest("PUT", strurl, body)
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, strurl)
	}
	request1.Header.Set("Accept", "application/json")

	fmt.Println("\n start to send to kubernetes api:")

	var resp1 *http.Response
	resp1, err = http.DefaultClient.Do(request1)
	if err != nil {
		//todo
		log.Println("http.Do failed,[err=%s][url=%s]", err, strurl)
	}
	defer resp1.Body.Close()
	fmt.Println("end to send to kubernetes api:")
	b, err = ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println("faild to send to kubernetes api:")
		panic(err)
	}
	fmt.Println("the return of kubernetes:")
	fmt.Printf("%s\n", b)
}

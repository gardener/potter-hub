package apiProxyTest

import (
	"crypto/tls"
	"errors"
	"fmt"
	"integration-test/pkg/helmProxyWrapper"
	"integration-test/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var releaseName = "test"
var chartName = "grafana"
var version = "7.9.11"
var apprepository = "bitnami"
var client = http.Client{}

func Run(helmProxyWrapper *helmProxyWrapper.HelmProxyWrapper) {
	helmProxyWrapper.InstallRelease(releaseName, chartName, version, apprepository)

	err := testGetNamespace(helmProxyWrapper.HTTPSBaseURL, helmProxyWrapper.Token)
	if err != nil {
		helmProxyWrapper.DeleteRelease(releaseName)
		log.Fatal(err)
	}
	err = testGetResources(helmProxyWrapper.HTTPSBaseURL, helmProxyWrapper.Token, helmProxyWrapper.Namespace)
	if err != nil {
		helmProxyWrapper.DeleteRelease(releaseName)
		log.Fatal(err)
	}
	err = testWS(helmProxyWrapper.BaseURL, helmProxyWrapper.Token, helmProxyWrapper.Namespace)
	if err != nil {
		helmProxyWrapper.DeleteRelease(releaseName)
		log.Fatal(err)
	}
	err = testGetAppRepossitories(helmProxyWrapper.HTTPSBaseURL, helmProxyWrapper.Token)
	if err != nil {
		helmProxyWrapper.DeleteRelease(releaseName)
		log.Fatal(err)
	}

	helmProxyWrapper.DeleteRelease(releaseName)
	log.Println("Test finished successfully")
}

func makeWSRequest(wsURL url.URL, token string) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	d := websocket.Dialer{}
	d.Subprotocols = append(d.Subprotocols, "base64url.bearer.authorization.k8s.io."+token+", binary.k8s.io")
	d.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	c, _, err := d.Dial(wsURL.String(), nil)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-timer.C:
			log.Println("timer hit. Closing connection.")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)

				return fmt.Errorf("connection closed unexpectedly")
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func testWS(baseURL, token, namespace string) error {
	queryParams := "watch=true&fieldSelector=metadata.name=" + releaseName + "-grafana" + "&Namespace=" + util.PROJECT_NAME + "&SecretName=" + util.SECRET_NAME

	deploymentWsPath := "/" + util.PROJECT_NAME + "/" + util.SECRET_NAME + util.APIKube + util.DEPLOYMENT_API + "/namespaces/" + namespace + "/deployments"
	deploymentWsURL := url.URL{Scheme: "wss", Host: baseURL, Path: deploymentWsPath, RawQuery: queryParams}

	serviceWsPath := "/" + util.PROJECT_NAME + "/" + util.SECRET_NAME + util.APIKube + util.CORE_V1_API + "/namespaces/" + namespace + "/services"
	serviceWsURL := url.URL{Scheme: "wss", Host: baseURL, Path: serviceWsPath, RawQuery: queryParams}

	secretWsPath := "/" + util.PROJECT_NAME + "/" + util.SECRET_NAME + util.APIKube + util.CORE_V1_API + "/namespaces/" + namespace + "/secrets"
	log.Println(secretWsPath)
	secretWsURL := url.URL{Scheme: "wss", Host: baseURL, Path: secretWsPath, RawQuery: queryParams}

	log.Info("Doing deployment ws request")
	if err := makeWSRequest(deploymentWsURL, token); err != nil {
		return err
	}
	fmt.Println()

	log.Info("Doing service ws request")
	if err := makeWSRequest(serviceWsURL, token); err != nil {
		return err
	}
	fmt.Println()

	log.Info("Doing secret ws request")
	if err := makeWSRequest(secretWsURL, token); err != nil {
		return err
	}
	fmt.Println()

	return nil
}

func testGetAppRepossitories(httpsBaseURL string, token string) error {
	apprepoURL := httpsBaseURL + "/api/ui-backend/apprepositories"

	log.Println("Doing apprepo request with correct auth header -> 200")
	req := util.BuildApiProxyRequestWithBearerAuth("GET", apprepoURL, token, nil)
	if err := makeRequestWithExpectedResult(req, 200); err != nil {
		return err
	}

	// TODO: Inline as soon as auth check in k8s-api-proxy is implemented
	//log.Println("Doing apprepo request with incorrect auth header -> 401")
	//req = util.BuildApiProxyRequestWithBearerAuth("GET", apprepoURL, "foobar", nil)
	//err = makeRequestWithExpectedResult(req, 401)
	//
	//log.Println("Doing apprepo request without auth header -> 400")
	//req = util.BuildApiProxyRequestWithoutAuth("GET", apprepoURL, nil)
	//err = makeRequestWithExpectedResult(req, 400)
	//
	//apprepoURL = httpsBaseURL + "/apis/kubeapps.com/v1alpha1/namespaces/" + util.TEST_NAMESPACE + "/apprepositories"
	//log.Println("Doing apprepo request in wrong namespace -> 403")
	//req = util.BuildApiProxyRequestWithoutAuth("GET", apprepoURL, nil)
	//err = makeRequestWithExpectedResult(req, 403)
	return nil
}

func testGetResources(httpsBaseURL, token, namespace string) error {
	deploymentURL := httpsBaseURL + util.APIKube + util.DEPLOYMENT_API + "/namespaces/" + namespace + "/deployments/" + releaseName + "-grafana"
	secretURL := httpsBaseURL + util.APIKube + util.CORE_V1_API + "/namespaces/" + namespace + "/secrets/" + releaseName + "-grafana-admin"
	serviceURL := httpsBaseURL + util.APIKube + util.CORE_V1_API + "/namespaces/" + namespace + "/services/" + releaseName + "-grafana"

	log.Infof("Doing get deployment request with auth")
	//log.Infof("DEBUG:\n  url=\n%s\n  token=\n%s", deploymentURL, token)
	req := util.BuildApiProxyRequestWithBearerAuth("GET", deploymentURL, token, nil)
	if err := makeRequestWithExpectedResult(req, 200); err != nil {
		return err
	}

	log.Infof("Doing get deployment request with invalid auth")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", deploymentURL, "foobar", nil)
	if err := makeRequestWithExpectedResult(req, 401); err != nil {
		return err
	}

	log.Infof("Doing get deployment request without auth")
	req = util.BuildApiProxyRequestWithoutAuth("GET", deploymentURL, nil)
	if err := makeRequestWithExpectedResult(req, 400); err != nil {
		return err
	}

	log.Infof("Doing get secret request with auth")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", secretURL, token, nil)
	if err := makeRequestWithExpectedResult(req, 200); err != nil {
		return err
	}

	log.Infof("Doing get secret request with invalid auth")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", secretURL, "foobar", nil)
	if err := makeRequestWithExpectedResult(req, 401); err != nil {
		return err
	}

	log.Infof("Doing get secret request without auth")
	req = util.BuildApiProxyRequestWithoutAuth("GET", secretURL, nil)
	if err := makeRequestWithExpectedResult(req, 400); err != nil {
		return err
	}

	log.Infof("Doing get service request with auth")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", serviceURL, token, nil)
	if err := makeRequestWithExpectedResult(req, 200); err != nil {
		return err
	}

	log.Infof("Doing get service request with invalid auth")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", serviceURL, "foobar", nil)
	if err := makeRequestWithExpectedResult(req, 401); err != nil {
		return err
	}

	log.Infof("Doing get service request without auth")
	req = util.BuildApiProxyRequestWithoutAuth("GET", serviceURL, nil)
	if err := makeRequestWithExpectedResult(req, 400); err != nil {
		return err
	}

	return nil
}

func testGetNamespace(httpsBaseURL string, token string) error {
	namespaceURL := httpsBaseURL + "/api/ui-backend/reverse-proxy/v1/api/v1/namespaces/"

	log.Infof("Doing get ns request with auth")
	req := util.BuildApiProxyRequestWithBearerAuth("GET", namespaceURL, token, nil)
	if err := makeRequestWithExpectedResult(req, 200); err != nil {
		return err
	}

	log.Infof("Doing get ns request with incorrect auth -> Unauthorized")
	req = util.BuildApiProxyRequestWithBearerAuth("GET", namespaceURL, "foobar", nil)
	if err := makeRequestWithExpectedResult(req, 401); err != nil {
		return err
	}

	log.Infof("Doing get ns request without auth -> Bad Request")
	req = util.BuildApiProxyRequestWithoutAuth("GET", namespaceURL, nil)
	if err := makeRequestWithExpectedResult(req, 400); err != nil {
		return err
	}

	return nil
}

func makeRequestWithExpectedResult(req *http.Request, expectedStatusCode int) error {
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	if res == nil {
		return errors.New("response is empty")
	}
	defer util.CloseBody(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %w", err)
	}

	if res.StatusCode != expectedStatusCode {
		errMsg := fmt.Sprintf("Expected statusCode does not match response statusCode. "+
			"Expected: %v, Actual: %v", expectedStatusCode, res.StatusCode)
		log.Error(errMsg)
		log.Errorf("Request body was: %s", string(body))
		return errors.New(errMsg)
	}

	log.Printf("Call returned with expected statusCode %v", res.StatusCode)

	defer fmt.Println()

	return nil
}

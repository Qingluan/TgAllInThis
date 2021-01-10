package actions

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/alyu/configparser"

	tdlib "github.com/Arman92/go-tdlib"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	Tp   string
	IP   string
	PORT string
}

var (
	Cha = make(chan os.Signal, 2)
)

func WaitInterrupt(client *tdlib.Client) {
	// Handle Ctrl+C , Gracefully exit and shutdown tdlib

	signal.Notify(Cha, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-Cha
		client.DestroyInstance()
		os.Exit(1)
	}()

}

func AuthClient(authSeciton *configparser.Section) (*tdlib.Client, error) {
	tdlib.SetLogVerbosityLevel(1)
	basedir := authSeciton.ValueOf("tddir")

	api, apihash, db, dbdir := authSeciton.ValueOf("api"), authSeciton.ValueOf("apihash"), authSeciton.ValueOf("tddb"), authSeciton.ValueOf("tddbdir")
	errorText := authSeciton.ValueOf("err_log")
	errorText = filepath.Join(basedir, errorText)
	db = filepath.Join(basedir, db)
	dbdir = filepath.Join(basedir, dbdir)

	logLevelStr := authSeciton.ValueOf("log_level")
	logLevel, _ := strconv.Atoi(logLevelStr)

	tdlib.SetFilePath(errorText)
	tdlib.SetLogVerbosityLevel(logLevel)
	// tdlib.SetFilePath(errorText)

	ptp := authSeciton.ValueOf("proxyTp")
	proxy := Proxy{}
	if ptp != "" {
		proxy.Tp = ptp
		proxy.IP = authSeciton.ValueOf("proxyIP")
		proxy.PORT = authSeciton.ValueOf("proxyPort")
	}

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               api,
		APIHash:             apihash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   db,
		FileDirectory:       dbdir,
		IgnoreFileNames:     false,
	})
	if proxy.Tp != "" {
		logrus.Info("Try Connect Proxy : ", proxy.Tp, "://", proxy.IP, ":", proxy.PORT)
		port, _ := strconv.Atoi(proxy.PORT)
		if _, err := client.AddProxy(proxy.IP, int32(port), true, tdlib.NewProxyTypeSocks5("", "")); err != nil {
			logrus.Error("Proxy add error:", err)
			// return nil, err
		} else {
			logrus.Info("Use Proxy: ", proxy.Tp, "://", proxy.IP, ":", proxy.PORT)

		}
	}

	// Wait while we get AuthorizationReady!
	// Note: See authorization example for complete auhtorization sequence example
	// currentState, _ := client.Authorize()
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	// for ; currentState.GetAuthorizationStateEnum() != tdlib.AuthorizationStateReadyType; currentState, _ = client.Authorize() {
	// 	time.Sleep(300 * time.Millisecond)
	// }
	return client, nil
}

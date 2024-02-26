package util

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/registry"
)

// key
const (
	currentUser registry.Key = registry.CURRENT_USER
)

// path
const (
	run             string = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	internetSetting string = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`
)

// sub key name
const (
	proxyEnable   string = "ProxyEnable"   //DWord
	proxyServer   string = "ProxyServer"   //string
	autoConfigURL string = "AutoConfigURL" //string
	proxyOverride string = "ProxyOverride" //string
)

func openKey(key registry.Key, path string) (registry.Key, error) {
	return registry.OpenKey(key, path, registry.ALL_ACCESS)
}

func ClearProxy() error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue(proxyEnable, 0)
	if err != nil {
		return err
	}
	err = key.DeleteValue(proxyServer)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.DeleteValue(autoConfigURL)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.DeleteValue(proxyOverride)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func ReloadPAC() error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	url, _, err := key.GetStringValue(autoConfigURL)
	if err != nil {
		return err
	}
	if url == "" {
		return nil
	}
	err = key.DeleteValue(autoConfigURL)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	go func() {
		// Make browser quickly request pac file again
		// Browser not send request if edit key synchronize
		defer key.Close()
		time.Sleep(1 * time.Second)
		err = key.SetStringValue(autoConfigURL, url)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func EnablePAC(ip string, port string) error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue(proxyEnable, 0)
	if err != nil {
		return err
	}
	err = key.DeleteValue(proxyServer)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.SetStringValue(autoConfigURL, "http://"+ip+":"+port)
	if err != nil {
		return err
	}
	err = key.DeleteValue(proxyOverride)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func EnableHTTP(ip string, port string) error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue(proxyEnable, 1)
	if err != nil {
		return err
	}
	err = key.SetStringValue(proxyServer, ip+":"+port)
	if err != nil {
		return err
	}
	err = key.DeleteValue(autoConfigURL)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.DeleteValue(proxyOverride)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func EnableSOCKS4(ip string, port string) error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue(proxyEnable, 1)
	if err != nil {
		return err
	}
	err = key.SetStringValue(proxyServer, "socks="+ip+":"+port)
	if err != nil {
		return err
	}
	err = key.DeleteValue(autoConfigURL)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.DeleteValue(proxyOverride)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func EnableSOCKS5(ip string, port string) error {
	key, err := openKey(currentUser, internetSetting)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetDWordValue(proxyEnable, 1)
	if err != nil {
		return err
	}
	err = key.SetStringValue(proxyServer, "socks://"+ip+":"+port)
	if err != nil {
		return err
	}
	err = key.DeleteValue(autoConfigURL)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	err = key.DeleteValue(proxyOverride)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func HookAutoStart(checked bool, appName string) error {
	key, err := openKey(currentUser, run)
	if err != nil {
		return err
	}
	defer key.Close()
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	run := "cmd /c cd " + filepath.Dir(ex) + " && start " + filepath.Base(ex)
	if checked {
		// Set the application to auto start
		if err := key.SetStringValue(appName, run); err != nil {
			return err
		}
	} else {
		// Remove the application from auto start
		if err := key.DeleteValue(appName); err != nil {
			if !errors.Is(err, registry.ErrNotExist) {
				return err
			}
		}
	}
	return nil
}

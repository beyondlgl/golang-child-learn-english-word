package service

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const youdaoDictVoiceURL = "https://dict.youdao.com/dictvoice"

// DictVoiceHandler 代理有道词典发音接口
func DictVoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	audio := r.URL.Query().Get("audio")
	if audio == "" {
		http.Error(w, "audio is required", http.StatusBadRequest)
		return
	}

	dictType := r.URL.Query().Get("type")
	if dictType == "" {
		dictType = "0"
	}

	params := url.Values{}
	params.Set("audio", audio)
	params.Set("type", dictType)
	targetURL := youdaoDictVoiceURL + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to call upstream service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	} else {
		w.Header().Set("Content-Type", "audio/mpeg")
	}

	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl != "" {
		w.Header().Set("Cache-Control", cacheControl)
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

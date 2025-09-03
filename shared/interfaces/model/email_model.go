package model

import (
	"bytes"
	"encoding/json"
	"text/template"
	"time"
)

type (
	Email struct {
		User_id uint
		Subject string
		Content string
	}
)

const EmailTemplate = `Terima kasih Anda telah menggunakan fasilitas kami. 
Berikut ini adalah informasi aktivitas yang telah Anda lakukan:

Waktu    : {{ .Time }} WIB
Aktivitas: {{ .Activity }}
Status   : {{ .Status }}

Semoga informasi ini bermanfaat bagi Anda.
Jika Anda merasa tidak melakukan aktivitas ini atau memiliki pertanyaan mengenai aktivitas Anda, silahkan hubungi kami.

Terima kasih.`

func (m *Email) ToJsonString() (string, error) {
	if m == nil {
		return "", nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func ToTemplateEmail(activity string, status string) (string, error) {
	data := struct {
		Time     string
		Activity string
		Status   string
	}{
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Activity: activity,
		Status:   status,
	}

	t, err := template.New("email").Parse(EmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
